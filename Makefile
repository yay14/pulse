# Define variables
APP_NAME=pulse
IMAGE_NAME=$(APP_NAME):latest
CONTAINER_NAME=container

# Default target
all: build

# Build the Docker image
build:
	docker build -t $(IMAGE_NAME) .

# Run the application in a container
run: build
	docker run --name $(CONTAINER_NAME) -p 8080:8080 $(IMAGE_NAME)

# Stop and remove the running container
stop:
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

# Remove the Docker image
clean:
	docker rmi $(IMAGE_NAME) || true

# Remove all stopped containers and unused images
cleanup:
	docker system prune -f
