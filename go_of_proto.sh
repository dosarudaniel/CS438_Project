#!/usr/bin/env bash

protoc services/chord_service/*.proto --go_out=plugins=grpc:$GOPATH/src
protoc services/client_service/*.proto --go_out=plugins=grpc:$GOPATH/src
protoc services/file_share_service/*.proto --go_out=plugins=grpc:$GOPATH/src

