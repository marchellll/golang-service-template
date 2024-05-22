# golang-service-template
A template to clone when building new service in Go


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

- docker
- redis
- postgres
- logging


- playwright
- ut

- terraform

- websocket
- grpc


- tracing
- metrics
- cron on kube