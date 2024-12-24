package core

import (
	"context"
	"io"
	"time"
)

type Clock interface {
	Now() time.Time
}

type Embeddings [][]float64

func Core(name string) string {
	result := "Core " + name
	return result
}

// Document struct represents a document that can be indexed and used as a source for RAG
type Document struct {
	ID      IDType
	Source  string // path to the document's underlying file
	Owner   User
	Content io.Reader
	Type    string
}

// Index is a collection of documents.
type Index struct {
	ID   IDType
	Name string
}

type DocumentRepository interface {
	AddDocument(ctx context.Context, d Document) (Document, error)
	UpdateDocument(ctx context.Context, d Document) error
	GetDocument(ctx context.Context, id IDType) (Document, error)
}

type UploadParams struct {
	Reader io.Reader
	Name   string
}

type StorageBackend interface {
	// Upload saves the contents of the reader into persistent storage.
	//
	// A successful upload returns a path which can be used to retrieve the contents and err == nil.
	Upload(ctx context.Context, params UploadParams) (path string, err error)
	Download(ctx context.Context, name string, w io.Writer) error
}

// Embedder creates embeddings for an input
type Embedder interface {
	Embed(ctx context.Context, content io.Reader) (Embeddings, error)
}

// ContentExtractor extracts the text content from a document
type ContentExtractor interface {
	// Extract extracts the document's content and writes the content to w
	Extract(ctx context.Context, doc Document, w io.Writer) error
}

// VectorStore stores embeddings for a document and searches based on embeddings
type VectorStore interface {
	Store(ctx context.Context, params StoreEmbeddingsParams) error
	Retrieve(ctx context.Context, embedding Embeddings) ([]RetrievalResult, error)
}

type StoreEmbeddingsParams struct {
	Doc        Document
	Embeddings Embeddings
	Index      Index
	Text       string
}

type RetrievalResult struct {
	Docs Document
}
