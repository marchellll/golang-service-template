# golang-service-template
A template to clone when building new service in Go

this should be a good starting point for building a new service in Go for small projects. It includes:
- http server
- db connection
- redis connection
- logging
- db migration

- Simple CRUD of Tasks
- Simple User Auth
- Simple RBAC


## DB migrations


### TLDR


```sh
# create a manual new migration
docker run -v ./migration/migrations:/migrations  --rm migrate/migrate create -ext sql -dir migrations create_users_table

# apply the migration
docker run -v .//migration/migrations:/migrations --network="host" migrate/migrate -path=/migrations/ -database "postgres://the_service_user:the_service_password@localhost:5432/the_service_database?sslmode=disable" up


# generate models and fluent query
gentool -dsn "host=localhost user=the_service_user password=the_service_password dbname=the_service_database port=5432 sslmode=disable" -outPath "./internal/dao/query"  -fieldNullable -fieldWithIndexTag -fieldWithTypeTag -fieldSignable -db postgres

```

### Long version


We can apply the migration using `go-migrate` (https://github.com/golang-migrate/migrate).


```sh
# create a manual new migration
# dont forget to add up and down sql scripts
docker run -v ./migration/migrations:/migrations  --rm migrate/migrate create -ext sql -dir migrations create_users_table

# apply the migration
# migrations must be tested in local before running anywhere else
# should use transaction
# `tcp(127.0.0.1:3306)`, the `tcp` is required, read here https://github.com/go-sql-driver/mysql/blob/af8d7931954ec21a96df9610a99c09c2887f2ee7/README.md#examples
docker run -v ./migration/migrations:/migrations --network="host" migrate/migrate -path=/migrations/ -database "postgres://the_service_user:the_service_password@localhost:5432/the_service_database?sslmode=disable" up


# check current version
docker run -v ./migration/migrations:/migrations --network="host" migrate/migrate -path=/migrations/ -database "postgres://the_service_user:the_service_password@localhost:5432/the_service_database?sslmode=disable" version

# rollback the migration
# this must always be tested in local before running anywhere else
docker run -v ./migration/migrations:/migrations --network="host" migrate/migrate -path=/migrations/ -database "postgres://the_service_user:the_service_password@localhost:5432/the_service_database?sslmode=disable" down 1




```


After that we use GORMS's GEN to generate the models and fluent query from the database. https://gorm.io/gen/gen_tool.html

```sh
go install gorm.io/gen/tools/gentool@latest

gentool -dsn "the_service_user:the_service_password@tcp(127.0.0.1:3306)/the_service_database" -outPath "./internal/dao/query"  -fieldNullable -fieldWithIndexTag -fieldWithTypeTag -fieldSignable -db mysql
gentool -dsn "host=localhost user=the_service_user password=the_service_password dbname=the_service_database port=5432 sslmode=disable" -outPath "./internal/dao/query"  -fieldNullable -fieldWithIndexTag -fieldWithTypeTag -fieldSignable -db postgres
```


## To run in docker

```sh
# dont forget the env file
cp .env.example .env
```

Make sure to change `DB_HOST` to `postgres` (and any other host) in `.env` file

```sh
DB_HOST=postgres # this one for docker-compose
# DB_HOST=localhost # this one for non-docker
```


```sh
# run all the dependencies and the service
docker compose --profile dev up

# run dependencies only, without the service
docker compose up

# build and run the service only
docker run -p 8080:8080 --env-file .env --rm -it $(docker build -q .)
```


## To run in local

Make sure to change `DB_HOST` to `localhost` (and any other host) in `.env` file

```sh
# dont forget the env file
cp .env.example .env
```


```bash
# run the service natively
go run ./cmd/server/main.go
```

## To build and run in local

```sh
# dont forget the env file
cp .env.example .env
```

```bash
# build and run the service natively
go build -o ./dist/run ./cmd/server
./dist/run
```


## To test the service

we can use Bruno to test the API

- https://www.usebruno.com/manifesto
- https://www.usebruno.com/blog/the-saas-dilemma


- run a collection: https://docs.usebruno.com/get-started/bruno-basics/run-a-collection

We can use the cli

```sh
npm install -g @usebruno/cli
```

Then
```sh
cd apitest
bru run --env local
```


## Ref

https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/


## TODO


- playwright
- ut

- terraform

- websocket
- grpc


- tracing
- metrics
- cron on kube

- github action