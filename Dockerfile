FROM golang

RUN apt update && apt install -y librsvg2-bin

EXPOSE 50051