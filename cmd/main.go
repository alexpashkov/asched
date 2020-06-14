package main

import (
	"context"
	"github.com/alexpashkov/asched/internal/amenities"
	"github.com/alexpashkov/asched/internal/config"
	"github.com/alexpashkov/asched/internal/photos"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/alexpashkov/asched/graph"
	"github.com/alexpashkov/asched/graph/generated"
)

func main() {
	conf, err := config.ReadConfig(log.Printf)
	if err != nil {
		log.Fatal(errors.Wrap(err, "invalid config"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.MongoDBRawConnString))
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create MongoDB client"))
	}
	cancel()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{Resolvers: &graph.Resolver{
			AmenitiesService: amenities.NewService(
				mongoClient,
				conf.MongoDBConnString.Database,
				photos.NewService(os.Getenv("PHOTOS_DIR")),
			)}},
	))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", conf.Port)
	log.Fatal(http.ListenAndServe(":"+conf.Port, nil))
}
