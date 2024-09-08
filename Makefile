.DEFAULT_GOAL=help

.SILENT: docker.build docker.run build run test-db

testdata: ## prepare test database with dummy 1 million records
	cd prepare-test-data && go mod download && go run main.go

docker.build: ## build docker image
	./scripts/build_image.sh

docker.run: ## run service in docker container
	docker run -p 8000:8000 gin-api

build: ## build go binary
	cd service && go mod download && \
	CGO_ENABLED=1 go build -o gin-api *.go

run: ## run service
	cd service && ./gin-api

help:
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ": "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
