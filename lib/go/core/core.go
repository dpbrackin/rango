package core

import (
	"context"
	"io"
	"time"
)

type IDType int32

type Clock interface {
	Now() time.Time
}

func Core(name string) string {
	result := "Core " + name
	return result
}

// Document struct represents a document that can be indexed and used as a source for RAG
type Document struct {
	ID     IDType
	Source string // path to the document's underlying file
	Owner  User
}

type DocumentRepository interface {
	AddDocument(ctx context.Context, d Document) (Document, error)
}

type UploadParams struct {
	Reader io.Reader
	Name   string
}

type StorageBackend interface {
	// Upload saves the contents of the reader into persistent storage.
	//
	// A successful upload returns a path which can be used to retrive the contents and err == nil.
	Upload(ctx context.Context, params UploadParams) (path string, err error)
	Download(ctx context.Context, name string, w io.Writer) error
}

// Index is takes a document and indexes it
type Indexer interface {
	Index(ctx context.Context, d Document) error
}
