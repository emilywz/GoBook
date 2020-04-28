# os

## 概述

os包封装了操作系统提供给Go的API

常用的包括：读取系统信息、读取环境变量、读取和修改文件信息等

```go
package main

import (
"os"
"fmt"
"time"
)

func main() {

    //获得当前工作路径
    //D:\BJBlockChain1801\demos\
    dir, _ := os.Getwd()
    fmt.Println(dir)

    //读取环境变量
    //D:\iWorkspace\GoPros\Go18DaysCode\Day13project\;C:\Users\sirouyang\go;D:\BJGo1801Pre\preWorks\predemos\W99\03标准库\38单元测试
    paths := os.Getenv("GOPATH")
    fmt.Println(paths)

    //修改文件的访问时间
    os.Chtimes("d:/temp/小黑子.avi",time.Now(),time.Now())

    //获得黄精变量
    environ := os.Environ()
    fmt.Println(environ)

    //获得主机名
    fmt.Println(os.Hostname())

    fmt.Println(os.IsPathSeparator('/'))//Linux认
    fmt.Println(os.IsPathSeparator('\\'))//Linux不认

    //获得文件状态信息
    fileInfo1, _ := os.Stat("d:/temp/小黑子2.avi")
    fileInfo2, _ := os.Stat("d:/temp/小黑子"+"2"+".avi")
    fmt.Println(os.SameFile(fileInfo1,fileInfo2))

    //获得用户临时文件夹所在位置
    fmt.Println(os.TempDir())

}

```

