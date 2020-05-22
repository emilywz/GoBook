# 执行HTTP的GET&POST请求

## 导入依赖包

```go
import (
    "fmt"
    "net/http"
    "io/ioutil"
    "strings"
)
```

## 提交GET请求并获得返回

```go
func main521() {
    url := "http://www.baidu.com/s?wd=肉"

    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("错误")
    }
    defer resp.Body.Close()

    bodyBytes, _ := ioutil.ReadAll(resp.Body) //读取信息
    fmt.Println(string(bodyBytes))            //读取网页源代码
}
```

## 提交POST请求并获得返回

```go
func main522() {
    //url := "http://www.baidu.com"
    url := "https://httpbin.org/post?name=张三"

    resp, err := http.Post(
        url,
        "application/x-www-form-urlencoded",
        strings.NewReader("id=nimei"))
    if err != nil {
        fmt.Println("错误")
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body) //读取信息
    fmt.Println(string(body))            //读取网页源代码

}
```

