# Go DuckDB API
A simple containerized Go API backed by [DuckDB](https://duckdb.org/).

## About
The service exposes `GET http://localhost:8000/users/:id` endpoint which returns user details from DuckDB running in [persistent mode](https://duckdb.org/docs/connect/overview.html#persistent-database).

**API response:**
```json
{
  "id": 1,
  "name": "Ken Adams",
  "email": "ken@dummy.org",
  "joined_date": "2021-09-12T05:13:37Z"
}
```
The service uses [go-duckdb](https://github.com/marcboeker/go-duckdb) pkg to interact with DuckDB C++ shared lib.

## Running the API

### 1. Generate Test Data
```bash
# populates 100K records in duckdb
make seed

# with custom seed size
SEED_COUNT=200000 make seed
```

The seed command generates `testdata/test.duckdb` duckdb file necessary to run the service. It contains `users` table which has following columns:

|               |                    |
|---------------|--------------------|
|    id         |     INTEGER        | 
|    name       |     VARCHAR        |
|    email      |     VARCHAR        |
|  joined_date  |     TIMESTAMP      |


### 2. Run Go Service

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

## Load testing
```bash
wrk -t12 -c100 -d10s --latency http://127.0.0.1:8000/users/100
```

_More about the load testing tool [wrk](https://github.com/wg/wrk)._
