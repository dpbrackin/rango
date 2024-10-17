-- name: GetUserByUsername :one
SELECT
  *
FROM
  users
WHERE
  username = $1;
--
-- name: AddUser :exec
INSERT INTO users(username, password) VALUES($1,$2);
--
-- name: GetSession :one
SELECT sessions.*, users.username
FROM sessions JOIN users on users.id = sessions.user_id
WHERE sessions.id = $1;
