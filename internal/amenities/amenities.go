package amenities

import (
	"context"
	"encoding/hex"
	"github.com/alexpashkov/asched/graph/model"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"os"
	"path/filepath"
)

type Service struct {
	mongoClient *mongo.Client
	mongoDBName string
	photosDir   string
}

func NewService(mongoClient *mongo.Client, mongoDBName, photosDir string) *Service {
	return &Service{
		mongoClient: mongoClient,
		mongoDBName: mongoDBName,
		photosDir:   photosDir,
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
			return errors.Wrap(s.AddPhoto(id, newAm.Photo.File), "failed to save the photo")
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

func (s *Service) DeleteAmenity(ctx context.Context, id string) error {
	var mongoID primitive.ObjectID
	if _, err := hex.Decode(mongoID[:], []byte(id)); err != nil {
		return err
	}
	_, err := s.mongoCollection().DeleteOne(ctx, bson.M{"_id": mongoID})
	return err
}

func (s *Service) AddPhotos(ctx context.Context, id string, files ...io.Reader) ([]string, error) {
	IDs := make([]string, len(files))
	for _, file := range files {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if id, err := s.AddPhoto(id, file); err != nil {
				return nil, err
			} else {
				IDs = append(IDs, id)
			}
		}
	}
	return IDs, nil
}

func (s *Service) AddPhoto(id string, file io.Reader) (string, error) {
	photoID := uuid.New().String()
	if photoID == "" {
		return photoID, errors.New("generated empty uuid")
	}
	dst, err := os.Create(filepath.Join(s.photosDir, id, photoID))
	if err != nil {
		return photoID, errors.Wrap(err, "failed to create a file")
	}
	_, err = io.Copy(dst, file)
	return photoID, errors.Wrap(err, "failed to write to the file")
}

func (s *Service) GetPhotoIDs(id string) ([]string, error) {
	var res []string
	return res, filepath.Walk(
		filepath.Join(s.photosDir, id),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			res = append(res, info.Name())
			return nil
		},
	)
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
