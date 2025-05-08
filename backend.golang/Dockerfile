FROM golang:1.23

WORKDIR /app

COPY . .

RUN go mod tidy

CMD ["make", "run"]
