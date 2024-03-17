# echo-test

REST API built with Go, Echo, PostgreSQL, Pgx, go-migrate and SQLC

## Requirements

- Go 1.22+
- Git
- Docker

## What is this about?

The goal here is to build tests first whenever possible so that
the API design and driven by tests. 

The `httpexpect` package is being used for testing the REST API.
It is a really cool package and has a realy cool API that makes
testing HTTP and REST api cool.

The tests take advantage of `TestMain` which resets the database
before the tests run. `go-migrate` is used to drop the schema
and then generate the schema whenever tests are run.

Go caches the tests but that is not desired here so ensure
cache is not enabled I have added a `make test` command
which ensures the code runs in test mode `APP_ENV=test` and
cache is disabled `-count=1`. The `test` command:

I was struggling to write http related tests in Go cuz I 
didn't like the API but things became clearer to me when 
I found out about `httpexpect` and saw the
[echo test](https://github.com/gavv/httpexpect/blob/master/_examples/echo_test.go) example.

The example from the repo as you see is really simple but it is enough
to get started. I have added a bunch of things on top of it like
`testutils` functions which reset the test database and create a testing user.
I have also refactored the example code to use newer code from the echo
documentation.

```sh
test:
    APP_ENV=test go test ./... -count=1
```

## Load testing locally

I installed [`wrk`](https://github.com/wg/wrk) to benchmark my API locally. That is how I did:

1. Build the binary and run it in production mode

```sh
make build-api
make start
```

2. Seed the database

```sh
make db-seed
```

3. Login to generate a jwt token

```sh
## TODO: use application/json instead. change handler to decode JSON from body
curl -X POST http://localhost:8080/login \
   -H "Content-Type: application/x-www-form-urlencoded" \
   -d "email=paulo@example.com&password=password12" 

# Result
{"token": "..."}
```

4. Copy the generated token and export an env variable

```sh
export JWT="<copied token>"
```

5. Run wrk against a protected route which will verify the token on every request

```sh
# Using 15 threads and 1000 concurrent connections during 1 minute
wrk -t15 -c1000 -d60s -H 'Content-Type: application/json' \
-H "Authorization: Bearer ${JWT}" \
http://localhost:8080/restricted/hello
```

6. Result I got. My specs: M2 Pro 16GB RAM

```sh
Running 1m test @ http://localhost:8080/restricted/hello
  15 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     4.50ms    4.22ms  62.80ms   71.86%
    Req/Sec     3.49k     1.64k    8.67k    70.98%
  3125402 requests in 1.00m, 482.86MB read
  Socket errors: connect 755, read 100, write 0, timeout 0
Requests/sec:  52053.16
Transfer/sec:      8.04MB
```

> I'm impressed! That is really cool. No performance tunning has been made.
