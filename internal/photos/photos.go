package photos

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

type Service struct {
	rootDir string
}

func NewService(rootDir string) *Service {
	return &Service{
		rootDir: rootDir,
	}
}

func (s *Service) SavePhoto(groupID string, photo io.Reader) error {
	photoID := uuid.New().String()
	if photoID == "" {
		return errors.New("generated empty uuid")
	}
	dst, err := os.Create(filepath.Join(s.rootDir, groupID, photoID))
	if err != nil {
		return errors.Wrap(err, "failed to create a file")
	}
	_, err = io.Copy(dst, photo)
	return errors.Wrap(err, "failed to write to the file")
}

func (s *Service) GetPhotoIDsByGroupID(groupID string) ([]string, error) {
	var res []string
	return res, filepath.Walk(
		filepath.Join(s.rootDir, groupID),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			res = append(res, info.Name())
			return nil
		},
	)
}
