@REM 构造linux amd64二进制文件
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o myapp-linux main.go

@REM 构造windows amd64二进制文件
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
go build -o myapp-windows.exe main.go