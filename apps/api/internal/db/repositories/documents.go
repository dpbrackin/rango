package repositories

import (
	"context"
	"fmt"
	"io"
	"rango/api/internal/core"
	"rango/api/internal/db/generated"
	"strings"

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
		UserID: pgtype.UUID{
			Bytes: d.ID,
			Valid: true,
		},
		Source: d.Source,
		Type:   d.Type,
	})

	document := core.Document{
		ID:      core.IDType(doc.ID.Bytes),
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
		ID: pgtype.UUID{
			Bytes: d.ID,
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: d.Owner.ID,
			Valid: true,
		},
		Source:  d.Source,
		Content: pgtype.Text{String: string(contentBytes), Valid: true},
		Type:    d.Type,
	})

	return err
}

func (p *PGDocumentsRepository) GetDocument(ctx context.Context, Id core.IDType) (core.Document, error) {
	doc, err := p.queries.GetDocument(ctx, pgtype.UUID{
		Bytes: Id,
		Valid: true,
	})

	if err != nil {
		return core.Document{}, fmt.Errorf("Failed to get document from DB: %w", err)
	}

	return core.Document{
		ID:     Id,
		Source: doc.Source,
		Owner: core.User{
			ID:       core.IDType(doc.UserID.Bytes),
			Username: doc.Username,
		},
		Content: strings.NewReader(doc.Content.String),
		Type:    doc.Type,
	}, nil
}
