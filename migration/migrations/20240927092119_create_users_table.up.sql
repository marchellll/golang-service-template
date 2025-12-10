-- THESE IS NO SANE WAY TO MAKE DDLS THAT WORKS FOR BOTH POSTGRESQL AND MYSQL
-- comment/uncomment the code below to make it work for your database


-- THIS IS MIGRATION FOR POSTGRESQL DATABASE

CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "email" text NOT NULL,
  "password" text NOT NULL,
  "created_at" timestamptz NULL DEFAULT now(),
  "updated_at" timestamptz NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);

CREATE INDEX "idx_users_created_at" ON "public"."users" ("created_at");
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");






-- THIS IS MIGRATION FOR MYSQL DATABASE

-- CREATE TABLE users (
--   id varchar(36) NOT NULL,
--   email varchar(255) NOT NULL,
--   password text NOT NULL,
--   created_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
--   updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
--   deleted_at DATETIME NULL,
--   PRIMARY KEY (id)
-- );

-- CREATE INDEX idx_users_created_at ON users (created_at);
-- CREATE INDEX idx_users_deleted_at ON users (deleted_at);
-- CREATE UNIQUE INDEX idx_users_email ON users (email);