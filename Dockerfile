FROM golang:latest

WORKDIR /app

COPY . .

RUN go get -d -v ./...

CMD ["go","run","main.go"] 