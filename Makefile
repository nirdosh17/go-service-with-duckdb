.DEFAULT_GOAL=help

.SILENT: docker.build docker.run build run test-db

seed: ## 1. generate test data
	rm -f testdata/test.duckdb*
	cd testdata && go mod download && go run seed.go

docker.build: ## 2. build docker image
	./scripts/build_image.sh

docker.run: ## 3. run service inside docker container
	docker run --rm --name duck-api -p 8000:8000 duck-api

build: ## build go binary
	cd service && go mod download && \
	CGO_ENABLED=1 go build -o duckapi *.go

run: ## run service without docker
	cd service && \
	DUCK_DB_FILE_PATH=../testdata/test.duckdb \
	./duckapi

help:
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ": "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
