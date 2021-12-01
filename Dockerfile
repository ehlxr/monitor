FROM golang:1.17
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w " -o monitor .

FROM ehlxr/alpine
LABEL maintainer="ehlxr <ehlxr.me@gmail.com>"

WORKDIR /app

COPY --from=0 /app/monitor .

ENTRYPOINT ["./monitor"]