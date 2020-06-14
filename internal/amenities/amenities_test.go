package amenities

import (
	"context"
	"encoding/json"
	"github.com/alexpashkov/asched/graph/model"
	"github.com/alexpashkov/asched/internal/config"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func TestAmenitiesService(t *testing.T) {
	if testing.Short() {
		t.Skip("short testing enabled")
	}
	conf, err := config.ReadConfig(t.Logf)
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	t.Cleanup(cancel)
	client, err := mongo.Connect(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, client.Disconnect(ctx))
	})
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	t.Cleanup(cancel)
	confBuf, err := json.MarshalIndent(conf, "", "\t")
	t.Logf("connecting to MongoDB, conf: %s", confBuf)
	require.NoError(t, client.Ping(ctx, nil))
	s := NewService(client, conf.MongoDBConnString.Database, nil)
	t.Run("add amenity", func(t *testing.T) {
		id, err := s.AddAmenity(ctx, model.NewAmenity{
			Name: time.Now().Format(time.UnixDate),
			Type: "TennisCourt",
			Lat:  0,
			Lon:  0,
		})
		require.NoError(t, err)
		t.Log("amenity", id, "created")
	})
}
