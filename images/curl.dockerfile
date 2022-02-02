FROM alpine:3.15.0

ENV HTTP_PROXY=""
ENV HTTPS_PROXY=""
ENV http_proxy=""
ENV https_proxy=""

RUN apk add curl bash
