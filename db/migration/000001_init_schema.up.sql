CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "role" varchar NOT NULL DEFAULT 'user',
  "password_changed_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "is_blocked" boolean NOT NULL DEFAULT false
);

CREATE TABLE "todos" (
  "id" bigserial PRIMARY KEY,
  "group_id" bigint NOT NULL,
  "todo_name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "status" varchar NOT NULL DEFAULT 'todo',
  "deadline" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "groups" (
  "id" bigserial PRIMARY KEY,
  "group_name" varchar NOT NULL
);

CREATE TABLE "users_groups" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "group_id" bigint NOT NULL,
  "role" varchar NOT NULL DEFAULT 'user'
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "todos" ("group_id");

CREATE INDEX ON "todos" ("created_at");

CREATE INDEX ON "todos" ("status");

CREATE INDEX ON "todos" ("deadline");

CREATE INDEX ON "users_groups" ("user_id");

CREATE INDEX ON "users_groups" ("group_id");

CREATE UNIQUE INDEX ON "users_groups" ("user_id", "group_id");

CREATE INDEX ON "sessions" ("user_id");

COMMENT ON COLUMN "users"."role" IS 'oneof(user, moderator, admin)';

COMMENT ON COLUMN "todos"."status" IS 'oneof(todo, suspended, completed)';

COMMENT ON COLUMN "users_groups"."role" IS 'oneof(user, owner)';

ALTER TABLE "todos" ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON DELETE CASCADE;

ALTER TABLE "users_groups" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "users_groups" ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON DELETE CASCADE;

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;