# syntax=docker/dockerfile:1

## Build
FROM golang:1.19.3-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /quotes

## Deploy
FROM alpine:3.17.0
ARG USER=goapp
ENV HOME /home/$USER

RUN apk add --update bash && rm -rf /var/cache/apk/*
COPY --from=build /quotes $HOME/quotes

# add new user
RUN adduser -D $USER \
        && chown $USER:$USER $HOME/quotes \
        && mkdir -p /etc/sudoers.d \
        && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
        && chmod 0440 /etc/sudoers.d/$USER

USER $USER
WORKDIR $HOME

EXPOSE 8080

ENTRYPOINT $HOME/quotes