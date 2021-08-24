FROM golang:latest

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o books ./cmd/books

# you may want to copy config file from elsewhere

ENTRYPOINT ["./books"]