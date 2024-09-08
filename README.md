# Go DuckDB API
A sample containerized Go API backed by [DuckDB](https://duckdb.org/) database.

## About the service
The service exposes `GET http://localhost:8000/user/:id` api which returns user details DuckDB database `test.duckdb`.

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


## Running without docker
```bash
make run
```
This will build and run Go service without using docker.


## Running as a container
```bash
make docker.build
make docker.run
```

## Test data generation
This is an *optional* step as there is already `test.duckdb` duckdb file necessary to run the service without setting up anything. It contains a table called `users` which has following columns:

| id (int32)| joined_date (timestamp) | name (varchar)|    email (varchar)      |
|-----------|--------------------|---------------|-------------------------|
|      1    |     2021-09-14     |  Jarret Kuhn  |  carsondooley@wolf.name |


**Command:**
```bash
make seed
```
- Creates a duckdb database file `test.duckdb` inside folder `testdata`
- Then creates 'users' table and populates specified number of dummy records.
- `test.duckdb` file is copied to the docker image and used by the service.
