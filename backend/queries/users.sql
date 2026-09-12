-- name: CreateUser :one
-- Create a new user
INSERT INTO users (full_name, username, password, stage)
VALUES (@full_name, @username, @password, @stage)
RETURNING id;

-- name: CreateUserProfile :exec
-- Create a user profile
INSERT INTO profiles (user_id, photo_url, bio)
VALUES (@user_id, @photo_url, @bio)
ON CONFLICT (user_id) DO UPDATE SET photo_url = @photo_url, bio = @bio;

-- name: GetUserProfile :one
-- Get user profile by user id
SELECT u.full_name, u.username, p.photo_url, p.bio, u.created_at, u.stage
FROM users u
LEFT JOIN profiles p ON p.user_id = u.id
WHERE u.id = @user_id
AND deleted_at IS NULL
;

-- name: GetUserByUsername :one
-- Get user id and password by username for login
SELECT id, password, banned
FROM users
WHERE username = @username
AND deleted_at IS NULL;
