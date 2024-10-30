IMAGE_NAME=forum-image

DOCKERFILE=Dockerfile.dev

CONTAINER_NAME=forum-container

build:
	sudo docker build -f $(DOCKERFILE) -t $(IMAGE_NAME) .

run:
	sudo docker run --name $(CONTAINER_NAME) -p 8080:8080 -v $(PWD):/app $(IMAGE_NAME)

stop:
	sudo docker stop $(CONTAINER_NAME) || true
	sudo docker rm $(CONTAINER_NAME) || true

clean:
	sudo rm -rf forum tmp || true
	sudo docker rm $(CONTAINER_NAME) || true
	sudo docker rmi $(IMAGE_NAME) || true

up: stop clean build run

.PHONY: build run stop clean up
