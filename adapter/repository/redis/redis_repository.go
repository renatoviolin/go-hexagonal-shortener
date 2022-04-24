package redis

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/renatoviolin/shortener/application/entity"
)

type RedisRepository struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewRedisRepository(redisURL string) (*RedisRepository, error) {
	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, err
	}

	return &RedisRepository{client: client}, nil
}

func (r *RedisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *RedisRepository) Find(code string) (*entity.Redirect, error) {
	redirect := &entity.Redirect{}
	key := r.generateKey(code)
	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, entity.ErrRedirectNotFound
	}

	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, err
	}

	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt
	return redirect, nil
}

func (r *RedisRepository) Store(redirect *entity.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}

	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return err
	}

	return nil
}
