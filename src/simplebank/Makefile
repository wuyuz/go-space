postgres: 
	docker run -d --name pgsql -p 5432:5432  -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -e PGDATA=/var/lib/postgresql/data/pgdata -e TZ=PRC -v /Users/alex/Desktop/Program\ Project/pg/docker/volume/pg/:/var/lib/postgresql/data postgres:12.1

createdb:
	docker exec -it pgsql createdb --username=root --owner=root simplebank

dropdb:
	docker exec -it pgsql dropdb  simplebank

migrateup:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simplebank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simplebank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test