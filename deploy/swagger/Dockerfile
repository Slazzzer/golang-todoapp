FROM golang:1.26.3-bookworm

RUN apt-get update \
    && apt-get install -y --no-install-recommends git \
    && rm -rf /var/lib/apt/lists/*

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6

WORKDIR /code

ENTRYPOINT ["swag"]
