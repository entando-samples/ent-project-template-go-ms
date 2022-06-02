FROM golang:alpine as builder

WORKDIR /app

COPY ./src/main/go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .

FROM  alpine

WORKDIR /app

COPY --from=builder /app/ /app/

ENTRYPOINT ["/app/entando-go-ms"]

EXPOSE 8081
