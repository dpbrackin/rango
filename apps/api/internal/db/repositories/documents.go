package repositories

import (
	"context"
	"fmt"
	"io"
	"rango/api/internal/db/generated"
	"rango/core"

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
func (p *PGDocumentsRepository) AddDocument(ctx context.Context, d core.Document) (core.Document, error) {
	doc, err := p.queries.AddDocument(ctx, generated.AddDocumentParams{
		UserID: pgtype.Int4{
			Int32: int32(d.Owner.ID),
			Valid: true,
		},
		Source: d.Source,
		Type:   d.Type,
	})

	document := core.Document{
		ID:      core.IDType(doc.ID),
		Source:  d.Source,
		Owner:   d.Owner,
		Content: d.Content,
		Type:    d.Type,
	}

	return document, err
}

func (p *PGDocumentsRepository) UpdateDocument(ctx context.Context, d core.Document) error {
	contentBytes, err := io.ReadAll(d.Content)

	if err != nil {
		fmt.Errorf("Failed to update document: %w", err)
	}

	err = p.queries.UpdateDocument(ctx, generated.UpdateDocumentParams{
		ID:      int32(d.ID),
		UserID:  pgtype.Int4{Int32: int32(d.Owner.ID), Valid: true},
		Source:  d.Source,
		Content: pgtype.Text{String: string(contentBytes), Valid: true},
		Type:    d.Type,
	})

	return err
}
