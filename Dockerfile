FROM golang:1.17-alpine

RUN apk update && apk upgrade && apk add --no-cache bash git && apk add --no-cache chromium

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY . .