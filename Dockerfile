FROM golang:1.18 as builder

LABEL maintainer="Nirdosh Gautam <nrdshgtm@gmail.com>"

ARG CPU_ARCH

RUN apt-get update; \
  apt-get -y install unzip

WORKDIR /service

RUN wget -nv https://github.com/duckdb/duckdb/releases/download/v0.7.0/libduckdb-linux-${CPU_ARCH}.zip -O libduckdb.zip; \
  unzip libduckdb.zip -d /tmp/libduckdb

# copying go.mod and go.sum first to leverage docker cache
COPY service/go.* /service/

RUN go mod download

COPY service/ .

RUN CGO_ENABLED=1 \
  CGO_LDFLAGS="-L/tmp/libduckdb" \
  go build -tags=duckdb_use_lib  \
  -o gin-api *.go

FROM debian:10.9-slim

WORKDIR /service/

# ~90 MB
COPY ./prepare-test-data/test.duckdb /prepare-test-data/test.duckdb

COPY --from=builder /service/gin-api .
COPY --from=builder /tmp/libduckdb/libduckdb.so ./libduckdb/libduckdb.so

ENV LD_LIBRARY_PATH ./libduckdb

EXPOSE 8000

CMD ["./gin-api"]
