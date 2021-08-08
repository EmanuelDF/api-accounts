FROM golang:alpine AS build-env

ARG APPNAME
ENV SRCPATH $GOPATH/src/github.com/emanueldf/$APPNAME

COPY ./ $SRCPATH

RUN go install github.com/emanueldf/$APPNAME/cmd/$APPNAME

FROM alpine

ARG APPNAME
ARG TAGS
ENV SERVICE_TAGS=$TAGS,active
ENV LOG_FORMAT=json
ENV AWS_SDK_LOAD_CONFIG=1
RUN apk add --no-cache dumb-init ca-certificates

WORKDIR /app
COPY build/package/$APPNAME/entrypoint.sh /app/
COPY --from=build-env /go/bin/$APPNAME /app/

COPY internal/app/$APPNAME/api /app/api

RUN addgroup -S appuser && adduser -S -G appuser appuser
RUN chown -R appuser:appuser /app

USER appuser

ENTRYPOINT ["./entrypoint.sh"]