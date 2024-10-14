postgres:
	docker run --name mypostgres -e POSTGRES_PASSWORD=letme1n -e POSTGRES_USER=root -p 5433:5432 -d postgres:16-alpine

redis:
	docker run --name myredis -p 6379:6379 -d redis:7.4-alpine

createdb:
	docker exec -it mypostgres createdb --username=root --owner=root vecsys

dropdb:
	docker exec -it mypostgres dropdb vecsys

deploy_dev:
	env CGO_ENABLED=0 go build -o vecsys
	scp -P 2285 vecsys root@182.253.225.226:~/vecsys

.PHONY: postgres createdb dropdb redis deploy_dev