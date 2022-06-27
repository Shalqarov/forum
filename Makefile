run:
	go run ./cmd/ 

docker-run:
	docker-compose up -d
	docker image prune

stop:
	docker-compose stop

delete:
	docker-compose down