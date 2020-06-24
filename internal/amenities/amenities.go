package amenities

import (
	"context"
	"encoding/hex"
	"github.com/alexpashkov/asched/graph/model"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"os"
	"path/filepath"
)

type Service struct {
	logger      *logrus.Entry
	mongoClient *mongo.Client
	mongoDBName string
	photosDir   string
}

func NewService(logger *logrus.Entry, mongoClient *mongo.Client, mongoDBName, photosDir string) *Service {
	return &Service{
		logger:      logger,
		mongoClient: mongoClient,
		mongoDBName: mongoDBName,
		photosDir:   photosDir,
	}
}

func (s *Service) Start(ctx context.Context) error {
	cur, err := s.mongoCollection().Indexes().List(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list indexes")
	}
	defer cur.Close(ctx)
	var indexes []bson.M
	if err := cur.All(ctx, &indexes); err != nil {
		return err
	}
	if !s.hasIndexes(indexes, "location_2dsphere") {
		s.logger.Info("creating 2dsphere index")
		if _, err := s.mongoCollection().Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.M{"location": "2dsphere"},
		}); err != nil {
			return errors.Wrap(err, "failed to create an index")
		}
	} else {
		s.logger.Info("found all required indexes")
	}
	return nil
}

func (s *Service) hasIndexes(indexes []bson.M, names ...string) bool {
	namesMap := make(map[string]bool)
	for _, index := range indexes {
		indexName, ok := index["name"].(string)
		if ok {
			namesMap[indexName] = true
			s.logger.Debugf("found %s index", indexName)
		} else {
			s.logger.Debugf("couldn't cast index name to string, got %T", index["name"])
		}
	}
	for _, name := range names {
		if !namesMap[name] {
			return false
		}
	}
	return true
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
			_, err := s.AddPhoto(id, newAm.Photo.File)
			return errors.Wrap(err, "failed to save the photo")
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
	dirPath := filepath.Join(s.photosDir, id)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", errors.Wrap(err, "failed to create a directory")
	}
	dst, err := os.Create(filepath.Join(dirPath, photoID))
	if err != nil {
		return photoID, errors.Wrap(err, "failed to create a file")
	}
	_, err = io.Copy(dst, file)
	return photoID, errors.Wrap(err, "failed to write to the file")
}

func (s *Service) GetPhotoIDs(id string) ([]string, error) {
	var res []string
	err := filepath.Walk(filepath.Join(s.photosDir, id),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			res = append(res, info.Name())
			return nil
		},
	)
	if err != nil && err != os.ErrNotExist {
		return nil, err
	}
	return res, nil
}

func (s *Service) mongoCollection() *mongo.Collection {
	return s.mongoClient.Database(s.mongoDBName).Collection("amenities")
}

// mongoAmenity is a structure used to persist Amenity in a MongoDB database
type mongoAmenity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Type        []string           `bson:"type"`
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
