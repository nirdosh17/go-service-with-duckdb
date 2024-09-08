FROM golang:1.22 AS builder

LABEL maintainer="Nirdosh Gautam <nrdshgtm@gmail.com>"

ARG CPU_ARCH
ARG DUCKDB_VERSION=v1.0.0

RUN apt-get update; \
  apt-get -y install unzip

WORKDIR /service

RUN wget -nv https://github.com/duckdb/duckdb/releases/download/${DUCKDB_VERSION}/libduckdb-linux-${CPU_ARCH}.zip -O libduckdb.zip; \
  unzip libduckdb.zip -d /tmp/libduckdb

# copying go.mod and go.sum first to leverage docker cache
COPY service/go.* /service/

RUN go mod download

COPY service/ .

RUN CGO_ENABLED=1 \
  CGO_LDFLAGS="-L/tmp/libduckdb" \
  go build -tags=duckdb_use_lib  \
  -o duckapi *.go

FROM ubuntu

WORKDIR /service/

# ~90 MB
COPY ./testdata/test.duckdb* /service

COPY --from=builder /service/duckapi .
COPY --from=builder /tmp/libduckdb/libduckdb.so ./libduckdb/libduckdb.so

ENV DUCK_DB_FILE_PATH=/service/test.duckdb
ENV LD_LIBRARY_PATH=./libduckdb

EXPOSE 8000

CMD ["./duckapi"]
