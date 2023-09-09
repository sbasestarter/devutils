
#
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bins/http-proxy.exe ./cmd/http-proxy/http-proxy.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bins/http-proxy ./cmd/http-proxy/http-proxy.go

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bins/tcpserver.exe ./cmd/tcpserver/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bins/tcpserver ./cmd/tcpserver/main.go