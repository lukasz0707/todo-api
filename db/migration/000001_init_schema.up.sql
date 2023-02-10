CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "is_blocked" boolean NOT NULL DEFAULT false
);

CREATE TABLE "todos" (
  "id" bigserial PRIMARY KEY,
  "group_id" bigint NOT NULL,
  "todo_name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "status" varchar NOT NULL DEFAULT 'ongoing',
  "deadline" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "groups" (
  "id" bigserial PRIMARY KEY,
  "group_name" varchar NOT NULL,
  "owner_id" bigint NOT NULL
);

CREATE TABLE "users_groups" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "group_id" bigint NOT NULL
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

CREATE INDEX ON "groups" ("owner_id");

CREATE INDEX ON "users_groups" ("user_id");

CREATE INDEX ON "users_groups" ("group_id");

CREATE UNIQUE INDEX ON "users_groups" ("user_id", "group_id");

CREATE INDEX ON "sessions" ("user_id");

COMMENT ON COLUMN "todos"."status" IS 'oneof(ongoing, suspended, completed)';

ALTER TABLE "todos" ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON DELETE CASCADE;

ALTER TABLE "groups" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "users_groups" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "users_groups" ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON DELETE CASCADE;

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;