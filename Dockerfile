FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apt-get update && apt-get install -y sqlite3

RUN CGO_ENABLED=1 go build -o /app/main .

EXPOSE 8080

CMD ["./main"]
