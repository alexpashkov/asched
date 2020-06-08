package amenities

import (
	"context"
	"github.com/alexpashkov/asched/graph/model"
	"github.com/alexpashkov/asched/internal/photos"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	mongoClient *mongo.Client
	mongoDBName string
	photosSrv   *photos.Service
}

func NewService(mongoClient *mongo.Client, mongoDBName string, photosSrv *photos.Service) *Service {
	return &Service{
		mongoClient: mongoClient,
		mongoDBName: mongoDBName,
		photosSrv:   photosSrv,
	}
}

func (s *Service) AddAmenity(ctx context.Context, newAm model.NewAmenity) (string, error) {
	var id string
	return id, s.mongoClient.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		res, err := s.mongoCollection().InsertOne(ctx, newAmenityToMongoAmenity(newAm))
		if err != nil {
			return err
		}
		oid, ok := res.InsertedID.(primitive.ObjectID)
		if !ok {
			return errors.New("InsertedID is not ObjectID")
		}
		id = oid.Hex()
		if newAm.Photo != nil && newAm.Photo.File != nil {
			return errors.Wrap(
				s.photosSrv.SavePhoto(id, newAm.Photo.File),
				"failed to save the photo",
			)
		}
		return nil
	})
}

func (s *Service) mongoCollection() *mongo.Collection {
	return s.mongoClient.Database(s.mongoDBName).Collection("amenities")
}

// mongoAmenity is a structure used to persist Amenity in a MongoDB database
type mongoAmenity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Type        string             `bson:"type"`
	Location    location           `bson:"location"`
	Description *string            `bson:"description"`
}

type location struct {
	Type        string     `json:"type"`
	Coordinates [2]float64 `json:"coordinates"`
}

func newAmenityToMongoAmenity(newAm model.NewAmenity) mongoAmenity {
	return mongoAmenity{
		Name: newAm.Name,
		Type: newAm.Type,
		Location: location{
			Type:        "Point",
			Coordinates: [2]float64{newAm.Latitude, newAm.Longitude},
		},
		Description: nil,
	}
}
