package amenities

import (
	"context"
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
	t.Logf("config is %#v", conf)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	t.Cleanup(cancel)
	client, err := mongo.Connect(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, client.Disconnect(ctx))
	})
	s := NewService(client, conf.MongoDBConnString.Database, nil)
	id, err := s.AddAmenity(ctx, model.NewAmenity{
		Name: time.Now().Format(time.UnixDate),
		Type: "TennisCourt",
		Lat:  0,
		Lon:  0,
	})
	require.NoError(t, err)
	t.Log(id)
}
