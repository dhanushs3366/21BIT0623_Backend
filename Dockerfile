# Dockerfile

FROM golang:1.22-alpine

RUN apk update && apk add --no-cache git

WORKDIR /app/

COPY . .

RUN go  install github.com/air-verse/air@latest
RUN go mod download

EXPOSE 8080


CMD ["air","-c",".air.toml"]
