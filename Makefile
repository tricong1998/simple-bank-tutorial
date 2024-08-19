createdb:
	docker exec -it postgres createdb --username=root --owner=root simple-bank-app

dropdb:
	docker exec -it postgres dropdb simple-bank-app	

postgres:
	docker run --name postgres --network simple-bank-app -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

migrateup:
	migrate -path db/migration -database "postgresql://postgres:cong1234@localhost:5432/simple-bank-app?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://postgres:cong1234@localhost:5432/simple-bank-app?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://postgres:cong1234@localhost:5432/simple-bank-app?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://postgres:cong1234@localhost:5432/simple-bank-app?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate	

test:
	go test -v -cover -short ./...

server:
	go run main.go

mockgen:
	mockgen -package mockdb -destination db/mock/store.go github.com/Sotatek-CongNguyen/simple-bank-practice/db/sqlc Store

.PHONY: postgres createdb dropdb postgres migrateup migratedown sqlc test server mockgen