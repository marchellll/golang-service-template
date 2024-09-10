# Build the application from source
FROM golang:alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum *.go ./
COPY cmd ./cmd
COPY internal ./internal
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /the-service  ./cmd/server

# Run the tests in the container
# FROM build-stage AS run-test-stage
# RUN go test -v ./...

# Deploy the application binary into a lean image
FROM scratch AS build-release-stage

WORKDIR /

COPY --from=build-stage /the-service /the-service

EXPOSE 8080

ENTRYPOINT ["/the-service"]