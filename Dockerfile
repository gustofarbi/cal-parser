FROM golang as builder

COPY parser /parser

WORKDIR /parser
RUN go mod vendor
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o renderer main.go

FROM alpine as runner

RUN adduser --system --shell /bin/ash cr

RUN apk add librsvg-dev

USER cr

COPY --chown=cr:cr --from=builder /parser/renderer /usr/local/bin/renderer

#ENTRYPOINT ["renderer"]

EXPOSE 80