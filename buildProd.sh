#!/bin/bash
docker build -t 172.21.220.44:5000/aida/econsentimento_service:latest .
docker push 172.21.220.44:5000/aida/econsentimento_service:latest
