schema "public" {
  comment = "standard public schema"
}

table "tasks" {
    schema = schema.public

    comment = "tasks table"

    column "id" {
      type = uuid
    }

    column "description" {
      type = text
    }

    column "state" {
      type = text
    }

    column "created_by" {
      type = timestamptz
      null = true
      default = sql("now()")
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
