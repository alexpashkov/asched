package main

import (
	"context"
	"github.com/alexpashkov/asched/internal/amenities"
	"github.com/alexpashkov/asched/internal/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/alexpashkov/asched/graph"
	"github.com/alexpashkov/asched/graph/generated"
)

func main() {
	logger := logrus.New()
	conf, err := config.ReadConfig(logger.Printf)
	if err != nil {
		logger.Fatal(errors.Wrap(err, "invalid config"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.MongoDBRawConnString))
	if err != nil {
		logger.Fatal(errors.Wrap(err, "failed to create MongoDB client"))
	}

	amenitiesService := amenities.NewService(
		logger.WithField("component", "AmenitiesService"),
		mongoClient,
		conf.MongoDBConnString.Database,
		os.Getenv("PHOTOS_DIR"),
	)
	if err := amenitiesService.Start(ctx); err != nil {
		logger.Fatal(errors.Wrap(err, "failed to start amenities service"))
	}
	cancel()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{Resolvers: &graph.Resolver{AmenitiesService: amenitiesService}},
	))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	logger.Printf("connect to http://localhost:%s/ for GraphQL playground", conf.Port)
	logger.Fatal(http.ListenAndServe(":"+conf.Port, nil))
}
