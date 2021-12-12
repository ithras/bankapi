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

# Deploy
FROM postgres:10.0-alpine

USER postgres

RUN chmod 0700 /var/lib/postgresql/data &&\
    initdb /var/lib/postgresql/data &&\
    echo "host all  all    0.0.0.0/0  md5" >> /var/lib/postgresql/data/pg_hba.conf &&\
    echo "listen_addresses='*'" >> /var/lib/postgresql/data/postgresql.conf &&\
    pg_ctl start &&\
    psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'bank'" | grep -q 1 || psql -U postgres -c "CREATE DATABASE bank" &&\
    psql -c "ALTER USER postgres WITH ENCRYPTED PASSWORD 'banktest';"

WORKDIR /

COPY --from=build /bankapi /bankapi

EXPOSE 8080

ENTRYPOINT [ "/bankapi" ]