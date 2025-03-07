FROM golang:latest

WORKDIR /app

COPY . .
RUN rm .env .env.example .gitignore README.md
RUN rm -rf .git/ docker/
RUN go mod download

EXPOSE 8080
