# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app/inventory-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
COPY config.docker.json config.json 
COPY config.docker.yaml config.yaml 

RUN go build -o ./out/inventory-service .

EXPOSE 8080

CMD [ "./out/inventory-service" ]