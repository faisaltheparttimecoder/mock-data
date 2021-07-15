IMAGE = "mock:latest"
PGDATABASE = "testmockdb"
PGUSER = "postgres"
PGHOST = "localhost"
PGPORT = 5432
PGPASSWORD = "postgres"
MOCK_DATA_TEST_RUNNER = "true"

# h - help
h help:
	@echo "h, help 			- Help menu"
	@echo "docker 				- Run docker image build"
	@echo "tests				- Run all the test suites (Unit & Integration)"
	@echo "unit_tests 			- Run mock data unit test suites"
	@echo "integration_tests 		- Run mock data integration test suites"
	@echo "build				- Build mock data binary"
	@echo "install_postgres		- Install postgres locally via brew"
	@echo "start_postgres		- Start postgres locally via brew"
	@echo "stop_postgres		- Stop postgres locally via brew"
	@echo "uninstall_postgres		- Uninstall postgres locally via brew"
	@echo "create_db			- create a mock database for testing"
	@echo "recreate_db			- drop if exists and create a mock database for testing"
	@echo "drop_db				- drop the mock database"
.PHONY: h

# docker build
docker:
	docker build -f ./build/Dockerfile -t $(IMAGE) .
.PHONY: docker

# Run Mock data test Suite
unit_tests:
	make recreate_db
	@echo "#### Running Mock Data Test"
	MOCK_DATA_TEST_RUNNER=$(MOCK_DATA_TEST_RUNNER) PGPASSWORD=$(PGPASSWORD) PGDATABASE=$(PGDATABASE) PGUSER=$(PGUSER) PGHOST=$(PGHOST) PGPORT=$(PGPORT) GOFLAGS="-count=1" go test -v . -race -cover
	@echo "#### Finished Mock Data Test"
	make drop_db
.PHONY: unit_tests

# Run Mock data loading test for all supported datatype
integration_tests:
	make recreate_db
	MOCK_DATA_TEST_RUNNER=$(MOCK_DATA_TEST_RUNNER) PGPASSWORD=$(PGPASSWORD) PGDATABASE=$(PGDATABASE) PGUSER=$(PGUSER) PGHOST=$(PGHOST) PGPORT=$(PGPORT) /bin/bash integration_test.sh
	make drop_db
.PHONY: integration_tests

# Run Mock data test codecov
codecov_tests:
	make recreate_db
	@echo "#### Running Mock Data Test"
	MOCK_DATA_TEST_RUNNER=$(MOCK_DATA_TEST_RUNNER) PGPASSWORD=$(PGPASSWORD) PGDATABASE=$(PGDATABASE) PGUSER=$(PGUSER) PGHOST=$(PGHOST) PGPORT=$(PGPORT) GOFLAGS="-count=1" go test -v . -coverprofile=coverage.txt -covermode=atomic
	@echo "#### Finished Mock Data Test"
	make drop_db
.PHONY: unit_tests

# Run all mock data test
tests:
	make unit_tests
	make integration_tests
.PHONY: tests

# Run Mock data build
build:
	/bin/sh build.sh
.PHONY: build

# Run brew install postgres script
install_postgres:
	brew install postgresql@12
.PHONY: install_postgres

# Run brew start postgres service
start_postgres:
	brew services start postgresql@12
.PHONY: start_postgres

# Run brew stop postgres service
stop_postgres:
	brew services stop postgresql@12
.PHONY: stop_postgres

# Run brew uninstall postgres script
uninstall_postgres:
	brew uninstall postgresql@12
.PHONY: uninstall_postgres

# Create Mock Database
create_db:
	@echo "#### Creating database"
	PGPASSWORD=$(PGPASSWORD) PGUSER=$(PGUSER) PGHOST=$(PGHOST) PGPORT=$(PGPORT) psql -Atc "CREATE DATABASE $(PGDATABASE);" template1
.PHONY: create_db

# Recreate Mock Database
recreate_db:
	@echo "#### Recreating database"
	make drop_db
	make create_db
.PHONY: recreate_db

# Drop Mock Database
drop_db:
	@echo "#### Dropping database if exists"
	PGPASSWORD=$(PGPASSWORD) PGUSER=$(PGUSER) PGHOST=$(PGHOST) PGPORT=$(PGPORT) psql -Atc "DROP DATABASE IF EXISTS $(PGDATABASE);" template1
.PHONY: drop_db