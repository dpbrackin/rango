-- name: AddDocument :exec
INSERT INTO documents(
  user_id,
  source
) values ($1, $2);
--
-- name: GetDocument :one
SELECT * FROM documents WHERE id = $1;
