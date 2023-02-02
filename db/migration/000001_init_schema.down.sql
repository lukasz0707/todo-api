ALTER TABLE IF EXISTS "users_todos" DROP CONSTRAINT IF EXISTS "users_todos_user_id_fkey";
ALTER TABLE IF EXISTS "users_todos" DROP CONSTRAINT IF EXISTS "users_todos_todos_id_fkey";

DROP TABLE IF EXISTS users_todos;
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS users;