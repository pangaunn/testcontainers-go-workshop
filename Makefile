run-api:
	ENV=dev go run cmd/main.go


# r1:
	# docker-compose up

# r2:
	# docker exec -i testcontainers-go-workshop_elasticsearch_1 sh < ./seed/es/es_init.sh
	# ./seed/es/es_init.sh

# run-docker: r1 r2


# run-docker:
# 	docker-compose up -d
#  	docker exec -i search-api_elasticsearch_1 sh < ./pre-test-script/es.sh
# 	# docker build --progress=plain -t testcontainers-go-workshop -f "./Dockerfile" .
# 	# docker run -it --env-file .env -p 3000:3000 testcontainers-go-workshop
	
