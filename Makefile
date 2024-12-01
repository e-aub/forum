IMAGE_NAME=forum-image

DOCKERFILE=Dockerfile.dev

CONTAINER_NAME=forum-container


PORT=8080
APP_ENV=local
DB_PATH=db/data.db

build:
	 docker build -f $(DOCKERFILE) -t $(IMAGE_NAME) .

run:
	 docker run --name $(CONTAINER_NAME) --network host -v $(PWD):/app $(IMAGE_NAME)

stop:
	 docker stop $(CONTAINER_NAME) || true
	 docker rm $(CONTAINER_NAME) || true

clean:
	 rm -rf forum tmp || true
	 docker rm $(CONTAINER_NAME) || true
	 docker rmi $(IMAGE_NAME) || true
push: clean
	@read -p "Enter commit message: " msg; \
	git add .; \
	git commit -m "$$msg"; \
	git push
all: stop clean build run

.PHONY: build run stop clean up

# Define the environment variable
test:

test:
	@echo "Running tests with DB_URL=$${DB_URL}"
	PORT=8080 DB_PATH=db/data.db go run cmd/main.go
