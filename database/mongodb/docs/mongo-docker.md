# MongoDB on Docker

Follow these steps to run MongoDB in a Docker container.

## Prerequisites

- Docker installed on your machine

## Installation Steps

1. Run the MongoDB container:

```bash
# Remove existing container if any
docker rm -f mongo

# Run MongoDB container
docker run -d -p 27017:27017 --name mongo mongo:latest
```

## MongoDB Container

The MongoDB container is running at `localhost:27017` with the following credentials:

- Username: root
- Password: example

## Cleanup (Docker)

To stop and remove the MongoDB container, run:

```bash
docker rm -f mongo
```
