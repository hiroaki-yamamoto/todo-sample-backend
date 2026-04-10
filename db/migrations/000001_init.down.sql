DROP INDEX IF EXISTS idx_todos_user_id ON todos;
DROP INDEX IF EXISTS idx_todos_completed_at ON todos;
DROP INDEX IF EXISTS idx_todos_wip_at ON todos;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS todos;
