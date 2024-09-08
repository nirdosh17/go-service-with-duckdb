.DEFAULT_GOAL=help

.SILENT: docker.build docker.run build run test-db

seed: ## generate test database with dummy records
	rm -f testdata/test.duckdb*
	cd testdata && go mod download && go run main.go

docker.build: ## build docker image
	./scripts/build_image.sh

docker.run: ## run service inside docker container
	docker run -p 8000:8000 duck-api

build: ## build go binary
	cd service && go mod download && \
	CGO_ENABLED=1 go build -o duckapi *.go

run: ## run service without docker
	cd service && ./duckapi

help:
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ": "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
