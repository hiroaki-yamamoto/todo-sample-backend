migrate: setup
	./db_migrate.sh up
createMigration:
	migrate create -ext sql -dir db/migrations -seq $(filter-out $@,$(MAKECMDGOALS))
updateGql:
	go tool gqlgen generate --verbose
setup:
	docker compose up -d --wait --wait-timeout 60
teardown:
	docker compose stop
clean:
	docker compose down
