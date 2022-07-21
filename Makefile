run:
	docker-compose up 

docker-run:
	docker-compose up 

stop:
	docker-compose stop

delete:
	docker-compose down
	docker volume rm forum_default
	docker volume rm forum_pg-data 

remove-images:
	docker rmi -f $(docker images -aq)