# Batch Transactions API

The current system allows us to create a transaction and issue it to the bank immediately. In the context of introducing transactions that are created as a result of a round-up, there will be potentially a lot of transactions worth a few cents issued to the bank.

In order to avoid spamming the user's bank account with a lot of low volume transactions, we want to "batch" those transactions into a single higher volume transactions.

## Setup

### Requirements

 * [Go Toolchain](https://golang.org/doc/install) versions 1.15+

### From Source

```shell
$ make start
// OR
$ go run main.go -n 1000 -d 24h
```

### Docker

Build and run from local Dockerfile:

```shell
$ docker-compose up
```

## Arguments
```shell
-m int
    minimal value of transactions to be batched (default 100).
-t duration
    batch timeout duration (default 0).
```

E.g.:
```shell
$ go run main.go -m 100000

$ go run main.go -m 10000 -d 0

$ go run main.go -m 10000 -d 24h

$ go run main.go -d 20s
```

_If you define the duration, a worker going to execute to dispatch all batch transaction if reach the duration limit_

## Docs

Swagger API docs provided at path `/swagger/index.html`

or you can install `go-swagger` and render it locally (macOS example)

Install:

```shell
$ brew tap go-swagger/go-swagger
$ brew install go-swagger
```

Render: 

```shell
$ swagger serve docs/swagger.yaml
```

Re-generate swag docs:

```shell
$ make swag
```

## Useful Commands

### Unit tests
```shell
$ make test
```

### Help

```shell
$ make help

 Choose a command run in batch:

  install            Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
  start              Clean, compile and start simulation.
  start-simulation   Start alian simulation from binary.
  stop               Stop the simulation.
  compile            Compile the project.
  exec               Run given command. e.g; make exec run="go test ./..."
  clean              Clean build files. Runs `go clean` internally.
  check              Run application check.
  test               Run all tests.
  unit               Run all unit tests.
  fmt                Run `go fmt` for all go files.
  govet              Run go vet.
  golint             Run golint.
  install-swag       Install go-swagger.
  swag               Install and run go-swagger.
```
