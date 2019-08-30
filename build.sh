#!/bin/bash -
declare -r Name="ecs-task-monitoring"

for GOOS in darwin linux; do
    GO111MODULE=on GOOS=$GOOS GOARCH=amd64 go build -o bin/cron-monitoring-$GOOS-amd64 *.go
done
