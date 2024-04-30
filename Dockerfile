FROM golang:latest

ENV GOOS=linux \
    CGO_ENABLED=0 \
    GOARCH=amd64

WORKDIR /app

COPY ./ ./

RUN go mod download

RUN go build -o app main.go

EXPOSE 8000

CMD ["./app"]   