run-docker:
	docker-compose up --build
	
int-test:
	ginkgo -r --label-filter="integration"