# golang-service-template
A template to clone when building new service in Go


## DB migrations


We use https://github.com/golang-migrate/migrate

```sh

# create a new migration
# dont forget to add up and down sql scripts
docker run -v ./migrations:/migrations  --rm migrate/migrate create -ext sql -dir migrations create_users_table

# apply the migration
# migrations must be tested in local before running anywhere else
# should use transaction
# `tcp(127.0.0.1:3306)`, the `tcp` is required, read here https://github.com/go-sql-driver/mysql/blob/af8d7931954ec21a96df9610a99c09c2887f2ee7/README.md#examples
docker run -v ./migrations:/migrations --network="host" migrate/migrate -path=/migrations/ -database "mysql://the_service_user:the_service_password@tcp(127.0.0.1:3306)/the_service_database" up

# rollback the migration
# this must always be tested in local before running anywhere else
docker run -v ./migrations:/migrations --network="host" migrate/migrate -path=/migrations/ -database "mysql://the_service_user:the_service_password@tcp(127.0.0.1:3306)/the_service_database" down 1
```


After that we use GORMS's GEN to generate the models and fluent query from the database. https://gorm.io/gen/gen_tool.html

```sh
go install gorm.io/gen/tools/gentool@latest

gentool -dsn "the_service_user:the_service_password@tcp(127.0.0.1:3306)/the_service_database" -outPath "./internal/dao/query"  -fieldNullable -fieldWithIndexTag -fieldWithTypeTag -withUnitTest -fieldSignable -db mysql
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

docker compose up --build the-service
```


TODO: add docker file
TODO: add github-action


## Ref

https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/


## TODO


- db migration✅

- redis ✅
- mysql: https://xata.io/pricing ✅
- logging ✅


- playwright
- ut

- terraform

- websocket
- grpc


- tracing
- metrics
- cron on kube

- github action