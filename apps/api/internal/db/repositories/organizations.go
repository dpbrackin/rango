package repositories

import (
	"context"
	"rango/api/internal/core"
	"rango/api/internal/db/generated"

	"github.com/jackc/pgx/v5/pgtype"
)

type PGOrganizationRepository struct {
	queries *generated.Queries
}

func NewPGOrganizationRepository(q *generated.Queries) *PGOrganizationRepository {
	return &PGOrganizationRepository{
		queries: q,
	}
}

// CreateMembership implements core.OrgatizationRepository.
func (p *PGOrganizationRepository) CreateMembership(ctx context.Context, params core.CreateMembershipParams) error {
	panic("unimplemented")
}

// CreateOrganization implements core.OrgatizationRepository.
func (p *PGOrganizationRepository) CreateOrganization(ctx context.Context, org core.Organization) (core.Organization, error) {
	o, err := p.queries.CreateOrganization(ctx, generated.CreateOrganizationParams{
		ID:   pgtype.UUID{Valid: false},
		Name: "",
	})

	if err != nil {
		return core.Organization{}, err
	}

	return core.Organization{
		ID:   o.ID.Bytes,
		Name: o.Name,
	}, nil
}

// GetOrganization implements core.OrgatizationRepository.
func (p *PGOrganizationRepository) GetOrganization(ctx context.Context, id core.IDType) (core.Organization, error) {
	panic("unimplemented")
}
