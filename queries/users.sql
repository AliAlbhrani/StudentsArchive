-- name: CreateUser :one
INSERT INTO users (full_name, username, password, stage)
VALUES (@full_name, @username, @password, @stage)
RETURNING id;
