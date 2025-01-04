cd src

go env -w GO111MODULE=on

go mod tidy

env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X bfv-bot/common/config.GitCommit=$(git rev-parse HEAD) -X bfv-bot/common/config.GitBranch=$(git rev-parse --abbrev-ref HEAD) -X bfv-bot/common/config.BuildTime=$(date +%Y-%m-%dT%H:%M:%S)" -tags release -o release/linux-amd64-bfv-bot

env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X bfv-bot/common/config.GitCommit=$(git rev-parse HEAD) -X bfv-bot/common/config.GitBranch=$(git rev-parse --abbrev-ref HEAD) -X bfv-bot/common/config.BuildTime=$(date +%Y-%m-%dT%H:%M:%S)" -tags release -o release/windows-x64-bfv-bot.exe