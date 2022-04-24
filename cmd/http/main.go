package main

import (
	"log"
	"sync"

	grpcAdapter "github.com/renatoviolin/shortener/adapter/grpc"
	chiAdapter "github.com/renatoviolin/shortener/adapter/http/chi"
	ginAdapter "github.com/renatoviolin/shortener/adapter/http/gin"
	"github.com/renatoviolin/shortener/adapter/repository/redis"
	"github.com/renatoviolin/shortener/application/shortener"
)

func main() {
	// 1. Create a Repository (redis or mongo)
	// repository, err := mongodb.NewMongoRepository("mongodb://localhost:27017", "go_projects", 5)
	repository, err := redis.NewRedisRepository("redis://localhost:6379")
	if err != nil {
		log.Fatal(err.Error())
	}

	// 2. Instantiate the Application Service
	service := shortener.NewRedirectService(repository)

	// 3. Instantiate the UseCase that uses the service
	useCase := shortener.NewUseCaseShortener(service)

	// Wait Group to wait each Goroutine
	wg := sync.WaitGroup{}
	wg.Add(3)

	// 4.1 Create a GRPC Server :7000
	grpcHandler := grpcAdapter.NewGrpcHandler(*useCase)
	go func() {
		grpcHandler.Run(":7000")
		wg.Done()
	}()

	// 4.2 Create a HTTP handler (CHI) :8000
	chiHandler := chiAdapter.NewChiHandler(*useCase)
	chiHandler.SetupRoutes()
	go func() {
		chiHandler.Run(":8000")
		wg.Done()
	}()

	// 4.3. Create another HTTP handler (GIN) :9000
	ginHandler := ginAdapter.NewGinHandler(*useCase)
	ginHandler.SetupRoutes()
	go func() {
		ginHandler.Run(":9000")
		wg.Done()
	}()

	wg.Wait()
}
