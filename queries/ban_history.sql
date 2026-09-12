-- name: GetBanHistory :many
-- Get ban history for a user
SELECT *
FROM ban_history
WHERE user_id = @user_id
;
