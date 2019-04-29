FROM golang:alpine AS build-env
ADD . /src

RUN apk add --no-cache git
RUN cd /src && go get "github.com/sawdustofmind/currency-arbitrage" \
        && go build -o arbapp

# final stage
FROM alpine
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

ENV PORT ${PORT:-12346}
WORKDIR /app
COPY --from=build-env /src/arbapp /app/
CMD ./arbapp $PORT