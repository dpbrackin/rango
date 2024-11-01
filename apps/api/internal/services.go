package internal

import (
	"context"
	"fmt"
	"io"
	"rango/core"
)

type DocumentService struct {
	Backend    core.StorageBackend
	Repository core.DocumentRepository
}

type AddDocumentParams struct {
	User   core.User
	Name   string
	Reader io.Reader
}

func (s *DocumentService) CreateDocument(ctx context.Context, params AddDocumentParams) (*core.Document, error) {
	fileName := fmt.Sprintf("%s-%s", params.User.Username, params.Name)

	path, err := s.Backend.Upload(ctx, core.UploadParams{
		Reader: params.Reader,
		Name:   fileName,
	})

	if err != nil {
		return nil, err
	}

	doc := core.Document{
		Source: path,
		Owner:  params.User,
	}

	err = s.Repository.AddDocument(ctx, doc)

	if err != nil {
		return nil, err
	}

	return &doc, err
}
