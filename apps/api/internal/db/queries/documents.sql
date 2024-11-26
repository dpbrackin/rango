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
SELECT documents.*, users.username
FROM documents JOIN users on users.id = documents.user_id
WHERE documents.id = $1;
--
-- name: UpdateDocument :exec
UPDATE documents
SET user_id = $2, source = $3, content = $4, type = $5
where id = $1;
