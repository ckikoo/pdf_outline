# 初始化项目 
> go mod init 
> 
> go mod tidy
> 
# 编译项目
## linux 下编译window
> CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o output.exe main.go
## window 编译window
> go build -o output.exe main.go


# 修改项目
>
> 