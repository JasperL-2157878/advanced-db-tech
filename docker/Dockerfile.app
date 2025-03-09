FROM golang:latest

WORKDIR /app

COPY ./public ./public
COPY ./source ./source
COPY go.mod ./go.mod 
COPY go.sum ./go.sum
RUN go mod download

EXPOSE 8080
