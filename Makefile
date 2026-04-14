migrate: setup
	./db_migrate.sh todo up
migrateTestDatabase: setup
	docker compose exec db psql -U postgres -c "DROP DATABASE IF EXISTS todo_test;"
	docker compose exec db psql -U postgres -c "CREATE DATABASE todo_test;"
	./db_migrate.sh todo_test up
teardownTestDatabase:
	docker compose exec db psql -U postgres -c "DROP DATABASE IF EXISTS todo_test;"
createMigration:
	migrate create \
		-ext sql \
		-dir db/migrations \
		-seq $(filter-out $@,$(MAKECMDGOALS))
test: migrateTestDatabase
	go test -v ./...
	make teardownTestDatabase
# gorm:
# 	gorm gen -i ./db/models -o ./db/repos
updateGql:
	make -C graph
	make -C auth
setup:
	docker compose up -d --wait --wait-timeout 60
teardown:
	docker compose stop
clean:
	docker compose down
mock:
	mockgen -source=./db/repos/todo/iface.go \
		-destination=./db/repos/todo/mock.go \
		-package=todo

%:
	@:
