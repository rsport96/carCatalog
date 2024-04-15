FROM golang:latest

COPY . /app

WORKDIR /app

RUN go build -o main .

EXPOSE 8080 80

CMD ["./main"]