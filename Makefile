migrate_up:
	migrate -database "mysql://bao:123@tcp(172.17.0.2:3306)/test" -path db/migration up

migrate_down:
	migrate -database "mysql://bao:123@tcp(172.17.0.2:3306)/test" -path db/migration down 

migrate_force:
	migrate -database "mysql://bao:123@tcp(172.17.0.2:3306)/test" -path db/migration force 1

sqlc:
	sqlc generate

.PHONY: migrate_up, migrate_down,  migrate_force, sqlc