.SILENT: docker.build docker.run build run prepare-test-db

prepare-test-db:
# prepare test database
	cd prepare-test-data && go mod download && go run main.go
# copy inside service so that it will used by the api
	cp prepare-test-data/test_data.db service/test_data.db

docker.build:
	./scripts/build_image.sh

# run with docker
docker.run: docker.build
	docker run -p 8000:8000 gin-api

build:
	cd service && go mod download && \
	CGO_ENABLED=1 go build -o gin-api *.go

# run without docker
run: build
	cd service && ./gin-api
