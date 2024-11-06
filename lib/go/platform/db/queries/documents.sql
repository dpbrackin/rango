-- name: AddDocument :one
INSERT INTO documents(
  user_id,
  source
) values ($1, $2)
RETURNING *;
--
-- name: GetDocument :one
SELECT * FROM documents WHERE id = $1;
