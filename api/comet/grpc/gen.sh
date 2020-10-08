#!/bin/zsh

protoc --proto_path=/Users/mike/.gvm/pkgsets/go1.14.9/global/src \
--proto_path=/Users/mike/.gvm/pkgsets/go1.14.9/global/src/github.com/go-kratos/kratos/third_party \
--proto_path=/Users/mike/golang/github/goim/api/common/grpc \
--proto_path=/Users/mike/golang/github/goim/api/comet/grpc \
--gofast_out=plugins=grpc:. api.proto