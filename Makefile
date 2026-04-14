.PHONY: build

build:
	go build

dockerUp:
	docker compose up -d --wait --wait-timeout 60

migrate: setup
	./db_migrate.sh todo up

setupTest: setup
	docker compose exec db psql -U postgres -c "DROP DATABASE IF EXISTS todo_test;"
	docker compose exec db psql -U postgres -c "CREATE DATABASE todo_test;"
	./db_migrate.sh todo_test up

cleanTest:
	docker compose exec db psql -U postgres -c "DROP DATABASE IF EXISTS todo_test;"

createMigration:
	migrate create \
		-ext sql \
		-dir db/migrations \
		-seq $(filter-out $@,$(MAKECMDGOALS))

testGo:
	go test -p 1 -v ./...

test: setupTest | testGo | cleanTest

updateGql:
	make -C graph
	make -C auth

setup: dockerUp | migrate

stop:
	docker compose stop

clean:
	docker compose down

mock:
	mockgen -source=./db/repos/todo/iface.go \
		-destination=./db/repos/todo/mock.go \
		-package=todo

%:
	@:
