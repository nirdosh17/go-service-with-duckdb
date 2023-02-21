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
# builds docker image downloading DuckDB's C dependencies.
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


## CPU Profiling
1. **CPU profiles for a single request**

    Command:
    ```bash
    bombardier -c 1 -n 1 http://localhost:8000/stats
    ```

    Normal mode:
    ```bash
    Showing top 10 nodes out of 23
      flat  flat%   sum%        cum   cum%
      40ms 57.14% 57.14%       40ms 57.14%  <unknown>
      30ms 42.86%   100%       30ms 42.86%  runtime.cgocall
    ```

    Running the service as container:
    ```bash
    (pprof) top
    Showing nodes accounting for 1.14s, 95.00% of 1.20s total
    Showing top 10 nodes out of 94
          flat  flat%   sum%        cum   cum%
        0.67s 55.83% 55.83%      0.67s 55.83%  [libduckdb.so]
        0.31s 25.83% 81.67%      0.31s 25.83%  runtime.cgocall
        0.09s  7.50% 89.17%      0.09s  7.50%  [libc-2.28.so]
    ```

2. **CPU profile for 10 requests with concurrency 2**

    Command: `bombardier -c 2 -n 10 http://localhost:8000/stats`

    Normal mode:
    ```bash
    (pprof) top
    Showing nodes accounting for 540ms, 100% of 540ms total
    Showing top 10 nodes out of 23
          flat  flat%   sum%        cum   cum%
        300ms 55.56% 55.56%      300ms 55.56%  <unknown>
        240ms 44.44%   100%      240ms 44.44%  runtime.cgocall
            0     0%   100%      240ms 44.44%  database/sql.(*DB).Query (inline)
    ```

    Running the service as container:
    ```bash
    (pprof) top
    Showing nodes accounting for 9.04s, 97.94% of 9.23s total
    Dropped 86 nodes (cum <= 0.05s)
    Showing top 10 nodes out of 26
          flat  flat%   sum%        cum   cum%
        5.66s 61.32% 61.32%      5.66s 61.32%  [libduckdb.so]
        2.59s 28.06% 89.38%      2.60s 28.17%  runtime.cgocall
        0.61s  6.61% 95.99%      0.61s  6.61%  [libc-2.28.so]
        0.17s  1.84% 97.83%      0.17s  1.84%  [libpthread-2.28.so]
    ```

    ### CPU Profiling Impressions
    go-duckdb lib is taking longer time in CGO calls as seen in the results. Even for a single request, CGO call inside container is ~10-16X times more.
