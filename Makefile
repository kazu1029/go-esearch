GOCMD=go
DOCKER_COMPOSE=docker-compose

.PHONY: start stop

start:
	$(DOCKER_COMPOSE) up

stop:
	$(DOCKER_COMPOSE) down

lint:
	golangci-lint run --disable-all --enable=goimports --enable=golint --enable=govet --enable=errcheck --enable=staticcheck
