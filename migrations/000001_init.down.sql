BEGIN;

DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS students;
DROP INDEX IF EXISTS idx_posts_title;
DROP INDEX IF EXISTS idx_posts_content;
DROP INDEX IF EXISTS idx_posts_user_id;

COMMIT;
