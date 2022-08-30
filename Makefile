run-api:
	ENV=dev go run cmd/api/main.go

run-docker:
	docker build --progress=plain -t testcontainers-go-workshop -f "./Dockerfile" .
	docker run -it --env-file .env -p 3000:3000 testcontainers-go-workshop

run-docker-compose:
	docker-compose up
	
