package mongo

import (
	"context"
	"github.com/alexpashkov/asched/internal/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context) (*mongo.Client, error) {
	return mongo.Connect(ctx, options.Client().ApplyURI(env.MONGODB_URI()))
}

func ConnectPing(ctx context.Context) (*mongo.Client, error) {
	c, err := Connect(ctx)
	if err != nil {
		return nil, err
	}
	return c, c.Ping(ctx, nil)
}
