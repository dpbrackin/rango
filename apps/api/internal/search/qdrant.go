package search

import (
	"context"
	"fmt"
	"log/slog"
	"rango/api/internal/core"

	"github.com/qdrant/go-client/qdrant"
)

type QDrantVectoreStore struct {
	client     *qdrant.Client
	vectorSize int32
}

// CreateIndex implements VectorStore.
func (q *QDrantVectoreStore) CreateIndex(ctx context.Context, idx core.Index) error {
	indexName, err := idx.ID.MarshalText()

	if err != nil {
		panic("failed to unmarshall uuid")
	}

	return q.client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: string(indexName),
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(q.vectorSize),
			Distance: qdrant.Distance_Cosine,
		}),
	})
}

// Retrieve implements VectorStore.
func (q *QDrantVectoreStore) Retrieve(ctx context.Context, params RetrieveParams) ([]RetrievalResult, error) {
	indexId, err := params.Index.ID.MarshalText()

	if err != nil {
		panic("failed to unmarshall uuid")
	}

	embedding := convertEmbeddingFloat32(params.Embedding)

	res, err := q.client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: string(indexId),
		Query:          qdrant.NewQuery(embedding...),
	})

	if err != nil {
		return []RetrievalResult{}, fmt.Errorf("Failed to query: %w", err)
	}

	fmt.Println(res)

	return []RetrievalResult{}, nil
}

// Store implements VectorStore.
func (q *QDrantVectoreStore) Store(ctx context.Context, params StoreEmbeddingsParams) error {
	vectors := make([]*qdrant.PointStruct, len(params.Embeddings))

	for i, e := range params.Embeddings {
		float32Vector := convertEmbeddingFloat32(e)
		vectors[i] = &qdrant.PointStruct{
			Vectors: qdrant.NewVectors(float32Vector...),
			Payload: qdrant.NewValueMap(map[string]any{
				"documentId": params.Doc.ID.String(),
			}),
		}
	}

	info, err := q.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: params.Index.ID.String(),
	})

	if err != nil {
		return fmt.Errorf("Failed to store document: %w", err)
	}

	slog.Debug("storage info", slog.String("info", info.String()))

	return nil
}

func convertEmbeddingFloat32(e Embedding) []float32 {
	vector := make([]float32, len(e))

	for i, v := range e {
		vector[i] = float32(v)
	}

	return vector
}

var _ VectorStore = &QDrantVectoreStore{}

type NewQdrantVectorStoreParams struct {
	Client     *qdrant.Client
	VectorSize int32
}

func NewQdrantVectorStore(params NewQdrantVectorStoreParams) *QDrantVectoreStore {
	return &QDrantVectoreStore{
		client:     params.Client,
		vectorSize: params.VectorSize,
	}
}
