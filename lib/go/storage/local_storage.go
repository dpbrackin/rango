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

// Download implements core.DocumentBackend.
func (s *DiskStorage) Download(ctx context.Context, w io.Writer) error {
	panic("unimplemented")
}

// Upload implements core.DocumentBackend.
func (s *DiskStorage) Upload(ctx context.Context, params core.UploadParams) (path string, err error) {
	path = filepath.Join(s.BasePath, params.Name)

	file, err := os.Create(path)

	if err != nil {
		return "", fmt.Errorf("Failed to create file on disk: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, params.Reader)

	return path, err
}
