package core

import "context"

type User struct {
	ID       IDType
	Username string
	Org      Organization
}

type Organization struct {
	ID   IDType
	Name string
}

type CreateMembershipParams struct {
	User User
	Org  Organization
}

type OrgatizationRepository interface {
	CreateOrganization(ctx context.Context, org Organization) (Organization, error)
	GetOrganization(ctx context.Context, id IDType) (Organization, error)
	CreateMembership(ctx context.Context, params CreateMembershipParams) error
}
