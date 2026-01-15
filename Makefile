include .env
export

service-run:
	@go run main.go

service-deploy:
	docker compose up -d app

service-undeploy:
	docker compose dowm -d app

run-hhtp-app:
	docker run -d -p 8080:8080 --name go-app corp-user2 && \
	docker run corp-user


	