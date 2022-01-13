FROM alpine:3.15.0

RUN apk update && apk add curl

ARG FILE=/etc/profile
