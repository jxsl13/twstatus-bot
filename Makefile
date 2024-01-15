.PHONY: build deploy redeploy start stop db db-stop


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

db:
	docker compose -f ./docker-compose.db.yaml up -d

db-stop:
	docker compose -f ./docker-compose.db.yaml down
