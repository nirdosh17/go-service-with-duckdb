FROM golang:1.18 as builder

LABEL maintainer="Nirdosh Gautam <nrdshgtm@gmail.com>"

RUN apt-get update; \
  apt-get -y install unzip

WORKDIR /service

RUN wget -nv https://github.com/duckdb/duckdb/releases/download/v0.7.0/libduckdb-linux-amd64.zip -O libduckdb.zip; \
  unzip libduckdb.zip -d /tmp/libduckdb

COPY service/ .

RUN go mod download

RUN CGO_ENABLED=1 \
  CGO_LDFLAGS="-L/tmp/libduckdb" \
  go build -tags=duckdb_use_lib  \
  -o gin-api *.go

FROM debian:10.9-slim

WORKDIR /service/

# ~25 MB
COPY ./service/test.db .

COPY --from=builder /service/gin-api .
COPY --from=builder /tmp/libduckdb/libduckdb.so ./libduckdb/libduckdb.so

ENV LD_LIBRARY_PATH ./libduckdb

EXPOSE 8000

CMD ["./gin-api"]
