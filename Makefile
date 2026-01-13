service-run:
	@go run main.go

some-target:
	docker run -d -p 8080:8080 --name go-app corp-user && \
	docker run corp-user
