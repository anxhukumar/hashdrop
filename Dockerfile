FROM ubuntu:24.04

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY ./server/build/hashdrop-server /app/hashdrop-server
COPY ./server/internal/sql/migrations /app/internal/sql/migrations

CMD ["/app/hashdrop-server"]