FROM golang:1.27-alpine

WORKDIR /app

COPY main.go .

RUN CGO_ENABLED=0 go build -o chall main.go

EXPOSE 1337

CMD ["./chall"]