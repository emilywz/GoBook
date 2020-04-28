# iota定义常量组

## 概述

第一个常量的iota，值为0

后面的常量自动沿用第一个常量的表达式

后面常量中的iota会自动递增

```go
package main

import "fmt"

//常量组中的iota会从0开始自动增长
const (
    Sunday    = iota //0
    Monday
    Tuesday
    Wednesday
    Thursday
    Friday
    Saturday
)

//常量组中的iota会从0开始自动增长
/*const (
    D = 1 << iota //1
    C = 1 << iota //10
    B = 1 << iota //100
    A = 1 << iota //100
    S = 1 << iota //1000
)*/
/*const (
    D = 1 << iota //1

    //沿用第一个变量的计算方式，iota每次+1
    C
    B
    A
    S
)*/
const (
    D = 2 * iota //0

    //沿用第一个变量的计算方式，iota每次+1
    C  //2
    B  //4
    A  //6
    S  //8
)

//定义一组变量
var (
    a int = 123
    b int = 456
)

func main() {
    fmt.Println(Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday)
    fmt.Println(S, A, B, C, D)
}

```

