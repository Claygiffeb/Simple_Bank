# Buiding stages
FROM golang:1.21.6-alpine3.19 AS building
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz  | tar xvz
       

#Running stage
FROM alpine:3.19
WORKDIR /app
COPY --from=building /app/main .
COPY --from=building /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
