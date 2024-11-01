package repositories

import (
	"context"
	"rango/core"
	"rango/db/generated"

	"github.com/jackc/pgx/v5/pgtype"
)

type PGDocumentsRepository struct {
	queries *generated.Queries
}

func NewPGDocumentRepository(queries *generated.Queries) *PGDocumentsRepository {
	return &PGDocumentsRepository{
		queries: queries,
	}
}

// AddDocument implements core.DocumentRepository.
func (p *PGDocumentsRepository) AddDocument(ctx context.Context, d core.Document) error {
	err := p.queries.AddDocument(ctx, generated.AddDocumentParams{
		UserID: pgtype.Int4{
			Int32: int32(d.Owner.ID),
			Valid: true,
		},
		Source: d.Source,
	})

	return err
}
