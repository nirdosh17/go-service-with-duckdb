# Containerized Go Service with DuckDB
Sample Containerized Go service using [DuckDB](https://duckdb.org/)

## About the service
The service exposes `GET http://localhost:8000/stats` api which returns aggregated user count from DuckDB database `test.duckdb`.

API response:
```json
[
	{
		"date": "2021-02-20",
		"users_joined": 2598
	},
	{
		"date": "2021-02-21",
		"users_joined": 2578
	}
]
```
The service uses [go-duckdb](https://github.com/marcboeker/go-duckdb) library to interact with DuckDB C libraries.


## Running without docker
```bash
make run
```
This will build and run the GIN service without using docker.


## Running as a container
```bash
# builds docker image downloading DuckDB dependencies
make docker.build

# runs docker image
make docker.run
```

## Test data generation
This is an *optional* step as there is already `test.duckdb` duckdb file necessary to run the service without setting up anything. It contains a table called `users` which has following columns:

| id (int32)| joined_date (date) | name (varchar)|    email (varchar)      |
|-----------|--------------------|---------------|-------------------------|
|      1    |     2021-09-14     |  Jarret Kuhn  |  carsondooley@wolf.name |


**Command:**
```bash
make test-db
```
- Creates a duckdb database file `test.duckdb` inside folder `prepare-test-data`
- Then creates 'users' table and populates 1 million dummy data. Takes around 2 mins.
- `test.duckdb` file is copied to the docker image and used by the service.
