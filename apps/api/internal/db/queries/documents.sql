-- name: AddDocument :one
INSERT INTO documents(
  user_id,
  source,
  content,
  type
  ) values ($1, $2, $3, $4)
RETURNING *;
--
-- name: GetDocument :one
SELECT * FROM documents WHERE id = $1;
--
-- name: UpdateDocument :exec
UPDATE documents
SET user_id = $2, source = $3, content = $4, type = $5
where id = $1;
