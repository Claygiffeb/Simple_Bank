# Buiding stages
FROM golang:1.21.6-alpine3.19 AS building
WORKDIR /app
COPY . .
RUN go build -o main main.go


#Running stage
FROM alpine:3.19
WORKDIR /app
COPY --from=building /app/main .
COPY app.env .

EXPOSE 8080
CMD ["/app/main"]
