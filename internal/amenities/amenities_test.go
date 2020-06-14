package amenities

import (
	"context"
	"github.com/alexpashkov/asched/graph/model"
	"github.com/alexpashkov/asched/internal/env"
	"github.com/alexpashkov/asched/internal/mongo"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAmenitiesService(t *testing.T) {
	if testing.Short() {
		t.Skip("short testing enabled")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	t.Cleanup(cancel)
	client, err := mongo.Connect(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, client.Disconnect(ctx))
	})
	s := NewService(client, env.MONGODB_DB_NAME(), nil)
	id, err := s.AddAmenity(ctx, model.NewAmenity{
		Name: time.Now().Format(time.UnixDate),
		Type: "TennisCourt",
		Lat:  0,
		Lon:  0,
	})
	require.NoError(t, err)
	t.Log(id)
}
