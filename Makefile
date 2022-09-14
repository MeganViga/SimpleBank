postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres
createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank 

dropdb:
	docker exec -it postgres dropdb simpdocker psle_bank

accessdb:
	docker exec -it postgres psql -U root simple_bank

migratecreate:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
		migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go


.PHONY: postgres createdb acessdb dropdb migratecreate migrateup migratedown sqlc test server