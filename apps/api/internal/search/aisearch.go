package search

import (
	"context"
	"io"
	"rango/api/internal/core"
)

type AISearchEngine struct {
	Embedder    Embedder
	VectorStore VectorStore
}

// AddToIndex implements core.SearchEngine.
func (a *AISearchEngine) AddToIndex(ctx context.Context, params core.AddToIndexParams) error {
	panic("unimplemented")
}

// CreateIndex implements core.SearchEngine.
func (a *AISearchEngine) CreateIndex(ctx context.Context, params core.CreateSearchIndexParams) error {
	return a.VectorStore.CreateIndex(ctx, params.Index)
}

// Search implements core.SearchEngine.
func (a *AISearchEngine) Search(ctx context.Context, params core.SearchParams) (core.SearchResult, error) {
	panic("unimplemented")
}

type Embedding []float64

// Embedder creates embeddings for an input
type Embedder interface {
	Embed(ctx context.Context, content io.Reader) ([]Embedding, error)
}

// VectorStore stores embeddings for a document and searches based on embeddings
type VectorStore interface {
	Store(ctx context.Context, params StoreEmbeddingsParams) error
	Retrieve(ctx context.Context, params RetrieveParams) ([]RetrievalResult, error)
	CreateIndex(ctx context.Context, idx core.Index) error
}

type StoreEmbeddingsParams struct {
	Doc        core.Document
	Embeddings []Embedding
	Index      core.Index
	Text       string
}

type RetrievalResult struct {
	Docs core.Document
}

type RetrieveParams struct {
	Index     core.Index
	Embedding Embedding
}
