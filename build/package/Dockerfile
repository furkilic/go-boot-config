FROM alpine:3.11.5

ARG EXEC_DIR="/bin/go-boot-config"
COPY $EXEC_DIR /data/app/go-boot-config
COPY /configs /data/app/configs

ENV GO_ENV=from-env

WORKDIR /data/app

CMD ["/data/app/go-boot-config", "--go.cmd=from-cmd"]