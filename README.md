# Question And Answers

## Requirement Tools

* [Go](https://golang.org/dl/)
* [Docker](https://www.docker.com/)
* [Air (Live loading)](https://github.com/cosmtrek/air)

## Development requirements

```bash
# Protobuf compiler
brew install protobuf
# Protobuf Compiler to Go
go install google.golang.org/protobuf/protoc-gen-go@latest
# Protobuf GRPC Compiler to Go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
# Live loading tool
go install github.com/cosmtrek/air@latest
```

## Running locally

### Build and Run

```bash
go build -o tmp/api cmd/api/main.go && ./tmp/api
```

### With live loading using Air

```bash
air
```

### With Docker Compose

```bash
docker-compose up -d
```

### Using Docker with make

```bash
# Build api image
make build
# Start docker api container (database should be up)
make run
# Stop and remove api container
make stop
```

## Testing locally

### Running unit tests

```bash
go test question-and-answers/test
```

### Testing HTTP server

```bash
curl http://localhost:3001/api/question
```

> Change port to the HTTP Server Port configured

### Testing GRPC server

```bash
go run cmd/testing/grpctest/main.go --test=1
```

## Development Tips

### Compile Protobuf

```bash
./third_party/protoc-gen.sh
```

## TODO

- [ ] Create module to manager user registration
- [ ] Add user authentication and authorization
- [ ] The GRPC Error handling is exposing too much information. It should hide sensitive information when in production
- [ ] Add model validation when creating and updating data entries
- [ ] Use UUID instead of primary keys to expose models identity
- [ ] Add more filtering options
- [ ] Add list pagination support
