FROM golang:1.24

WORKDIR /app

RUN apt-get update && apt-get install -y sqlite3

COPY go.mod go.sum ./

RUN GOPROXY=https://goproxy.cn,direct go mod download

COPY . .


RUN CGO_ENABLED=1 go build -o /app/main .

EXPOSE 8080

CMD ["./main"]
