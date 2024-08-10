createdb:
	docker exec -it postgres createdb --username=root --owner=root simple-bank-app

dropdb:
	docker exec -it postgres dropdb simple-bank-app	

postgres:
	docker run --name postgres --network simple-bank-app -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

migrateup:
	migrate -path db/migration -database "postgresql://postgres:cong1234@localhost:5432/simple-bank-app?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:cong1234@localhost:5432/simple-bank-app?sslmode=disable" -verbose down

sqlc:
	sqlc generate	

test:
	go test -v -cover -short ./...