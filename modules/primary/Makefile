CONTAINER_NAME=analytics_postgres
DB_NAME=dev_primary

db:
	docker exec -it ${CONTAINER_NAME} psql -U postgres -c "DROP DATABASE IF EXISTS ${DB_NAME};"
	docker exec -it ${CONTAINER_NAME} psql -U postgres -c "CREATE DATABASE ${DB_NAME};"
	docker run -v ${CURDIR}/migrations:/migrations --network host migrate/migrate \
    	-path=/migrations/ -database postgres://postgres@localhost:5440/${DB_NAME}?sslmode=disable up
.PHONY: db