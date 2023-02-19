.SILENT: docker.build docker.run build run test-db

# prepares test database with dummy 1 million records
test-db:
	cd prepare-test-data && go mod download && go run main.go
# copy inside service so that it will used by the api
	cp prepare-test-data/test.db service/test.db

# builds docker image
docker.build:
	./scripts/build_image.sh

# runs service as a container
docker.run:
	docker run -p 8000:8000 gin-api

# build go binary
build:
	cd service && go mod download && \
	CGO_ENABLED=1 go build -o gin-api *.go

# runs service without docker
run: build
	cd service && ./gin-api
