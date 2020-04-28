# Go语言开发环境

## 开发环境搭建

### 下载地址

- Go编译器 
  https://golang.google.cn/dl/
- Goland官网下载 
  https://www.jetbrains.com/go/download/#section=windows
- 在线激活地址 
  [http://idea.youbbs.org](http://idea.youbbs.org/)

### 配置环境变量

Go 语言依赖一个重要的环境变量：$GOPATH 

GOPATH允许多个目录，当有多个目录时，请注意分隔符，多个目录的时候Windows是分号，Linux系统是冒号，当有多个GOPATH时，默认会将go get的内容放在第一个目录下。 
$GOPATH 目录约定有三个子目录：

src 存放源代码（比如：.go .c .h .s等）
pkg 编译后生成的文件（比如：.a）
bin 编译后生成的可执行文件（为了方便，可以把此目录加入到 $PATH 变量中）

查看变量是否设置成功  ：

```shell
go env
```

