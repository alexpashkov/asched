package amenities

import (
	"context"
	"github.com/alexpashkov/asched/graph/model"
	"github.com/alexpashkov/asched/internal/photos"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
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

func (s *Service) SearchAmenities(ctx context.Context, latitude, longitude float64, maxDistance int, typeArg *string) ([]*model.Amenity, error) {
	filter := bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry":    makeLocation(latitude, longitude),
				"$maxDistance": maxDistance,
			},
		},
	}
	if typeArg != nil {
		filter["type"] = *typeArg
	}
	cur, err := s.mongoCollection().Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "db query failed")
	}
	defer cur.Close(ctx)
	var res []*model.Amenity
	for cur.Next(ctx) {
		var mongoAm mongoAmenity
		if err := cur.Decode(&mongoAm); err != nil {
			return nil, err
		}
		am := mongoAmenityToAmenity(mongoAm)
		res = append(res, &am)
	}
	return res, cur.Err()
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

func makeLocation(latitude, longitude float64) location {
	return location{
		Type:        "Point",
		Coordinates: [2]float64{latitude, longitude},
	}
}

type location struct {
	Type        string     `json:"type"`
	Coordinates [2]float64 `json:"coordinates"`
}

func newAmenityToMongoAmenity(newAm model.NewAmenity) mongoAmenity {
	return mongoAmenity{
		Name:        newAm.Name,
		Type:        newAm.Type,
		Location:    makeLocation(newAm.Lat, newAm.Lon),
		Description: nil,
	}
}

func mongoAmenityToAmenity(mongoAm mongoAmenity) model.Amenity {
	var am model.Amenity
	am.ID = mongoAm.ID.Hex()
	am.Name = mongoAm.Name
	am.Type = mongoAm.Type
	am.Lat = mongoAm.Location.Coordinates[0]
	am.Lon = mongoAm.Location.Coordinates[1]
	am.Description = mongoAm.Description
	return am
}
