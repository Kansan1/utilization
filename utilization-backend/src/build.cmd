隐藏控制台
go build -ldflags="-H=windowsgui" -o ../src.exe main.go
不隐藏
go build -o ../src.exe main.go