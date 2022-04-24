package mongodb

import (
	"context"
	"time"

	"github.com/renatoviolin/shortener/application/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newMongoClient(mongoURL string, timeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, err
}

func NewMongoRepository(mongoURL, mongoDB string, timeout int) (*MongoRepository, error) {
	repo := &MongoRepository{
		database: mongoDB,
		timeout:  time.Duration(timeout) * time.Second,
	}

	client, err := newMongoClient(mongoURL, timeout)
	if err != nil {
		return nil, err
	}

	repo.client = client
	return repo, nil
}

func (m *MongoRepository) Find(code string) (*entity.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.timeout)*time.Second)
	defer cancel()

	redirect := &entity.Redirect{}
	collection := m.client.Database(m.database).Collection("redirects")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, entity.ErrRedirectNotFound
		}
		return nil, err
	}

	return redirect, nil
}

func (m *MongoRepository) Store(redirect *entity.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.timeout)*time.Second)
	defer cancel()

	collection := m.client.Database(m.database).Collection("redirects")
	_, err := collection.InsertOne(ctx, bson.M{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}
