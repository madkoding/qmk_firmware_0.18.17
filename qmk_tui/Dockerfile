FROM alpine:3.19

RUN apk add --no-cache \
    go \
    gcc \
    musl-dev \
    libx11 \
    libxext \
    libxrender \
    fontconfig \
    ttf-dejavu \
    ncurses-libs

WORKDIR /app

COPY qmk_tui .
COPY util/qmk_wrapper.sh /qmk_wrapper.sh

ENV TERM=xterm-256color
ENV QMK_PATH=/qmk

ENTRYPOINT ["./qmk_tui"]
