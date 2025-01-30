-- name: CreateIndex :one
INSERT INTO indexes(
  id,
  org_id,
  name,
  engine
) VALUES ($1, $2, $3, $4)
  RETURNING *;
--
-- name: GetIndex :one
SELECT *
FROM indexes
WHERE id = $1;
--
-- name: UpdateIndex :exec
UPDATE indexes
SET name = $1
WHERE id = $2;
