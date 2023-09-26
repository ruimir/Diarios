#!/bin/bash
go get -u
go mod tidy
docker build -t 172.21.220.44:5000/aida/diarios:staging .
docker push 172.21.220.44:5000/aida/diarios:staging
