migrate: setup
	./db_migrate.sh up
createMigration:
	migrate create \
		-ext sql \
		-dir db/migrations \
		-seq $(filter-out $@,$(MAKECMDGOALS))
test:
	go test -v ./...
# gorm:
# 	gorm gen -i ./db/models -o ./db/repos
updateGql:
	go tool gqlgen generate --verbose
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
