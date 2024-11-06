package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"rango/core"
)

type DiskStorage struct {
	BasePath string
}

type NewDiscStorageParams struct {
	BasePath string
}

func NewDiscStorage(params NewDiscStorageParams) *DiskStorage {
	return &DiskStorage{
		BasePath: params.BasePath,
	}
}

func (s *DiskStorage) getAbsolutePath(name string) string {
	return filepath.Join(s.BasePath, name)
}

// Download implements core.DocumentBackend.
func (s *DiskStorage) Download(ctx context.Context, name string, w io.Writer) error {
	path := s.getAbsolutePath(name)

	f, err := os.Open(path)

	if err != nil {
		return err
	}

	_, err = io.Copy(w, f)

	if err != nil {
		return err
	}

	return nil
}

// Upload implements core.DocumentBackend.
func (s *DiskStorage) Upload(ctx context.Context, params core.UploadParams) (path string, err error) {
	path = s.getAbsolutePath(params.Name)

	file, err := os.Create(path)

	if err != nil {
		return "", fmt.Errorf("Failed to create file on disk: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, params.Reader)

	return path, err
}
