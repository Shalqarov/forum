.PHONY: build

build:
	docker-compose build
	docker image prune

run: build
	docker-compose up -d

stop:
	docker-compose stop

delete:
	docker-compose down