# pdf_outline
# 初始化项目 
> go mod init file
> 
> go mod tidy
> 
# 编译项目
## linux 下编译window
> CC=x86_64-w64-mingw64-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o output.exe main.go
## window 编译window
> go build -o output.exe main.go


# 修改项目
>
> 