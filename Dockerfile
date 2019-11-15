FROM ehlxr/alpine

LABEL maintainer="ehlxr <ehlxr.me@gmail.com>"

COPY ./dist/monitor_linux_amd64 /usr/local/bin/monitor


ENTRYPOINT ["/usr/local/bin/monitor"]