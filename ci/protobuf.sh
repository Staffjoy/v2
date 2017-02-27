#!/bin/bash
set -e

# this file handles recompilation of files after a protocol buffer
# file gets modified

# HACK- generate non-gogo proto for gateway due to bugs
# (custom types, plus gateway json encoder doesn't support time.Time')
protoc \
    -I ./protobuf/ \
    -I ./vendor/ \
    -I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:../ \
    ./protobuf/account.proto
mv account/account.pb.go account/api/
protoc \
    -I ./protobuf/ \
    -I ./vendor/ \
    -I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --grpc-gateway_out=logtostderr=true:./account/ \
    ./protobuf/account.proto
mv ./account/account.pb.gw.go ./account/api/
sed -i "s/package account/package main/g" account/api/account.pb.go
sed -i "s/package account/package main/g" account/api/account.pb.gw.go

# Main account package
# account.proto -> account
# model
protoc \
    -I ./protobuf/ \
    -I ./vendor/ \
    -I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --gogo_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:../ \
    ./protobuf/account.proto
# gateway
# swagger
protoc \
    -I ./protobuf/ \
    -I ./vendor/ \
    -I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --swagger_out=logtostderr=true:./account/api/ \
    ./protobuf/account.proto
# Encode swagger
cd ./account/api/
go-bindata account.swagger.json
gofmt -s -w bindata.go
sed -i "s/Json/JSON/g" bindata.go
cd ../..

# email.proto -> email
protoc \
    -I ./protobuf/ \
    -I ./vendor/ \
    --gogo_out=plugins=grpc:.. \
    ./protobuf/email.proto
    
# sms.proto -> sms
protoc \
    -I ./protobuf/ \
    -I ./vendor/ \
    --gogo_out=plugins=grpc:.. \
    ./protobuf/sms.proto

# bot.proto -> bot
protoc \
    -I ./protobuf/ \
    -I ./vendor/ \
    -I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --gogo_out=plugins=grpc:.. \
    ./protobuf/bot.proto

# HACK- generate non-gogo proto for gateway due to bugs
# (custom types, plus gateway json encoder doesn't support time.Time')
protoc \
    -I ./protobuf/ \
    -I ./vendor/ \
    -I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:../ \
    ./protobuf/company.proto
mv company/company.pb.go company/api/
protoc \
    -I ./protobuf/ \
    -I ./vendor/ \
    -I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --grpc-gateway_out=logtostderr=true:./company/ \
    ./protobuf/company.proto
mv ./company/company.pb.gw.go ./company/api/
sed -i "s/package company/package main/g" company/api/company.pb.go
sed -i "s/package company/package main/g" company/api/company.pb.gw.go

# company.proto -> company
# model
protoc \
    -I ./protobuf/ \
    -I ./vendor/ \
    -I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --gogo_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:../ \
    ./protobuf/company.proto
# swagger
protoc \
    -I ./protobuf/ \
    -I ./vendor/ \
    -I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --swagger_out=logtostderr=true:./company/api/ \
    ./protobuf/company.proto
# Encode swagger
cd ./company/api/

go-bindata company.swagger.json
gofmt -s -w bindata.go
sed -i "s/Json/JSON/g" bindata.go
cd ../..
