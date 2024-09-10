schema "public" {
  comment = "standard public schema"
}

table "todos" {
    schema = schema.public

    comment = "todos table"

    column "id" {
      type = bigserial
    }

    column "text" {
      type = text
    }

    column "created_at" {
      type = timestamptz
      null = true
      default = sql("now()")
    }

    column "updated_at" {
      type = timestamptz
      null = true
      default = sql("now()")
    }

    column "deleted_at" {
      type = timestamptz
      null = true
      default = sql("null")
    }

    primary_key {
      columns = [
        column.id
      ]
    }

    index "idx_deleted_at" {
      columns = [column.deleted_at]
    }

    index "idx_created_at" {
      columns = [column.created_at]
    }
}
