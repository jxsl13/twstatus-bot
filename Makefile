.PHONY: build deploy redeploy start stop


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
