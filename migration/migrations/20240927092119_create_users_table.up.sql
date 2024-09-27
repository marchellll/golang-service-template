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