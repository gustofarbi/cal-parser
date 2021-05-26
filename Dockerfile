FROM golang as builder

COPY parser /parser

WORKDIR /parser
RUN go mod vendor
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o renderer main.go

FROM alpine as runner

RUN addgroup -g 1000 cr && adduser -u 1000 -G cr -h /home/cr -s /bin/ash -D cr

RUN apk add librsvg-dev curl

ARG FONTS_PATH="/var/fonts"
ARG FONTS_VERSION="v1.0.2"
ARG DOWNLOAD_URL="https://codeload.github.com/myposter-de/canvas-fonts/tar.gz/${FONTS_VERSION}"

RUN mkdir $FONTS_PATH && curl $DOWNLOAD_URL | tar xz -C $FONTS_PATH --strip-components=2
ENV FONTS_PATH=$FONTS_PATH

USER cr
ARG PORT=8080
ENV PORT=$PORT
COPY --chown=cr:cr --from=builder /parser/renderer /usr/local/bin/renderer

ENTRYPOINT ["renderer"]

EXPOSE $PORT