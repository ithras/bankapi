# syntax=docker/dockerfile:1

FROM golang:alpine AS build

WORKDIR /app

COPY /src .

RUN go mod download
RUN go build -o /bankapi

# Deploy
FROM alpine

WORKDIR /

COPY --from=build /bankapi /bankapi

EXPOSE 8080

ENTRYPOINT [ "/bankapi" ]