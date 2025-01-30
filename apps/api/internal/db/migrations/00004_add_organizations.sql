-- +goose Up
-- +goose StatementBegin
CREATE TABLE organizations(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL
);

CREATE TABLE organization_memberships(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  org_id UUID REFERENCES organizations(id) ON DELETE RESTRICT,
  is_default BOOLEAN
);

CREATE UNIQUE INDEX unique_default_org_per_user
ON organization_memberships (user_id)
WHERE is_default;

CREATE UNIQUE INDEX unique_user_per_org
ON organization_memberships (user_id, org_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS organizations;
DROP TABLE IF EXISTS organization_memberships;
-- +goose StatementEnd
