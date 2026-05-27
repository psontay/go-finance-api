postgres:
	docker run --name postgres-container -p 5432:5432 -e POSTGRES_USER=elvis -e POSTGRES_PASSWORD=elvis -d postgres:15

createdb:
	docker exec -it postgres-container createdb --username=elvis --owner=elvis simple_bank

dropdb:
	docker exec -it postgres-container dropdb --username=elvis simple_bank

migrateup:
	migrate -path database/migration -database "postgresql://elvis:elvis@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path database/migration -database "postgresql://elvis:elvis@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination database/mock/store.go SimpleBank/database/sqlc Store

.PHONY: postgres createdb dropdb migratedown migrateup sqlc test server mock