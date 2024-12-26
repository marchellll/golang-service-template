-- THESE IS NO SANE WAY TO MAKE DDLS THAT WORKS FOR BOTH POSTGRESQL AND MYSQL
-- comment/uncomment the code below to make it work for your database


-- THIS IS MIGRATION FOR POSTGRESQL DATABASE

-- create "tasks" table
-- CREATE TABLE "public"."tasks" (
--   "id" uuid NOT NULL,
--   "description" text NOT NULL,
--   "state" text NOT NULL,
--   "created_by" uuid NOT NULL,
--   "created_at" timestamptz NULL DEFAULT now(),
--   "updated_at" timestamptz NULL DEFAULT now(),
--   "deleted_at" timestamptz NULL,
-- );
-- CREATE INDEX "idx_created_at" ON "public"."tasks" ("created_at");
-- CREATE INDEX "idx_deleted_at" ON "public"."tasks" ("deleted_at");
-- COMMENT ON TABLE "tasks" IS 'tasks table';



-- THIS IS MIGRATION FOR MYSQL DATABASE

CREATE TABLE tasks (
  id varchar(36) NOT NULL,
  description text NOT NULL,
  state text NOT NULL,
  created_by varchar(36) NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  PRIMARY KEY (id)
);
CREATE INDEX idx_created_at ON tasks (created_at);
CREATE INDEX idx_deleted_at ON tasks (deleted_at);