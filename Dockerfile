FROM golang:1.23-alpine AS builder

COPY . /avito-shop/source/
WORKDIR /avito-shop/source/

RUN go build -o ./bin/avito-shop cmd/app/main.go

FROM alpine:3.13

WORKDIR /root/

COPY --from=builder /avito-shop/source/bin/ .
COPY --from=builder /avito-shop/source/docker.env .

CMD ["sh", "-c", "./avito-shop --env-path=docker.env"]
