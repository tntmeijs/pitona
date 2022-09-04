export GOOS="linux"
export GOARCH=arm
export GOARM=7

go build -buildmode=exe -o pitona main.go
