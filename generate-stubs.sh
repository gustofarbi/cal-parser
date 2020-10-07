#!/bin/bash

protoc --go_out=svg/generated --go_opt=paths=source_relative \
      --go-grpc_out=svg/generated --go-grpc_opt=paths=source_relative car.proto
