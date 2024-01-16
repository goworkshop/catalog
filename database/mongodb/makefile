.PHONY: run-mongo cleanup-mongo

run-mongo:
	docker rm -f mongo || true
	docker run -d -p 27017:27017 --name mongo mongo:latest

cleanup-mongo:
	docker rm -f mongo || true

