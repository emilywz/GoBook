# JSON序列化

## 概述

 序列化：将Go数据转化为json字符串 

## 定义结构体

```go
type Person struct {
    Name    string
    Age     int
    Sex     bool
    Hobbies []string
}
```

## 定义数据

```go
var (
    zhangqd, mutz Person
    ps            []Person
    mMap, yMap    map[string]interface{}
    mSlice        []map[string]interface{}
)
```

## 初始化数据

```go
func init() {
    //初始化结构体数据
    zhangqd = Person{"张全蛋", 20, true, []string{"质检", "维护流水线", "打饭"}}
    mutz = Person{"穆铁花", 23, false, []string{"搞破坏", "维护流水线", "打人"}}
    ps = make([]Person, 0)
    ps = append(ps, zhangqd, mutz)

    //初始化map数据
    mMap = make(map[string]interface{})
    mMap["name"] = "张全蛋"
    mMap["age"] = 20
    mMap["sex"] = true
    mMap["hobbies"] = []string{"质检", "维护流水线", "打饭"}
    yMap = make(map[string]interface{})
    yMap["name"] = "穆铁柱"
    yMap["age"] = 23
    yMap["sex"] = false
    yMap["hobbies"] = []string{"破坏", "维护流水线", "打人"}

    //初始化切片数据
    mSlice = make([]map[string]interface{}, 0)
    mSlice = append(mSlice, mMap, yMap)

}
```

## 序列化结构体

```go
func main11() {
    //bytes, err := json.Marshal(zhangqd)
    bytes, err := json.Marshal(ps)
    if err != nil {
        fmt.Println("序列化失败,err=", err)
    } else {
        jsonStr := string(bytes)
        fmt.Println("序列化结果：", jsonStr)
    }
}
```

## 序列化map

```go
func main12() {
    bytes, err := json.Marshal(mMap)
    if err != nil {
        fmt.Println("序列化失败,err=", err)
    } else {
        jsonStr := string(bytes)
        fmt.Println("序列化结果：", jsonStr)
    }
}
```

## 序列化切片

```go
func main13() {
    bytes, err := json.Marshal(mSlice)
    if err != nil {
        fmt.Println("序列化失败,err=", err)
    } else {
        jsonStr := string(bytes)
        fmt.Println("序列化结果：", jsonStr)
    }
}
```