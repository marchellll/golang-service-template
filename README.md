# golang-service-template
A template to clone when building new service in Go


## DB migrations

we use Atlas migration tool to manage db migrations

https://atlasgo.io/cli-reference#atlas-schema-inspect

This command below will inspect the schema of the database and generate a schema file in HCL format. This command should be sun in in the initiation of the project and whenever there is a change in the schema by other methods (ORM/SQL scripts) ouside of Atlas

```bash
atlas schema inspect --format "{{ sql . }}" -u "postgres://root:password@localhost:3306/db" > db/schema.sql
```

To create changes in the db schema, we can edit the schema.hcl file and then run the following command to generate the migration file
https://atlasgo.io/versioned/diff

```bash

# check migration status, should not be run on the initial setup
atlas migrate status --url "postgres://root:password@localhost:3306/db"  --dir "file://db/migrations"



# append, add user table to schema.sql
# schema.sql is the FINAL intended schema of our DB. It should be updated manually
# Atlas recommended using the HCL format, but I found it easier to use SQL format.
# Using SQL format also decouples the schema from the Atlas tool, no vendor lock in

echo "
CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  created_at DATETIME DEFAULT NOW(),
  updated_at DATETIME DEFAULT NOW() ON UPDATE NOW(),
  deleted_at DATETIME
);
" >> db/schema.sql



# generate new migration file. no changes to DB yet
# this will look up the schema file and the current schema in the DB and generate a migration file
# Atlas just jere to help us to generate these migration files, we can create them manually if we want
# If we were to leave Atlas, we can still use the migration files to apply the changes manually to the DB
atlas migrate diff --dev-url "postgres://root:password@localhost:3306/db" --dir file://db/migrations --to "file://db/schema.sql" add_users_table



# apply the migration, MUST use dry run first
atlas migrate apply --dry-run --dir "file://db/migrations" --url "postgres://root:password@localhost:3306/db" 1



# WARNING: this will apply the migration to the DB without confirmation
atlas migrate apply --dir "file://db/migrations" --url "postgres://root:password@localhost:3306/db" 1
```


## To run in local

```bash
cp env.example .env
go run main.go
```
## To build and run in local

```bash
go build -o ./dist/run
./dist/run
```

## To run in docker

```
docker run -p 8080:8080 -e PORT=8080 --rm -it $(docker build -q .)
```


TODO: add docker file
TODO: add github-action


## Ref

https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/


## TODO


- db migration✅

- redis
- postgres: https://xata.io/pricing
- logging


- playwright
- ut

- terraform

- websocket
- grpc


- tracing
- metrics
- cron on kube

- github action