package main

import (
	"fmt"
	"log"

	"github.com/renatoviolin/shortener/adapter/repository/redis"
	"github.com/renatoviolin/shortener/adapter/serializer"
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

	// 4. Execute the UseCase to encode the URL
	url := "http://www.google.com"
	fmt.Printf("Generating Code from URL: %s\n", url)
	redirectOut, err := useCase.UrlToCode(url)
	if err != nil {
		fmt.Printf("error: %+v", err.Error())
	}
	jsonSerializer := serializer.JSONSerializer{}
	bytes, _ := jsonSerializer.Encode(redirectOut)
	fmt.Printf("%+v\n\n", string(bytes))

	// 5. Execute the UseCase to get back the URL from the code
	fmt.Printf("Retrieve the URL from code: %s\n", redirectOut.Code)
	out, err := useCase.CodeToUrl(redirectOut.Code)
	if err != nil {
		fmt.Printf("error: %+v", err.Error())
	}
	bytes, _ = jsonSerializer.Encode(out)
	fmt.Printf("%+v\n\n", string(bytes))
}
