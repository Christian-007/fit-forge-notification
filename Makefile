include .env

test:
	sh -c 'env $$(cat .env | xargs) go test ./...'

run:
	sh -c 'env $$(cat .env | xargs) go run ./cmd'

run_prod:
	sh -c 'env $$(cat .env.prod | xargs) go run ./cmd'

.PHONY: test run run_prod
