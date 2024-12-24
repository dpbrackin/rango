-- name: CreateOrganization :one
INSERT INTO organizations(id, name)
VALUES ($1, $2)
RETURNING *;
--
-- name: CreateMembership :one
INSERT INTO organization_memberships(user_id, org_id, is_default)
VALUES ($1, $2, $3)
RETURNING *;
--
-- name: GetOrganization :one
SELECT *
FROM organizations
WHERE id = $1;
--
-- name: GetUserOrganizations :many
SELECT org.*
FROM organizations org JOIN organization_memberships om ON org.id = om.org_id
WHERE om.user_id = $1;
--
-- name: GetDefaultUserOrganization :one
SELECT org.*
FROM organizations org JOIN organization_memberships om ON org.id = om.org_id
WHERE om.user_id = $1 AND is_default = true;
--
-- name: UpdateMembership :exec
UPDATE organization_memberships
SET is_default = $1
WHERE user_id = $2 and org_id = $3;
