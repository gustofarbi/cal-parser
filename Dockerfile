FROM golang

RUN apt update && apt install rsvg

EXPOSE 50051