FROM golang as builder

COPY parser /parser

RUN cd /parser && \
    go mod vendor && \
    go build -o renderer main.go
#RUN apt update && apt install -y librsvg2-bin

FROM alpine as runner

COPY --from=builder /parser/renderer /usr/local/bin/renderer

RUN apk add librsvg-dev

ENTRYPOINT ["render"]

EXPOSE 80