help:
	@echo "docker-build - to build docker image, name of the builded image forum-docker"
	@echo "docker-run - to run docker container, name of the container forum-container, on port 8080:8080"
	@echo "docker-stop - to stop docker container"
	@echo "docker-start - to start docker container"
	@echo "go-run - to run server without docker"

docker-build:
	docker image build -f Dockerfile . -t forum-docker
	docker rmi $$(docker images -f "dangling=true" -q)
	docker images

docker-run:
	docker container run -p 8080:8080 --rm -d --name forum-container forum-docker
	docker ps -a

docker-stop:
	docker stop forum-container
	docker ps -a

docker-start:
	docker start forum-container
	docker ps -a

go-run:
	go run cmd/main.go