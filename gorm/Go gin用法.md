# Go gin用法

Gin 是一个 go 写的 web 框架，具有高性能的优点。官方地址：https://github.com/gin-gonic/gin

带目录请移步 http://xf.shuangdeyu.com/movie/content.html?mid=25，简书markdown不支持目录生成

**目录**

[TOC]

# 安装

要安装Gin包，首先需要安装Go并设置Go工作区

1、下载并安装

> ```shell
> $ go get -u github.com/gin-gonic/gin
> ```

2、在代码中导入它

> ```go
> import "github.com/gin-gonic/gin"
> ```

**使用包管理工具Govendor安装**

1、`go get` govendor(安装)

> ```go
> $ go get github.com/kardianos/govendor
> ```

2、创建项目文件夹并进入文件夹

> ```shell
> mkdir -p GOPATH/src/github.com/myusername/project && cd "$_"
> ```

3、初始化项目并添加 gin

> $ govendor init
> 
> $ govendor fetch [github.com/gin-gonic/gin@v1.3](http://github.com/gin-gonic/gin@v1.3)

4、复制一个模板到你的项目

> ```shell
> $ curl https://raw.githubusercontent.com/gin-gonic/gin/master/examples/basic/main.go > main.go
> ```

5、运行项目

> ```shell
> $ go run main.go
> ```

# 前提

使用gin需要Go的版本号为1.6或更高

# 快速入门

运行这段代码并在浏览器中访问 http://localhost:8080

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.Run() // listen and serve on 0.0.0.0:8080
}
```

# 代码示例

## 使用 GET, POST, PUT, PATCH, DELETE, OPTIONS

```go
func main() {
    // Disable Console Color
    // gin.DisableConsoleColor()

    // 使用默认中间件创建一个gin路由器
    // logger and recovery (crash-free) 中间件
    router := gin.Default()

    router.GET("/someGet", getting)
    router.POST("/somePost", posting)
    router.PUT("/somePut", putting)
    router.DELETE("/someDelete", deleting)
    router.PATCH("/somePatch", patching)
    router.HEAD("/someHead", head)
    router.OPTIONS("/someOptions", options)

    // 默认启动的是 8080端口，也可以自己定义启动端口
    router.Run()
    // router.Run(":3000") for a hard coded port
}
```

## 获取路径中的参数

```go
func main() {
    router := gin.Default()

    // 此规则能够匹配/user/john这种格式，但不能匹配/user/ 或 /user这种格式
    router.GET("/user/:name", func(c *gin.Context) {
        name := c.Param("name")
        c.String(http.StatusOK, "Hello %s", name)
    })

    // 但是，这个规则既能匹配/user/john/格式也能匹配/user/john/send这种格式
    // 如果没有其他路由器匹配/user/john，它将重定向到/user/john/
    router.GET("/user/:name/*action", func(c *gin.Context) {
        name := c.Param("name")
        action := c.Param("action")
        message := name + " is " + action
        c.String(http.StatusOK, message)
    })

    router.Run(":8080")
}
```

## 获取Get参数

```go
func main() {
    router := gin.Default()

    // 匹配的url格式:  /welcome?firstname=Jane&lastname=Doe
    router.GET("/welcome", func(c *gin.Context) {
        firstname := c.DefaultQuery("firstname", "Guest")
        lastname := c.Query("lastname") // 是 c.Request.URL.Query().Get("lastname") 的简写

        c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
    })
    router.Run(":8080")
}
```

## 获取Post参数

```go
func main() {
    router := gin.Default()

    router.POST("/form_post", func(c *gin.Context) {
        message := c.PostForm("message")
        nick := c.DefaultPostForm("nick", "anonymous") // 此方法可以设置默认值

        c.JSON(200, gin.H{
            "status":  "posted",
            "message": message,
            "nick":    nick,
        })
    })
    router.Run(":8080")
}
```

## Get + Post 混合

```go
示例：
POST /post?id=1234&page=1 HTTP/1.1
Content-Type: application/x-www-form-urlencoded

name=manu&message=this_is_great
```

```go
func main() {
    router := gin.Default()

    router.POST("/post", func(c *gin.Context) {

        id := c.Query("id")
        page := c.DefaultQuery("page", "0")
        name := c.PostForm("name")
        message := c.PostForm("message")

        fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
    })
    router.Run(":8080")
}
```

```
结果：id: 1234; page: 1; name: manu; message: this_is_great
```

## 上传文件

### 单文件上传

参考问题 [#774](https://github.com/gin-gonic/gin/issues/774)，细节 [example code](https://github.com/gin-gonic/gin/tree/master/examples/upload-file/single)

慎用 `file.Filename` ，参考 [Content-Disposition on MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition#Directives) 和 [#1693](https://github.com/gin-gonic/gin/issues/1693)

> 上传文件的文件名可以由用户自定义，所以可能包含非法字符串，为了安全起见，应该由服务端统一文件名规则

```go
func main() {
    router := gin.Default()
    // 给表单限制上传大小 (默认 32 MiB)
    // router.MaxMultipartMemory = 8 << 20  // 8 MiB
    router.POST("/upload", func(c *gin.Context) {
        // 单文件
        file, _ := c.FormFile("file")
        log.Println(file.Filename)

        // 上传文件到指定的路径
        // c.SaveUploadedFile(file, dst)

        c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
    })
    router.Run(":8080")
}
```

`curl` 测试：

```shell
curl -X POST http://localhost:8080/upload \
  -F "file=@/Users/appleboy/test.zip" \
  -H "Content-Type: multipart/form-data"
```

### 多文件上传

详细示例：[example code](https://github.com/gin-gonic/gin/tree/master/examples/upload-file/multiple)

```go
func main() {
    router := gin.Default()
    // 给表单限制上传大小 (默认 32 MiB)
    // router.MaxMultipartMemory = 8 << 20  // 8 MiB
    router.POST("/upload", func(c *gin.Context) {
        // 多文件
        form, _ := c.MultipartForm()
        files := form.File["upload[]"]

        for _, file := range files {
            log.Println(file.Filename)

            // 上传文件到指定的路径
            // c.SaveUploadedFile(file, dst)
        }
        c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
    })
    router.Run(":8080")
}
```

`curl` 测试：

```shell
curl -X POST http://localhost:8080/upload \
  -F "upload[]=@/Users/appleboy/test1.zip" \
  -F "upload[]=@/Users/appleboy/test2.zip" \
  -H "Content-Type: multipart/form-data"
```

## 路由分组

```go
func main() {
    router := gin.Default()

    // Simple group: v1
    v1 := router.Group("/v1")
    {
        v1.POST("/login", loginEndpoint)
        v1.POST("/submit", submitEndpoint)
        v1.POST("/read", readEndpoint)
    }

    // Simple group: v2
    v2 := router.Group("/v2")
    {
        v2.POST("/login", loginEndpoint)
        v2.POST("/submit", submitEndpoint)
        v2.POST("/read", readEndpoint)
    }

    router.Run(":8080")
}
```

## 无中间件启动

使用

```go
r := gin.New()
```

代替

```go
// 默认启动方式，包含 Logger、Recovery 中间件
r := gin.Default()
```

## 使用中间件

```go
func main() {
    // 创建一个不包含中间件的路由器
    r := gin.New()

    // 全局中间件
    // 使用 Logger 中间件
    r.Use(gin.Logger())

    // 使用 Recovery 中间件
    r.Use(gin.Recovery())

    // 路由添加中间件，可以添加任意多个
    r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

    // 路由组中添加中间件
    // authorized := r.Group("/", AuthRequired())
    // exactly the same as:
    authorized := r.Group("/")
    // per group middleware! in this case we use the custom created
    // AuthRequired() middleware just in the "authorized" group.
    authorized.Use(AuthRequired())
    {
        authorized.POST("/login", loginEndpoint)
        authorized.POST("/submit", submitEndpoint)
        authorized.POST("/read", readEndpoint)

        // nested group
        testing := authorized.Group("testing")
        testing.GET("/analytics", analyticsEndpoint)
    }

    // Listen and serve on 0.0.0.0:8080
    r.Run(":8080")
}
```

## 写日志文件

```go
func main() {
    // 禁用控制台颜色
    gin.DisableConsoleColor()

    // 创建记录日志的文件
    f, _ := os.Create("gin.log")
    gin.DefaultWriter = io.MultiWriter(f)

    // 如果需要将日志同时写入文件和控制台，请使用以下代码
    // gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

    router := gin.Default()
    router.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })

    router.Run(":8080")
}
```

## 自定义日志格式

```go
func main() {
    router := gin.New()

    // LoggerWithFormatter 中间件会将日志写入 gin.DefaultWriter
    // By default gin.DefaultWriter = os.Stdout
    router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

        // 你的自定义格式
        return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
                param.ClientIP,
                param.TimeStamp.Format(time.RFC1123),
                param.Method,
                param.Path,
                param.Request.Proto,
                param.StatusCode,
                param.Latency,
                param.Request.UserAgent(),
                param.ErrorMessage,
        )
    }))
    router.Use(gin.Recovery())

    router.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })

    router.Run(":8080")
}
```

**输出示例：**

```
::1 - [Fri, 07 Dec 2018 17:04:38 JST] "GET /ping HTTP/1.1 200 122.767µs "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.80 Safari/537.36" "
```

## 模型绑定和验证

若要将请求主体绑定到结构体中，请使用模型绑定，目前支持JSON、XML、YAML和标准表单值(foo=bar&boo=baz)的绑定。

Gin使用 [go-playground/validator.v8](https://github.com/go-playground/validator) 验证参数，[查看完整文档](https://godoc.org/gopkg.in/go-playground/validator.v8#hdr-Baked_In_Validators_and_Tags)。

需要在绑定的字段上设置tag，比如，绑定格式为json，需要这样设置 `json:"fieldname"` 。

此外，**Gin还提供了两套绑定方法**：

- **Must bind**
- - Methods  - `Bind`, `BindJSON`, `BindXML`, `BindQuery`, `BindYAML` 
- - Behavior - 这些方法底层使用 `MustBindWith`，如果存在绑定错误，请求将被以下指令中止 `c.AbortWithError(400, err).SetType(ErrorTypeBind)`，响应状态代码会被设置为400，请求头`Content-Type`被设置为`text/plain; charset=utf-8`。注意，如果你试图在此之后设置响应代码，将会发出一个警告 `[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 400 with 422`，如果你希望更好地控制行为，请使用`ShouldBind`相关的方法
- **Should bind**
- - Methods  - `ShouldBind`, `ShouldBindJSON`, `ShouldBindXML`, `ShouldBindQuery`, `ShouldBindYAML` 
- - Behavior - 这些方法底层使用 `ShouldBindWith`，如果存在绑定错误，则返回错误，开发人员可以正确处理请求和错误。

当我们使用绑定方法时，Gin会根据Content-Type推断出使用哪种绑定器，如果你确定你绑定的是什么，你可以使用`MustBindWith`或者`BindingWith`。

给字段指定特定规则的修饰符，如果一个字段用`binding:"required"`修饰，并且在绑定时该字段的值为空，那么将返回一个错误。

```go
// 绑定为json
type Login struct {
    User     string `form:"user" json:"user" xml:"user"  binding:"required"`
    Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func main() {
    router := gin.Default()

    // Example for binding JSON ({"user": "manu", "password": "123"})
    router.POST("/loginJSON", func(c *gin.Context) {
        var json Login
        if err := c.ShouldBindJSON(&json); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        if json.User != "manu" || json.Password != "123" {
            c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
            return
        } 
        
        c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
    })

    // Example for binding XML (
    //  <?xml version="1.0" encoding="UTF-8"?>
    //  <root>
    //      <user>user</user>
    //      <password>123</password>
    //  </root>)
    router.POST("/loginXML", func(c *gin.Context) {
        var xml Login
        if err := c.ShouldBindXML(&xml); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        if xml.User != "manu" || xml.Password != "123" {
            c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
            return
        } 
        
        c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
    })

    // Example for binding a HTML form (user=manu&password=123)
    router.POST("/loginForm", func(c *gin.Context) {
        var form Login
        // This will infer what binder to use depending on the content-type header.
        if err := c.ShouldBind(&form); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        if form.User != "manu" || form.Password != "123" {
            c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
            return
        } 
        
        c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
    })

    // Listen and serve on 0.0.0.0:8080
    router.Run(":8080")
}
```

**请求示例：**

```shell
$ curl -v -X POST http://localhost:8080/loginJSON -H 'content-type: application/json' -d '{ "user": "manu" }'
> POST /loginJSON HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.51.0
> Accept: */*
> content-type: application/json
> Content-Length: 18
>
* upload completely sent off: 18 out of 18 bytes
< HTTP/1.1 400 Bad Request
< Content-Type: application/json; charset=utf-8
< Date: Fri, 04 Aug 2017 03:51:31 GMT
< Content-Length: 100
<
{"error":"Key: 'Login.Password' Error:Field validation for 'Password' failed on the 'required' tag"}
```

**跳过验证：**

当使用上面的curl命令运行上面的示例时，返回错误，因为示例中`Password`字段使用了`binding:"required"`，如果我们使用`binding:"-"`，那么它就不会报错。

## 自定义验证器

Gin允许我们自定义参数验证器，[参考1](https://github.com/gin-gonic/gin/blob/master/examples/custom-validation/server.go)，[参考2](https://github.com/go-playground/validator/releases/tag/v8.7)，[参考3](https://github.com/gin-gonic/gin/tree/master/examples/struct-lvl-validations)

```go
package main

import (
    "net/http"
    "reflect"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "gopkg.in/go-playground/validator.v8"
)

// Booking contains binded and validated data.
type Booking struct {
    CheckIn  time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
    CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn" time_format:"2006-01-02"`
}

func bookableDate(
    v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
    field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
    if date, ok := field.Interface().(time.Time); ok {
        today := time.Now()
        if today.Year() > date.Year() || today.YearDay() > date.YearDay() {
            return false
        }
    }
    return true
}

func main() {
    route := gin.Default()

    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterValidation("bookabledate", bookableDate)
    }

    route.GET("/bookable", getBookable)
    route.Run(":8085")
}

func getBookable(c *gin.Context) {
    var b Booking
    if err := c.ShouldBindWith(&b, binding.Query); err == nil {
        c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }
}
```

```shell
$ curl "localhost:8085/bookable?check_in=2018-04-16&check_out=2018-04-17"
{"message":"Booking dates are valid!"}

$ curl "localhost:8085/bookable?check_in=2018-03-08&check_out=2018-03-09"
{"error":"Key: 'Booking.CheckIn' Error:Field validation for 'CheckIn' failed on the 'bookabledate' tag"}
```

## 只绑定Get参数

`ShouldBindQuery` 函数只绑定Get参数，不绑定post数据，[查看详细信息](https://github.com/gin-gonic/gin/issues/742#issuecomment-315953017)

```go
package main

import (
    "log"

    "github.com/gin-gonic/gin"
)

type Person struct {
    Name    string `form:"name"`
    Address string `form:"address"`
}

func main() {
    route := gin.Default()
    route.Any("/testing", startPage)
    route.Run(":8085")
}

func startPage(c *gin.Context) {
    var person Person
    if c.ShouldBindQuery(&person) == nil {
        log.Println("====== Only Bind By Query String ======")
        log.Println(person.Name)
        log.Println(person.Address)
    }
    c.String(200, "Success")
}
```

## 绑定Get参数或者Post参数

[查看详细信息](https://github.com/gin-gonic/gin/issues/742#issuecomment-264681292)，这个例子很有用，可以自己实践一下

```go
package main

import (
    "log"
    "time"

    "github.com/gin-gonic/gin"
)

type Person struct {
    Name     string    `form:"name"`
    Address  string    `form:"address"`
    Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
}

func main() {
    route := gin.Default()
    route.GET("/testing", startPage)
    route.Run(":8085")
}

func startPage(c *gin.Context) {
    var person Person
    // If `GET`, only `Form` binding engine (`query`) used.
    // 如果是Get，那么接收不到请求中的Post的数据？？
    // 如果是Post, 首先判断 `content-type` 的类型 `JSON` or `XML`, 然后使用对应的绑定器获取数据.
    // See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
    if c.ShouldBind(&person) == nil {
        log.Println(person.Name)
        log.Println(person.Address)
        log.Println(person.Birthday)
    }

    c.String(200, "Success")
}
```

## 绑定uri

[查看详细信息](https://github.com/gin-gonic/gin/issues/846)

```go
package main

import "github.com/gin-gonic/gin"

type Person struct {
    ID string `uri:"id" binding:"required,uuid"`
    Name string `uri:"name" binding:"required"`
}

func main() {
    route := gin.Default()
    route.GET("/:name/:id", func(c *gin.Context) {
        var person Person
        if err := c.ShouldBindUri(&person); err != nil {
            c.JSON(400, gin.H{"msg": err})
            return
        }
        c.JSON(200, gin.H{"name": person.Name, "uuid": person.ID})
    })
    route.Run(":8088")
}
```

测试用例：

```shell
$ curl -v localhost:8088/thinkerou/987fbc97-4bed-5078-9f07-9141ba07c9f3
$ curl -v localhost:8088/thinkerou/not-uuid
```

## 绑定HTML复选框

[查看详细信息](https://github.com/gin-gonic/gin/issues/129#issuecomment-124260092)

main.go

```go
type myForm struct {
    Colors []string `form:"colors[]"`
}


func formHandler(c *gin.Context) {
    var fakeForm myForm
    c.ShouldBind(&fakeForm)
    c.JSON(200, gin.H{"color": fakeForm.Colors})
}
```

form.html

```html
<form action="/" method="POST">
    <p>Check some colors</p>
    <label for="red">Red</label>
    <input type="checkbox" name="colors[]" value="red" id="red">
    <label for="green">Green</label>
    <input type="checkbox" name="colors[]" value="green" id="green">
    <label for="blue">Blue</label>
    <input type="checkbox" name="colors[]" value="blue" id="blue">
    <input type="submit">
</form>
```

result:

```
{"color":["red","green","blue"]}
```

## 绑定Post参数

```go
package main

import (
    "github.com/gin-gonic/gin"
)

type LoginForm struct {
    User     string `form:"user" binding:"required"`
    Password string `form:"password" binding:"required"`
}

func main() {
    router := gin.Default()
    router.POST("/login", func(c *gin.Context) {
        // you can bind multipart form with explicit binding declaration:
        // c.ShouldBindWith(&form, binding.Form)
        // or you can simply use autobinding with ShouldBind method:
        var form LoginForm
        // in this case proper binding will be automatically selected
        if c.ShouldBind(&form) == nil {
            if form.User == "user" && form.Password == "password" {
                c.JSON(200, gin.H{"status": "you are logged in"})
            } else {
                c.JSON(401, gin.H{"status": "unauthorized"})
            }
        }
    })
    router.Run(":8080")
}
```

测试用例：

```shell
$ curl -v --form user=user --form password=password http://localhost:8080/login
```

## XML、JSON、YAML和ProtoBuf 渲染（输出格式）

即接口返回的数据格式

```go
func main() {
    r := gin.Default()

    // gin.H is a shortcut for map[string]interface{}
    r.GET("/someJSON", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
    })

    r.GET("/moreJSON", func(c *gin.Context) {
        // You also can use a struct
        var msg struct {
            Name    string `json:"user"`
            Message string
            Number  int
        }
        msg.Name = "Lena"
        msg.Message = "hey"
        msg.Number = 123
        // Note that msg.Name becomes "user" in the JSON
        // Will output  :   {"user": "Lena", "Message": "hey", "Number": 123}
        c.JSON(http.StatusOK, msg)
    })

    r.GET("/someXML", func(c *gin.Context) {
        c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
    })

    r.GET("/someYAML", func(c *gin.Context) {
        c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
    })

    r.GET("/someProtoBuf", func(c *gin.Context) {
        reps := []int64{int64(1), int64(2)}
        label := "test"
        // The specific definition of protobuf is written in the testdata/protoexample file.
        data := &protoexample.Test{
            Label: &label,
            Reps:  reps,
        }
        // Note that data becomes binary data in the response
        // Will output protoexample.Test protobuf serialized data
        c.ProtoBuf(http.StatusOK, data)
    })

    // Listen and serve on 0.0.0.0:8080
    r.Run(":8080")
}
```

**SecureJSON**

使用SecureJSON可以防止json劫持，如果返回的数据是数组，则会默认在返回值前加上`"while(1)"`

```
func main() {
    r := gin.Default()

    // 可以自定义返回的json数据前缀
    // r.SecureJsonPrefix(")]}',\n")

    r.GET("/someJSON", func(c *gin.Context) {
        names := []string{"lena", "austin", "foo"}

        // 将会输出:   while(1);["lena","austin","foo"]
        c.SecureJSON(http.StatusOK, names)
    })

    // Listen and serve on 0.0.0.0:8080
    r.Run(":8080")
}
```

**JSONP**

使用JSONP可以跨域传输，如果参数中存在回调参数，那么返回的参数将是回调函数的形式

```
func main() {
    r := gin.Default()

    r.GET("/JSONP", func(c *gin.Context) {
        data := map[string]interface{}{
            "foo": "bar",
        }
        
        // 访问 http://localhost:8080/JSONP?callback=call
        // 将会输出:   call({foo:"bar"})
        c.JSONP(http.StatusOK, data)
    })

    // Listen and serve on 0.0.0.0:8080
    r.Run(":8080")
}
```

**AsciiJSON**

使用AsciiJSON将使特殊字符编码

```
func main() {
    r := gin.Default()

    r.GET("/someJSON", func(c *gin.Context) {
        data := map[string]interface{}{
            "lang": "GO语言",
            "tag":  "<br>",
        }

        // 将输出: {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
        c.AsciiJSON(http.StatusOK, data)
    })

    // Listen and serve on 0.0.0.0:8080
    r.Run(":8080")
}
```

**PureJSON**

通常情况下，JSON会将特殊的HTML字符替换为对应的unicode字符，比如`<`替换为`\u003c`，如果想原样输出html，则使用PureJSON，这个特性在Go 1.6及以下版本中无法使用。

```
func main() {
    r := gin.Default()
    
    // Serves unicode entities
    r.GET("/json", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "html": "<b>Hello, world!</b>",
        })
    })
    
    // Serves literal characters
    r.GET("/purejson", func(c *gin.Context) {
        c.PureJSON(200, gin.H{
            "html": "<b>Hello, world!</b>",
        })
    })
    
    // listen and serve on 0.0.0.0:8080
    r.Run(":8080")
}
```

## 设置静态文件路径

访问静态文件需要先设置路径

```
func main() {
    router := gin.Default()
    router.Static("/assets", "./assets")
    router.StaticFS("/more_static", http.Dir("my_file_system"))
    router.StaticFile("/favicon.ico", "./resources/favicon.ico")

    // Listen and serve on 0.0.0.0:8080
    router.Run(":8080")
}
```

## 返回第三方获取的数据

```
func main() {
    router := gin.Default()
    router.GET("/someDataFromReader", func(c *gin.Context) {
        response, err := http.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
        if err != nil || response.StatusCode != http.StatusOK {
            c.Status(http.StatusServiceUnavailable)
            return
        }

        reader := response.Body
        contentLength := response.ContentLength
        contentType := response.Header.Get("Content-Type")

        extraHeaders := map[string]string{
            "Content-Disposition": `attachment; filename="gopher.png"`,
        }

        c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
    })
    router.Run(":8080")
}
```

## HTML渲染

使用`LoadHTMLGlob()` 或者 `LoadHTMLFiles()`

```
func main() {
    router := gin.Default()
    router.LoadHTMLGlob("templates/*")
    //router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
    router.GET("/index", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl", gin.H{
            "title": "Main website",
        })
    })
    router.Run(":8080")
}
```

templates/index.tmpl

```
<html>
    <h1>
        {{ .title }}
    </h1>
</html>
```

在不同目录中使用具有相同名称的模板

```
func main() {
    router := gin.Default()
    router.LoadHTMLGlob("templates/**/*")
    router.GET("/posts/index", func(c *gin.Context) {
        c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
            "title": "Posts",
        })
    })
    router.GET("/users/index", func(c *gin.Context) {
        c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
            "title": "Users",
        })
    })
    router.Run(":8080")
}
```

templates/posts/index.tmpl

```
{{ define "posts/index.tmpl" }}
<html><h1>
    {{ .title }}
</h1>
<p>Using posts/index.tmpl</p>
</html>
{{ end }}
```

templates/users/index.tmpl

```
{{ define "users/index.tmpl" }}
<html><h1>
    {{ .title }}
</h1>
<p>Using users/index.tmpl</p>
</html>
{{ end }}
```

**自定义模板渲染器**

```
import "html/template"

func main() {
    router := gin.Default()
    html := template.Must(template.ParseFiles("file1", "file2"))
    router.SetHTMLTemplate(html)
    router.Run(":8080")
}
```

**自定义渲染分隔符**

```
r := gin.Default()
    r.Delims("{[{", "}]}")
    r.LoadHTMLGlob("/path/to/templates")
```

**自定义模板函数**

[详细信息](https://github.com/gin-gonic/gin/blob/master/examples/template)

main.go

```
import (
    "fmt"
    "html/template"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

func formatAsDate(t time.Time) string {
    year, month, day := t.Date()
    return fmt.Sprintf("%d%02d/%02d", year, month, day)
}

func main() {
    router := gin.Default()
    router.Delims("{[{", "}]}")
    router.SetFuncMap(template.FuncMap{
        "formatAsDate": formatAsDate,
    })
    router.LoadHTMLFiles("./testdata/template/raw.tmpl")

    router.GET("/raw", func(c *gin.Context) {
        c.HTML(http.StatusOK, "raw.tmpl", map[string]interface{}{
            "now": time.Date(2017, 07, 01, 0, 0, 0, 0, time.UTC),
        })
    })

    router.Run(":8080")
}
```

raw.tmpl

然后就可以在html中直接使用formatAsDate函数了

```
Date: {[{.now | formatAsDate}]}
```

Result:

```
Date: 2017/07/01
```

## 多个模板文件

Gin默认情况下只允许使用一个html模板文件（即一次可以加载多个模板文件），点击[这里](https://github.com/gin-contrib/multitemplate)查看实现案例

## 重定向

发布HTTP重定向很容易，支持内部和外部链接

```
r.GET("/test", func(c *gin.Context) {
    c.Redirect(http.StatusMovedPermanently, "http://www.google.com/")
})
```

Gin路由重定向，使用如下的`HandleContext`

```
r.GET("/test", func(c *gin.Context) {
    c.Request.URL.Path = "/test2"
    r.HandleContext(c)
})
r.GET("/test2", func(c *gin.Context) {
    c.JSON(200, gin.H{"hello": "world"})
})
```

## 自定义中间件

```
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        t := time.Now()

        // Set example variable
        c.Set("example", "12345")

        // before request

        c.Next()

        // after request
        latency := time.Since(t)
        log.Print(latency)

        // access the status we are sending
        status := c.Writer.Status()
        log.Println(status)
    }
}

func main() {
    r := gin.New()
    r.Use(Logger())

    r.GET("/test", func(c *gin.Context) {
        example := c.MustGet("example").(string)

        // it would print: "12345"
        log.Println(example)
    })

    // Listen and serve on 0.0.0.0:8080
    r.Run(":8080")
}
```

## 使用BasicAuth()（验证）中间件

```
// simulate some private data
var secrets = gin.H{
    "foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
    "austin": gin.H{"email": "austin@example.com", "phone": "666"},
    "lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func main() {
    r := gin.Default()

    // Group using gin.BasicAuth() middleware
    // gin.Accounts is a shortcut for map[string]string
    authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
        "foo":    "bar",
        "austin": "1234",
        "lena":   "hello2",
        "manu":   "4321",
    }))

    // /admin/secrets endpoint
    // hit "localhost:8080/admin/secrets
    authorized.GET("/secrets", func(c *gin.Context) {
        // get user, it was set by the BasicAuth middleware
        user := c.MustGet(gin.AuthUserKey).(string)
        if secret, ok := secrets[user]; ok {
            c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
        } else {
            c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
        }
    })

    // Listen and serve on 0.0.0.0:8080
    r.Run(":8080")
}
```

## 中间件中使用Goroutines

在中间件或处理程序中启动新的Goroutines时，你不应该使用其中的原始上下文，你必须使用只读副本（`c.Copy()`）

```
func main() {
    r := gin.Default()

    r.GET("/long_async", func(c *gin.Context) {
        // 创建要在goroutine中使用的副本
        cCp := c.Copy()
        go func() {
            // simulate a long task with time.Sleep(). 5 seconds
            time.Sleep(5 * time.Second)

            // 这里使用你创建的副本
            log.Println("Done! in path " + cCp.Request.URL.Path)
        }()
    })

    r.GET("/long_sync", func(c *gin.Context) {
        // simulate a long task with time.Sleep(). 5 seconds
        time.Sleep(5 * time.Second)

        // 这里没有使用goroutine，所以不用使用副本
        log.Println("Done! in path " + c.Request.URL.Path)
    })

    // Listen and serve on 0.0.0.0:8080
    r.Run(":8080")
}
```

## 自定义HTTP配置

直接像这样使用`http.ListenAndServe()`

```
func main() {
    router := gin.Default()
    http.ListenAndServe(":8080", router)
}
```

或者

```
func main() {
    router := gin.Default()

    s := &http.Server{
        Addr:           ":8080",
        Handler:        router,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    s.ListenAndServe()
}
```

## 支持Let's Encrypt证书

1行代码实现LetsEncrypt HTTPS服务器

```
package main

import (
    "log"

    "github.com/gin-gonic/autotls"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // Ping handler
    r.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })

    log.Fatal(autotls.Run(r, "example1.com", "example2.com"))
}
```

自定义autocert管理器的示例

```
package main

import (
    "log"

    "github.com/gin-gonic/autotls"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/acme/autocert"
)

func main() {
    r := gin.Default()

    // Ping handler
    r.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })

    m := autocert.Manager{
        Prompt:     autocert.AcceptTOS,
        HostPolicy: autocert.HostWhitelist("example1.com", "example2.com"),
        Cache:      autocert.DirCache("/var/www/.cache"),
    }

    log.Fatal(autotls.RunWithManager(r, &m))
}
```

## Gin运行多个服务

请参阅[问题](https://github.com/gin-gonic/gin/issues/346)并尝试以下示例

```
package main

import (
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "golang.org/x/sync/errgroup"
)

var (
    g errgroup.Group
)

func router01() http.Handler {
    e := gin.New()
    e.Use(gin.Recovery())
    e.GET("/", func(c *gin.Context) {
        c.JSON(
            http.StatusOK,
            gin.H{
                "code":  http.StatusOK,
                "error": "Welcome server 01",
            },
        )
    })

    return e
}

func router02() http.Handler {
    e := gin.New()
    e.Use(gin.Recovery())
    e.GET("/", func(c *gin.Context) {
        c.JSON(
            http.StatusOK,
            gin.H{
                "code":  http.StatusOK,
                "error": "Welcome server 02",
            },
        )
    })

    return e
}

func main() {
    server01 := &http.Server{
        Addr:         ":8080",
        Handler:      router01(),
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    server02 := &http.Server{
        Addr:         ":8081",
        Handler:      router02(),
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    g.Go(func() error {
        return server01.ListenAndServe()
    })

    g.Go(func() error {
        return server02.ListenAndServe()
    })

    if err := g.Wait(); err != nil {
        log.Fatal(err)
    }
}
```

## 优雅重启或停止

想要优雅地重启或停止你的Web服务器，使用下面的方法

我们可以使用[fvbock/endless](https://github.com/fvbock/endless)来替换默认的`ListenAndServe`，有关详细信息，请参阅问题[＃296](https://github.com/gin-gonic/gin/issues/296)

```
router := gin.Default()
router.GET("/", handler)
// [...]
endless.ListenAndServe(":4242", router)
```

一个替换方案

-  [manners](https://github.com/braintree/manners)：一个Go HTTP服务器，能优雅的关闭
-  [graceful](https://github.com/tylerb/graceful)：Graceful是一个go的包，支持优雅地关闭http.Handler服务器
-  [grace](https://github.com/facebookgo/grace)：对Go服务器进行优雅的重启和零停机部署

如果你的Go版本是1.8，你可能不需要使用这个库，考虑使用http.Server内置的[Shutdown()](https://golang.org/pkg/net/http/#Server.Shutdown)方法进行优雅关闭，查看[例子](https://github.com/gin-gonic/gin/tree/master/examples/graceful-shutdown)

```
// +build go1.8

package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"

    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
    router.GET("/", func(c *gin.Context) {
        time.Sleep(5 * time.Second)
        c.String(http.StatusOK, "Welcome Gin Server")
    })

    srv := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    go func() {
        // service connections
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    // Wait for interrupt signal to gracefully shutdown the server with
    // a timeout of 5 seconds.
    quit := make(chan os.Signal)
    signal.Notify(quit, os.Interrupt)
    <-quit
    log.Println("Shutdown Server ...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server Shutdown:", err)
    }
    log.Println("Server exiting")
}
```

## 构建包含模板的二进制文件

你可以使用[go-assets](https://github.com/jessevdk/go-assets)将服务器构建成一个包含模板的二进制文件

```
func main() {
    r := gin.New()

    t, err := loadTemplate()
    if err != nil {
        panic(err)
    }
    r.SetHTMLTemplate(t)

    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "/html/index.tmpl",nil)
    })
    r.Run(":8080")
}

// loadTemplate loads templates embedded by go-assets-builder
func loadTemplate() (*template.Template, error) {
    t := template.New("")
    for name, file := range Assets.Files {
        if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
            continue
        }
        h, err := ioutil.ReadAll(file)
        if err != nil {
            return nil, err
        }
        t, err = t.New(name).Parse(string(h))
        if err != nil {
            return nil, err
        }
    }
    return t, nil
}
```

请参见`examples/assets-in-binary`目录中的例子

## 使用自定义结构绑定表单数据

以下示例使用自定义结构

```
type StructA struct {
    FieldA string `form:"field_a"`
}

type StructB struct {
    NestedStruct StructA
    FieldB string `form:"field_b"`
}

type StructC struct {
    NestedStructPointer *StructA
    FieldC string `form:"field_c"`
}

type StructD struct {
    NestedAnonyStruct struct {
        FieldX string `form:"field_x"`
    }
    FieldD string `form:"field_d"`
}

func GetDataB(c *gin.Context) {
    var b StructB
    c.Bind(&b)
    c.JSON(200, gin.H{
        "a": b.NestedStruct,
        "b": b.FieldB,
    })
}

func GetDataC(c *gin.Context) {
    var b StructC
    c.Bind(&b)
    c.JSON(200, gin.H{
        "a": b.NestedStructPointer,
        "c": b.FieldC,
    })
}

func GetDataD(c *gin.Context) {
    var b StructD
    c.Bind(&b)
    c.JSON(200, gin.H{
        "x": b.NestedAnonyStruct,
        "d": b.FieldD,
    })
}

func main() {
    r := gin.Default()
    r.GET("/getb", GetDataB)
    r.GET("/getc", GetDataC)
    r.GET("/getd", GetDataD)

    r.Run()
}
```

运行示例：

```
$ curl "http://localhost:8080/getb?field_a=hello&field_b=world"
{"a":{"FieldA":"hello"},"b":"world"}
$ curl "http://localhost:8080/getc?field_a=hello&field_c=world"
{"a":{"FieldA":"hello"},"c":"world"}
$ curl "http://localhost:8080/getd?field_x=hello&field_d=world"
{"d":"world","x":{"FieldX":"hello"}}
```

**注意**：不支持以下样式结构

```
type StructX struct {
    X struct {} `form:"name_x"` // HERE have form
}

type StructY struct {
    Y StructX `form:"name_y"` // HERE have form
}

type StructZ struct {
    Z *StructZ `form:"name_z"` // HERE have form
}
```

总之，现在只支持现在没有`form`标签的自定义结构

## 将请求体绑定到不同的结构体中

绑定请求体的常规方法使用`c.Request.Body`，并且不能多次调用

```
type formA struct {
  Foo string `json:"foo" xml:"foo" binding:"required"`
}

type formB struct {
  Bar string `json:"bar" xml:"bar" binding:"required"`
}

func SomeHandler(c *gin.Context) {
  objA := formA{}
  objB := formB{}
  // This c.ShouldBind consumes c.Request.Body and it cannot be reused.
  if errA := c.ShouldBind(&objA); errA == nil {
    c.String(http.StatusOK, `the body should be formA`)
  // Always an error is occurred by this because c.Request.Body is EOF now.
  } else if errB := c.ShouldBind(&objB); errB == nil {
    c.String(http.StatusOK, `the body should be formB`)
  } else {
    ...
  }
}
```

同样，你能使用`c.ShouldBindBodyWith`

```
func SomeHandler(c *gin.Context) {
  objA := formA{}
  objB := formB{}
  // This reads c.Request.Body and stores the result into the context.
  if errA := c.ShouldBindBodyWith(&objA, binding.JSON); errA == nil {
    c.String(http.StatusOK, `the body should be formA`)
  // At this time, it reuses body stored in the context.
  } else if errB := c.ShouldBindBodyWith(&objB, binding.JSON); errB == nil {
    c.String(http.StatusOK, `the body should be formB JSON`)
  // And it can accepts other formats
  } else if errB2 := c.ShouldBindBodyWith(&objB, binding.XML); errB2 == nil {
    c.String(http.StatusOK, `the body should be formB XML`)
  } else {
    ...
  }
}
```

-  `c.ShouldBindBodyWith` 在绑定之前将body存储到上下文中，这对性能有轻微影响，因此如果你要立即调用，则不应使用此方法
- 此功能仅适用于这些格式 -- `JSON`, `XML`, `MsgPack`, `ProtoBuf`。对于其他格式，`Query`, `Form`, `FormPost`, `FormMultipart`, 可以被`c.ShouldBind()`多次调用而不影响性能（参考 [#1341](https://github.com/gin-gonic/gin/pull/1341)）

## HTTP/2 服务器推送

`http.Pusher`只支持Go 1.8或更高版本，有关详细信息，请参阅[golang博客](https://blog.golang.org/h2push)

```
package main

import (
    "html/template"
    "log"

    "github.com/gin-gonic/gin"
)

var html = template.Must(template.New("https").Parse(`
<html>
<head>
  <title>Https Test</title>
  <script src="/assets/app.js"></script>
</head>
<body>
  <h1 style="color:red;">Welcome, Ginner!</h1>
</body>
</html>
`))

func main() {
    r := gin.Default()
    r.Static("/assets", "./assets")
    r.SetHTMLTemplate(html)

    r.GET("/", func(c *gin.Context) {
        if pusher := c.Writer.Pusher(); pusher != nil {
            // use pusher.Push() to do server push
            if err := pusher.Push("/assets/app.js", nil); err != nil {
                log.Printf("Failed to push: %v", err)
            }
        }
        c.HTML(200, "https", gin.H{
            "status": "success",
        })
    })

    // Listen and Server in https://127.0.0.1:8080
    r.RunTLS(":8080", "./testdata/server.pem", "./testdata/server.key")
}
```

## 自定义路由日志的格式

默认的路由日志是这样的：

```
[GIN-debug] POST   /foo                      --> main.main.func1 (3 handlers)
[GIN-debug] GET    /bar                      --> main.main.func2 (3 handlers)
[GIN-debug] GET    /status                   --> main.main.func3 (3 handlers)
```

如果你想以给定的格式记录这些信息（例如 JSON，键值对或其他格式），你可以使用`gin.DebugPrintRouteFunc`来定义格式，在下面的示例中，我们使用标准日志包记录路由日志，你可以使用其他适合你需求的日志工具

```
import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
        log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
    }

    r.POST("/foo", func(c *gin.Context) {
        c.JSON(http.StatusOK, "foo")
    })

    r.GET("/bar", func(c *gin.Context) {
        c.JSON(http.StatusOK, "bar")
    })

    r.GET("/status", func(c *gin.Context) {
        c.JSON(http.StatusOK, "ok")
    })

    // Listen and Server in http://0.0.0.0:8080
    r.Run()
}
```

## 设置并获取cookie

```
import (
    "fmt"

    "github.com/gin-gonic/gin"
)

func main() {

    router := gin.Default()

    router.GET("/cookie", func(c *gin.Context) {

        cookie, err := c.Cookie("gin_cookie")

        if err != nil {
            cookie = "NotSet"
            c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
        }

        fmt.Printf("Cookie value: %s \n", cookie)
    })

    router.Run()
}
```

## 测试

`net/http/httptest`包是http测试的首选方式

```
package main

func setupRouter() *gin.Engine {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })
    return r
}

func main() {
    r := setupRouter()
    r.Run(":8080")
}
```

测试上面的示例代码

```
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
    router := setupRouter()

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/ping", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
    assert.Equal(t, "pong", w.Body.String())
}
```

## 用户

以下是使用Gin的一些用户

-  [drone](https://github.com/drone/drone): Drone is a Continuous Delivery platform built on Docker, written in Go.
-  [gorush](https://github.com/appleboy/gorush): A push notification server written in Go.
-  [fnproject](https://github.com/fnproject/fn): The container native, cloud agnostic serverless platform.
-  [photoprism](https://github.com/photoprism/photoprism): Personal photo management powered by Go and Google TensorFlow.
-  [krakend](https://github.com/devopsfaith/krakend): Ultra performant API Gateway with middlewares.
-  [picfit](https://github.com/thoas/picfit): An image resizing server written in Go.