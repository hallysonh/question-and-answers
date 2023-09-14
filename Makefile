VERSION ?= 0.0.1
NAME ?= question-and-answers-api
AUTHOR ?= Hallyson Almeida
REST_PORT ?= 3001
GRPC_PORT ?= 3002
NO_CACHE ?= true

.PHONY: build run stop clean

build:
	docker build -f cmd/api/Dockerfile -t $(NAME):$(VERSION) --no-cache=$(NO_CACHE) .

run:
	docker run --name $(NAME) -d \
	-p $(REST_PORT):$(REST_PORT) -p $(GRPC_PORT):$(GRPC_PORT) \
	--env-file .env -e DB_HOST=postgres --link $(NAME)-postgres-1 \
	--net $(NAME)_default $(NAME):$(VERSION) \
	&& docker ps -a --format "{{.ID}}\t{{.Names}}" |grep $(NAME)

stop:
	docker rm $$(docker stop $$(docker ps -a -q --filter "ancestor=$(NAME):$(VERSION)"))

clean:
	@rm -f main

DEFAULT: build