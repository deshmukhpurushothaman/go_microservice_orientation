package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/deshmukhpurushothaman/golang_microservice_orientation/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create context to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// start web server
	// go app.serve()
	app.serve()

	// log.Println("Starting service on port ", webPort)
	// srv := &http.Server{
	// 	Addr:    fmt.Sprintf(":%s", webPort),
	// 	Handler: app.routes(),
	// }

	// err := srv.ListenAndServe()
	// if err != nil {
	// 	log.Panic(err)
	// }
}

func (app *Config) serve() {
	log.Println("Starting service on port ", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connection:", err)
		return nil, err
	}

	log.Println("Connected to mongo!")

	return c, nil
}
