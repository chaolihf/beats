#!/bin/sh
GOARCH=amd64 go build -ldflags '-w -s'  -o oneagent_amd64_1.0.2
GOARCH=arm64 go build -ldflags '-w -s'  -o oneagent_arm64_1.0.2