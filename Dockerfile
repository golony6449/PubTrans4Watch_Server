FROM golang:1.17-bullseye
EXPOSE 8080

RUN mkdir app
WORKDIR /app

COPY . .

ENV SECRET = ${SECRET}

RUN go build .

# TODO: HEATHCHECK
# https://docs.docker.com/engine/reference/builder/#healthcheck

ENTRYPOINT ["./PubTrans4Watch_Server"]