FROM golang:1.22-alpine as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o ./bin/server .

FROM alpine:3.19.1
COPY --from=build-base /app/bin/server /app/bin/server

EXPOSE 8080

CMD ["/app/bin/server"]