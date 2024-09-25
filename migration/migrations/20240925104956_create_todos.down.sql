-- reverse: set comment to table: "tasks"
COMMENT ON TABLE "public"."tasks" IS '';
-- reverse: create index "idx_deleted_at" to table: "tasks"
DROP INDEX "public"."idx_deleted_at";
-- reverse: create index "idx_created_at" to table: "tasks"
DROP INDEX "public"."idx_created_at";
-- reverse: create "tasks" table
DROP TABLE "public"."tasks";
