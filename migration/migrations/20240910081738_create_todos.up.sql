-- create "todos" table
CREATE TABLE "public"."todos" (
  "id" bigserial NOT NULL,
  "text" text NOT NULL,
  "created_at" timestamptz NULL DEFAULT now(),
  "updated_at" timestamptz NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_created_at" to table: "todos"
CREATE INDEX "idx_created_at" ON "public"."todos" ("created_at");
-- create index "idx_deleted_at" to table: "todos"
CREATE INDEX "idx_deleted_at" ON "public"."todos" ("deleted_at");
-- set comment to table: "todos"
COMMENT ON TABLE "public"."todos" IS 'todos table';
