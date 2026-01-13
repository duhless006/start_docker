FROM golang:1.24.6-bookworm

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o /app/exe main.go

EXPOSE 8080

CMD ["/app/exe"]