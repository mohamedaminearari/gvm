# Windows
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o build/gvm-windows-amd64.exe .
$env:GOOS="windows"; $env:GOARCH="arm64"; go build -o build/gvm-windows-arm64.exe .

# Linux
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o build/gvm-linux-amd64 .
$env:GOOS="linux"; $env:GOARCH="arm64"; go build -o build/gvm-linux-arm64 .

# macOS
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o build/gvm-darwin-amd64 .
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o build/gvm-darwin-arm64 .