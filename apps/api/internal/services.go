package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"rango/core"
	"rango/platform/eventbus"
	"rango/platform/extractors"
)

type DocumentService struct {
	Storage    core.StorageBackend
	Repository core.DocumentRepository
	EventBus   *eventbus.EventBus
}

type AddDocumentParams struct {
	User   core.User
	Name   string
	Reader io.Reader
}

type IndexService struct {
	Embeder core.Embedder
}

type IndexDocumentParams struct {
	Document core.Document
}

func (s *DocumentService) CreateDocument(ctx context.Context, params AddDocumentParams) (*core.Document, error) {
	fileName := fmt.Sprintf("%s-%s", params.User.Username, params.Name)

	path, err := s.Storage.Upload(ctx, core.UploadParams{
		Reader: params.Reader,
		Name:   fileName,
	})

	if err != nil {
		return nil, err
	}

	doc := core.Document{
		Source: path,
		Owner:  params.User,
		Type:   filepath.Ext(path)[1:],
	}

	doc, err = s.Repository.AddDocument(ctx, doc)

	if err != nil {
		return nil, err
	}

	s.EventBus.Publish(eventbus.Event{
		Topic: DOCUMENT_CREATED_TOPIC,
		Data:  doc,
	})

	return &doc, err
}

// ExtractContent extracts the text content from the document and persists it
func (s *DocumentService) ExtractContent(ctx context.Context, doc *core.Document) error {
	var extractor core.ContentExtractor

	if doc.Type == "txt" {
		extractor = &extractors.TextExtractor{
			Storage: s.Storage,
		}
	}

	if extractor == nil {
		return errors.New("Failed to create an extractor")
	}

	content := bytes.NewBuffer(nil)
	extractor.Extract(ctx, *doc, content)
	contentReader := bytes.NewReader(content.Bytes())

	doc.Content = contentReader

	err := s.Repository.UpdateDocument(ctx, *doc)

	// Allow the content to be read again
	contentReader.Seek(0, 0)

	if err != nil {
		return fmt.Errorf("Failed to update document: %w", err)
	}

	return nil
}

func (s *IndexService) IndexDocument(ctx context.Context, params IndexDocumentParams) error {
	embeddings, err := s.Embeder.Embed(ctx, params.Document.Content)

	if err != nil {
		return fmt.Errorf("Failed to get embeddings: %w", err)
	}

	fmt.Print(embeddings)

	return nil
}
