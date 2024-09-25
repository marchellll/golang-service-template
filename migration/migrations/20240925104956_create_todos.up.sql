-- create "tasks" table
CREATE TABLE "public"."tasks" (
  "id" uuid NOT NULL,
  "description" text NOT NULL,
  "state" text NOT NULL,
  "created_by" uuid NOT NULL,
  "created_at" timestamptz NULL DEFAULT now(),
  "updated_at" timestamptz NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_created_at" to table: "tasks"
CREATE INDEX "idx_created_at" ON "public"."tasks" ("created_at");
-- create index "idx_deleted_at" to table: "tasks"
CREATE INDEX "idx_deleted_at" ON "public"."tasks" ("deleted_at");
-- set comment to table: "tasks"
COMMENT ON TABLE "public"."tasks" IS 'tasks table';
