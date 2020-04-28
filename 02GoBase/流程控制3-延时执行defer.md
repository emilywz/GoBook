# 延时执行defer

## defer概述

defer XXX() 延时执行，将xxx放在函数的最后执行

多个defer xxx()时，所有defer倒序执行，即最在声明的defer会最后执行

## 用途

- IO(数据库读写，文件读写，网络数据读写)
- IO资源=数据库连接，打开的文件对象，网路连接
- IO资源开销（CPU，内存，磁盘…）巨大
- IO资源，随用随开，用完【即】关

## 场景

- IO只是场景之一
- 其它场景：消费完毕要买单，运动完要洗澡，自习结束要关灯…

## 读写数据库，并在程序的最后关闭数据库

```go
func demo41() {
    //随用随开，用完【即】关
    fmt.Println("打开数据库")
    //defer所唤起的函数将在函数结束前才执行
    defer closeDatabase()

    fmt.Println("愉快地读写数据")
    fmt.Println("读写完毕")
}
```

## 多个defer时，最早defer的操作最后执行

```go
func demo42() {
    //打开数据库
    fmt.Println("打开数据库")
    defer closeDatabase()
    //读入DB数据
    fmt.Println("读入DB数据")
    //打开文件
    fmt.Println("打开文件")
    defer closeFile()
    //向文件中写出DB中的数据
    fmt.Println("读入DB数据")
    //关闭文件
    //继续操作数据库
    fmt.Println("继续操作数据库")
    //关闭数据库
}

func closeDatabase() {
    fmt.Println("关闭数据库")
}

func closeFile() {
    fmt.Println("关闭文件")
}
```

