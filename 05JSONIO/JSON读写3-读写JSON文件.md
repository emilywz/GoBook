# 读写JSON文件

## 概述

本地文件是存储JSOP数据的一个重要方案
将Go数据以JSON字符串的形式写出到文件中的过程称之为编码
从文件中读取JSON字符串为Go数据的过程称之为解码

## 导入依赖

```go
import (
    "encoding/json"
    "os"
    "fmt"
)
```

## 定义结构体

```go
type PersonII struct {
    Name    string
    Age     int
    Sex     bool
    Hobbies []string
}
```

## 将Go数据编码到json文件 

将结构体实例以JSON格式写出到文件的过程称为编码，其一般步骤是：

准备Go数据
打开文件并创建基于该文件的JSON编码器
将Go数据编码到文件中

```go
func main31() {
    //创建PersonII实例小黑子
    xhz := PersonII{"小黑子", 60, true, []string{"撸代码", "撸项目", "撸其它的"}}

//创建文件（并打开）
filePtr, err := os.Create("d:/小黑子.json")
if err!=nil{
    fmt.Println("创建文件失败，err=",err)
    return
}
defer filePtr.Close()

//创建基于文件的JSON编码器
encoder := json.NewEncoder(filePtr)

//将小黑子实例编码到文件中
err = encoder.Encode(xhz)
if err!=nil{
    fmt.Println("编码失败，err=",err)
}else{
    fmt.Println("编码成功")
}

}
```

## 将json文件解码为Go数据 

从文件中读取JSON数据为Go语言数据的过程称之为解码，其一般步骤为：

打开含有JSON字符串的文件
创建基于改文件的JSON解码器
解码文件中的JSON数据到相应的Go数据指针中

```go
func main32() {

//预定义解码结果
var xhz PersonII

//打开文件
filePtr, _ := os.Open("d:/小黑子.avi")
defer filePtr.Close()

//创建该文件的json解码器
decoder := json.NewDecoder(filePtr)

//把解码的结果存在xhz的地址中
err := decoder.Decode(&xhz)
if err!=nil{
    fmt.Println("解码失败，err=",err)
}else {
    fmt.Printf("解码成功:%#v\n",xhz)
}

}
```


