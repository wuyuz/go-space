

CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "email" varchar DEFAULT NULL,
  "pwd" varchar NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "staff" boolean DEFAULT true
);

CREATE TABLE "items" (
  "id" SERIAL PRIMARY KEY,
  "content" bigint,
  "user_id" int
);

ALTER TABLE "items" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

CREATE UNIQUE INDEX ON "users" ("id", "username", "pwd");

CREATE UNIQUE INDEX ON "items" ("id", "user_id");