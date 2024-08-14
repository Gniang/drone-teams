GOOS=linux
GOARCH=amd64
CGO_ENABLED=0
GO111MODULE=on
go build -a -tags netgo -o release/linux/amd64/drone-teams ./cmd/drone-teams
docker build --rm -t gniang/drone-teams -f ./docker/Dockerfile.linux.amd64 .