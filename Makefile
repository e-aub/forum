IMAGE_NAME=forum-image

DOCKERFILE=Dockerfile.dev

CONTAINER_NAME=forum-container

build:
	sudo docker build -f $(DOCKERFILE) -t $(IMAGE_NAME) .

run:
	sudo docker run --name $(CONTAINER_NAME) --network host -v $(PWD):/app $(IMAGE_NAME)

stop:
	sudo docker stop $(CONTAINER_NAME) || true
	sudo docker rm $(CONTAINER_NAME) || true

clean:
	sudo rm -rf forum tmp || true
	sudo docker rm $(CONTAINER_NAME) || true
	sudo docker rmi $(IMAGE_NAME) || true
push: clean
	@read -p "Enter commit message: " msg; \
	git add .; \
	git commit -m "$$msg"; \
	git push
all: stop clean build run

.PHONY: build run stop clean up
