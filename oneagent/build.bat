SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-w -s"  -o oneagent_amd64_1.0.7
SET GOARCH=arm64
go build -ldflags "-w -s"  -o oneagent_arm64_1.0.7