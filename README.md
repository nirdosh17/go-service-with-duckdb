# Go DuckDB API
A simple containerized Go API backed by [DuckDB](https://duckdb.org/) database.

## About
The service exposes `GET http://localhost:8000/users/:id` api which returns user details DuckDB database `test.duckdb`.

API response:
```json
{
  "id": 100000,
  "name": "Maximillian Flatley",
  "email": "darronkoepp@pouros.org",
  "joined_date": "2021-09-12T05:13:37Z"
}
```
The service uses [go-duckdb](https://github.com/marcboeker/go-duckdb) library to interact with DuckDB C++ shared library.

## Running the API

### 1. Generating test data
```bash
# populates 100K records
make seed

# with custom seed size
SEED_COUNT=200000 make seed
```

The seed command generates `testdata/test.duckdb` duckdb file necessary to run the service. It contains `users` table which has following columns:

| id (int32)| joined_date (timestamp) | name (varchar)|    email (varchar)      |
|-----------|--------------------|---------------|-------------------------|
|      1    |     2021-09-14     |  Jarret Kuhn  |  carsondooley@wolf.name |


### 2. Running service

Normally:
```bash
make build
make run
```

As a docker container:
```bash
make docker.build
make docker.run
```

### Load testing
```bash
wrk -t12 -c100 -d10s --latency http://127.0.0.1:8000/users/100
```

_More about the load testing tool [wrk](https://github.com/wg/wrk)._
