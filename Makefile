.PHONY: build deploy redeploy start stop down db db-stop db-down


build:
	go build .

start: deploy

# automatically builds and runs this in a docker container
deploy:
	docker compose up -d --build

redeploy:
	docker compose up -d --force-recreate --build

stop:
	docker compose down

down: stop

db:
	docker compose -f ./docker-compose.db.yaml up -d

db-stop:
	docker compose -f ./docker-compose.db.yaml down

db-down: db-stop
