package repositories

import (
	"context"
	"rango/api/internal/core"
	"rango/api/internal/db/generated"

	"github.com/jackc/pgx/v5/pgtype"
)

type PGIndexRepository struct {
	queries *generated.Queries
}

// CreateIndex implements core.IndexRepository.
func (p *PGIndexRepository) CreateIndex(ctx context.Context, i core.Index) (core.Index, error) {
	idx, err := p.queries.CreateIndex(ctx, generated.CreateIndexParams{
		ID: pgtype.UUID{
			Bytes: i.ID,
			Valid: true,
		},
		OrgID: pgtype.UUID{
			Bytes: i.Org.ID,
			Valid: true,
		},
		Name:   i.Name,
		Engine: i.Engine,
	})

	return core.Index{
		ID:     core.IDType(idx.ID.Bytes),
		Name:   idx.Name,
		Engine: idx.Engine,
		Org:    i.Org,
	}, err
}

// GetIndex implements core.IndexRepository.
func (p *PGIndexRepository) GetIndex(ctx context.Context, id core.IDType) (core.Index, error) {
	idx, err := p.queries.GetIndex(ctx, pgtype.UUID{
		Bytes: id,
		Valid: true,
	})

	return core.Index{
		ID:     core.IDType(idx.ID.Bytes),
		Name:   idx.Name,
		Engine: idx.Engine,
	}, err
}

// ListDocuments implements core.IndexRepository.
func (*PGIndexRepository) ListDocuments(ctx context.Context, i core.Index, p core.Pagination) ([]core.Document, error) {
	panic("unimplemented")
}

func NewPGIndexRepository(queries *generated.Queries) *PGIndexRepository {
	return &PGIndexRepository{
		queries: queries,
	}
}

var _ core.IndexRepository = &PGIndexRepository{}
