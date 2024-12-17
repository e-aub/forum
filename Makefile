IMAGE_NAME=forum-image

DOCKERFILE=Dockerfile

CONTAINER_NAME=forum-container

build:
	  docker build -f $(DOCKERFILE) -t $(IMAGE_NAME) .

run:
	  docker run --name $(CONTAINER_NAME) -v $(PWD):/app $(IMAGE_NAME)

stop:
	  docker stop $(CONTAINER_NAME) || true
	  docker rm $(CONTAINER_NAME) || true

clean:
	  rm -rf forum tmp || true
	 docker rm -f $(CONTAINER_NAME) || true
	 docker rmi -f $(IMAGE_NAME) || true
push: clean
	@read -p "Enter commit message: " msg; \
	git add .; \
	git commit -m "$$msg"; \
	git push
all: stop clean build run

.PHONY: build run stop clean up

test:
	PORT=8080 DB_PATH=db/data.db go run cmd/main.go