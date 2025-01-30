package services

import (
	"context"
	"fmt"
	"rango/api/internal/core"
	"rango/api/internal/db/generated"
	"rango/api/internal/db/repositories"
	"rango/api/internal/search"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

type SearchService struct {
	Engine          core.SearchEngine
	IndexRepository core.IndexRepository
}

type CreateIndexParams struct {
	Name   string
	Engine string
	Org    core.Organization
}

func (s *SearchService) CreateIndex(ctx context.Context, params CreateIndexParams) (core.Index, error) {
	idx := core.Index{
		ID:     core.IDType(uuid.New()),
		Name:   params.Name,
		Engine: params.Engine,
		Org:    params.Org,
	}

	idx, err := s.IndexRepository.CreateIndex(ctx, idx)

	if err != nil {
		return core.Index{}, err
	}

	err = s.Engine.CreateIndex(ctx, core.CreateSearchIndexParams{
		Index: idx,
	})

	if err != nil {
		return core.Index{}, err
	}

	return idx, nil
}

type NewAISearchServiceParams struct {
	DB generated.DBTX
}

func NewAISearchService(params NewAISearchServiceParams) *SearchService {
	embedder := search.NewOpenAIEmbedder()

	qdrantClient, err := qdrant.NewClient(&qdrant.Config{Host: "localhost", Port: 6334})

	if err != nil {
		panic("failed to initialize qdrant client")
	}

	vectorStore := search.NewQdrantVectorStore(search.NewQdrantVectorStoreParams{
		Client:     qdrantClient,
		VectorSize: 1536,
	})

	engine := &search.AISearchEngine{
		Embedder:    embedder,
		VectorStore: vectorStore,
	}

	queries := generated.New(params.DB)

	repository := repositories.NewPGIndexRepository(queries)

	return &SearchService{
		Engine:          engine,
		IndexRepository: repository,
	}
}

type NewSearchServiceParams struct {
	DB     generated.DBTX
	Engine search.EngineType
}

func NewSearchService(params NewSearchServiceParams) (*SearchService, error) {
	if params.Engine == search.AISearchEngineType {
		return NewAISearchService(NewAISearchServiceParams{
			DB: params.DB,
		}), nil
	}

	return nil, fmt.Errorf("Unsupported engine")
}
