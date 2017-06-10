#!/usr/bin/env bash
pkill main
source .env
go run api/main.go &

inotifywait -m -r -e close_write adapters domain core providers | while read line
do
  pkill main && go run api/main.go &
done
