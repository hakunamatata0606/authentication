db_url = "mysql://bao:123@tcp(172.17.0.2:3306)/test"

migrate_up:
	migrate -database $(db_url) -path db/migration up

migrate_down:
	migrate -database $(db_url) -path db/migration down 

migrate_force:
	migrate -database $(db_url) -path db/migration force 1

sqlc:
	sqlc generate

test:
	ROOT_DIR=/home/bao/go_proj/authentication go test ./...

.PHONY: migrate_up, migrate_down,  migrate_force, sqlc, test