#!/usr/bin/env bash
pkill main
source .env
go run api/main.go &

fswatch api adapters domain core providers | while read line
do
    killall -9 main;
    go run api/main.go &
done
