FROM golang:1.24.2-alpine

RUN apk add --no-cache curl

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main cmd/xanny-go-template/main.go

EXPOSE 8013

CMD ["./main"]
