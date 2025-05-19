postgres:
	podman run -d --name onetap-postgres -e POSTGRES_USER=onetapuser -e POSTGRES_PASSWORD=onetappassword -e POSTGRES_DB=onetapdb -p 5433:5432 docker.io/library/postgres:latest
mig:
	migrate -database "postgres://onetapuser:onetappassword@localhost:5433/onetapdb?sslmode=disable" -path migrations up

migdown:
	migrate -database "postgres://onetapuser:onetappassword@localhost:5433/onetapdb?sslmode=disable" -path migrations down

migreset:
	migrate -database "postgres://onetapuser:onetappassword@localhost:5433/onetapdb?sslmode=disable" -path migrations reset