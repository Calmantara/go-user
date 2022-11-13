# this command is for up all containers using docker compose
up:
	docker-compose -f ./deployment/docker-compose.yaml up -d

# this command is for down all docker compose containers 
down:
	docker-compose -f ./deployment/docker-compose.yaml down