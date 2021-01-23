ifeq "${count}" ""
	count=1
endif

ifeq "${run}" ""
	run=""
endif

# `make test count=50` to run `go test -race -count=50 ./...`
# `make test run=TestXXX` to run `go test -race -run=TestXXX ./...`
test:
	go test -race -run=${run} -count=${count} ./...

.PHONY: test

nats-test:
	docker-compose -f .docker/nats-test.yml up --build --abort-on-container-exit --remove-orphans; \
	docker-compose -f .docker/nats-test.yml down

.PHONY: nats-test

mongo-test:
	docker-compose -f .docker/mongostore-test.yml up --build --abort-on-container-exit --remove-orphans; \
	docker-compose -f .docker/mongostore-test.yml down

.PHONY: mongo-test

coverage:
	docker-compose -f .docker/coverage.yml up --build --abort-on-container-exit --remove-orphans; \
	docker-compose -f .docker/coverage.yml down; \
	go tool cover -html=out/coverage.out

.PHONY: coverage

bench:
	go test -v -bench=. -run=${run} -count=${count} ./...

.PHONY: bench
