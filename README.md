# golang-service-template

A template to clone when building new service in Go

this should be a good starting point for building a new service in Go for small projects. It includes:

- http server
- [Dependency Injection](github.com/samber/do)
- db connection
- redis connection
- logging
- db migration

- Simple CRUD of Tasks
- Simple User Auth
- Simple rbac

## Pre-Commit

```sh

go fmt $(go list ./... | grep -v /vendor/)
go vet $(go list ./... | grep -v /vendor/)
go test -coverprofile=coverage.out $(go list ./... | grep -v /vendor/) ;    go tool cover -html=coverage.out

docker run -t --rm -v $(pwd):/app -v ~/.cache/golangci-lint/v1.62.2:/root/.cache -w /app golangci/golangci-lint:v1.62.2 golangci-lint run -v --timeout 10m

```

## DB migrations

### TLDR

```sh
# create a manual new migration
docker run -v ./migration/migrations:/migrations  --rm migrate/migrate create -ext sql -dir migrations create_users_table

# apply the migration
docker run --rm -v ./migration/migrations:/migrations --network="host" migrate/migrate -path=/migrations/ -database "postgres://the_service_user:the_service_password@localhost:5432/the_service_database?sslmode=disable" up

## or for mysql

docker run --rm -v ./migration/migrations:/migrations --network="host" migrate/migrate -path=/migrations/ -database "mysql://the_service_user:the_service_password@tcp(localhost:3306)/the_service_database?charset=utf8mb4&parseTime=True" up


# generate models and fluent query
# go install gorm.io/gen/tools/gentool@latest
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

## Add DAO Model from Database

### Install Gorm Gentool
https://gorm.io/gen/gen_tool.html#Install

### Execute command

```
gentool -c gentool.yaml
```

see [./gentool.yaml](./gentool.yaml)
1. Change DB dialect user, password and database name
3. Change tables params to the table name

## Quick Crud Generator

After Making Migration and Generating Models, we can use the `sergen` to generate the CRUD for the model.

```sh
# run from the root of the project
go run ./cmd/sergen -ModuleName "golang-service-template" -EntityName Goose -EntityNamePlural Geese
```

That command will generate the CRUD for the `Goose` model. It will generate the following:

- `internal/service/goose.go`
- `internal/handler/goose.go`
- add endpoints to `internal/app/routes.go`
- register dependensy in `internal/app/di.go`

Of course, we can manually create the CRUD, but this is a good starting point.

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

## Unit Test

install [ginkgo](https://onsi.github.io/ginkgo/#getting-started). A BDD testing framework for Go.

```sh
go get github.com/onsi/ginkgo/v2/ginkgo
go install github.com/onsi/ginkgo/v2/ginkgo
go get github.com/onsi/gomega/...
```

you can run the test using `go test ./...` like usual or using ginkgo

```sh
go test ./...`

# or
ginkgo  ./...
```

to bootstrap the package to use ginko, you can use `ginkgo bootstrap`. (Only need to run once per package)

```sh
# cd to your package
cd repository
# bootstrap the package
ginkgo bootstrap
```

to generate test file, you can use `ginkgo generate [file_name]`

```sh
# cd to your package
cd repository
# generate test file
ginkgo generate loyalty_card
```

## Generate Mock

We are using [vektra.github.io/mockery](https://vektra.github.io/mockery/latest/) to generate mock files for testing using [github.com/stretchr/testify](https://github.com/stretchr/testify).

Install mockery by running `brew install mockery`

After creating a new interface  (a new repository/service), run `mockery` to generate the mock file for the new interfaces.

[github.com/stretchr/testify](https://github.com/stretchr/testify) main advantage over the other mocking is allowing expectation using `mock.Anything`


## Integration Test

we can use Bruno to test the API

- <https://www.usebruno.com/manifesto>
- <https://www.usebruno.com/blog/the-saas-dilemma>

- run a collection: <https://docs.usebruno.com/get-started/bruno-basics/run-a-collection>

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

<https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/>

## TODO

- ut

- socket.io
- temporal
- queue

- tracing
- metrics
- cron on kube

- playwright
- grpc
- terraform
- github action
