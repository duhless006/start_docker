FROM golang:1.24.6-bookworm

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o /app/exe main.go

RUN /app/exe --help 2>/dev/null || echo "Build successful"

EXPOSE 5050

CMD ["/app/exe"] 