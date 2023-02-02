CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "todos" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "group" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "status" varchar NOT NULL DEFAULT 'ongoing',
  "deadline" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "users_todos" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "todos_id" bigint NOT NULL,
  "has_permissions" bool NOT NULL DEFAULT false
);

CREATE INDEX ON "todos" ("created_at");

CREATE INDEX ON "todos" ("deadline");

CREATE INDEX ON "todos" ("status");

CREATE INDEX ON "todos" ("name");

CREATE INDEX ON "users_todos" ("user_id");

CREATE INDEX ON "users_todos" ("todos_id");

CREATE UNIQUE INDEX ON "users_todos" ("user_id", "todos_id");

COMMENT ON COLUMN "todos"."status" IS 'oneof(ongoing, suspended, completed)';

ALTER TABLE "users_todos" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "users_todos" ADD FOREIGN KEY ("todos_id") REFERENCES "todos" ("id");