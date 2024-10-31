IMAGE_NAME=forum-image

DOCKERFILE=Dockerfile.dev

CONTAINER_NAME=forum-container

build:
	docker build -f $(DOCKERFILE) -t $(IMAGE_NAME) .

run:
	docker run --name $(CONTAINER_NAME) -p 8080:8080 -v $(PWD):/app $(IMAGE_NAME)

stop:
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

clean:
	rm -rf forum tmp || true
	docker rm $(CONTAINER_NAME) || true
	docker rmi $(IMAGE_NAME) || true

all: stop clean build run

.PHONY: build run stop clean up
