# JSON反序列化

## 导入依赖

```go
import (
    "encoding/json"
    "fmt"
)
```

## 定义结构体

```go
type Human struct {
    Name    string
    Age     int
    Sex     bool
    Hobbies []string
}
```

## 定义数据

```
var jsonStr,jsonStr2 string
var zqd Human
var theMap map[string]interface{}
var theSlice []map[string]interface{}
```

## 数据初始化

```go
func init() {
    jsonStr = "{\"age\":20,\"hobbies\":[\"质检\",\"维护流水线\",\"打饭\"],\"name\":\"张全蛋\",\"sex\":true}"
    jsonStr2 = "[{\"age\":20,\"hobbies\":[\"质检\",\"维护流水线\",\"打饭\"],\"name\":\"张全蛋\",\"sex\":true},{\"age\":23,\"hobbies\":[\"破坏\",\"维护流水线\",\"打人\"],\"name\":\"穆铁柱\",\"sex\":false}]"
}
```

## 反序列化为结构体

```go
func main21() {
    err := json.Unmarshal([]byte(jsonStr), &zqd)
    if err!=nil{
        fmt.Println("反序列化失败，err=",err)
    }else {
        fmt.Printf("反序列化成功：%#v\n",zqd)
    }
}
```

## 反序列化为map

```go
func main22() {
    err := json.Unmarshal([]byte(jsonStr), &theMap)
    if err!=nil{
        fmt.Println("反序列化失败，err=",err)
    }else {
        fmt.Printf("反序列化成功：%#v\n", theMap)
    }
}
```

## 反序列化为切片

```go
func main23() {
    err := json.Unmarshal([]byte(jsonStr2), &theSlice)
    if err!=nil{
        fmt.Println("反序列化失败，err=",err)
    }else {
        fmt.Printf("反序列化成功：%v\n", theSlice)
    }
}
```

