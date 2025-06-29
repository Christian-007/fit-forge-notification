## Prerequisites

Before getting started, ensure you have done the following things:

1. Install `Go`
2. Clone [Core API repository](https://github.com/Christian-007/fit-forge)
3. Execute `make generate_jwks` from the Core API repository
4. Copy and paste `private.pem` file to this repository

## Getting Started

To start this repo on your machine, do the following:

1.  Clone this repo
2.  Go to the repo directory on your machine
3.  Execute `go mod tidy && go mod vendor`
4.  Setup the environment variables in `./.env` file (see below for details)
5.  Finally, `make run` to run the Go app with `.env` file

## Environment Variables

Setup the following environment variables in `./.env` file:

```
ENV=localhost
REDIS_DSN=localhost:6380
REDIS_PASSWORD=secret
PUBSUB_PROJECT_ID=local-project
PUBSUB_EMULATOR_HOST=localhost:8085
FRONTEND_URL=http://localhost:3000
```
