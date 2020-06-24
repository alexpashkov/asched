package amenities

import (
	"bytes"
	"context"
	"crypto/rand"
	"github.com/alexpashkov/asched/graph/model"
	"github.com/alexpashkov/asched/internal/config"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"testing"
	"time"
)

const (
	photosDir = "photos"
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
	t.Logf("connecting to MongoDB, conf: %#v", conf)
	require.NoError(t, client.Ping(ctx, nil))
	s := NewService(client, conf.MongoDBConnString.Database, photosDir)
	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(photosDir))
	})
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
	t.Run("add photo", func(t *testing.T) {
		_, err := s.AddPhoto("foobar", bytes.NewReader(randBytes(t, 1024)))
		require.NoError(t, err)
	})
}

func randBytes(t testing.TB, n int) []byte {
	t.Helper()
	buf := make([]byte, n)
	_, err := rand.Read(buf)
	require.NoError(t, err)
	return buf
}
