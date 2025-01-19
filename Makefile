postgres:
	docker run --name mypostgres -e POSTGRES_PASSWORD=letme1n -e POSTGRES_USER=root -p 5433:5432 -d postgres:16-alpine

redis:
	docker run --name myredis -p 6379:6379 -d redis:7.4-alpine

createdb:
	docker exec -it mypostgres createdb --username=root --owner=root vecsys

createextension:
	docker exec -it mypostgres psql -U root -d vecsys -c 'CREATE EXTENSION IF NOT EXISTS "uuid-ossp";'

dropdb:
	docker exec -it mypostgres dropdb vecsys

.PHONY: postgres createdb dropdb redis createextension deploy_dev