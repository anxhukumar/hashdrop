#!/bin/bash

cd server
GOOS=linux GOARCH=amd64 go build -o ./build/hashdrop-server