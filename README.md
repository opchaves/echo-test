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
