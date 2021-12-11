# syntax=docker/dockerfile:1

FROM golang:1.16-alpine AS build

WORKDIR /app

COPY /src/go.mod ./
COPY /src/go.sum ./
RUN go mod download

COPY /src .

RUN go get -d -v
RUN go build -v
RUN go build -o /bankapi

#
# Deploy
#

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /bankapi /bankapi

EXPOSE 8080

ENTRYPOINT [ "/bankapi" ]