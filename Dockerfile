FROM golang

ARG app_env
ENV APP_ENV $app_env

COPY ./app /go/src/goapi/app
WORKDIR /go/src/goapi/app

RUN go get ./
RUN go build

ENTRYPOINT /go/bin/app

EXPOSE 3000
