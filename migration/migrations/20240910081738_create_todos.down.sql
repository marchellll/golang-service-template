-- reverse: set comment to table: "todos"
COMMENT ON TABLE "public"."todos" IS '';
-- reverse: create index "idx_deleted_at" to table: "todos"
DROP INDEX "public"."idx_deleted_at";
-- reverse: create index "idx_created_at" to table: "todos"
DROP INDEX "public"."idx_created_at";
-- reverse: create "todos" table
DROP TABLE "public"."todos";
