if not exist "scripts" (cd ..)
if not exist "scripts" (cd ..)
if not exist "scripts" (cd ..)
docker run --rm ^
-v %cd%:/go/src/github.com/aura-studio/excel-converter ^
-w /go/src/github.com/aura-studio/excel-converter ^
-e GOBIN=/go/bin ^
-e GOPATH=/go ^
-e GOPROXY=https://goproxy.cn,direct ^
-e CGO_ENABLED=0 ^
-e GOARCH=amd64 ^
-w /go/src/github.com/aura-studio/excel-converter golang:latest  go build -o converter-unix .
