# 搭建HTTP服务器

## 1、Web工作方式

我们平时浏览网页的时候,会打开浏览器，输入网址后按下回车键，然后就会显示出你想要 浏览的内容。在这个看似简单的用户行为背后，到底隐藏了些什么呢？

对于普通的上网过程，系统其实是这样做的：浏览器本身是一个客户端，当你输入URL的 时候，首先浏览器会去请求DNS服务器，通过DNS获取相应的域名对应的IP，然后通过 IP地址找到IP对应的服务器后，要求建立TCP连接，等浏览器发送完HTTP Request （请求）包后，服务器接收到请求包之后才开始处理请求包，服务器调用自身服务，返回 HTTP Response（响应）包；客户端收到来自服务器的响应后开始渲染这个Response包 里的主体（body），等收到全部的内容随后断开与该服务器之间的TCP连接。 

![20180913084404299](..\image\20180913084404299.png)


Web服务器的工作原理可以简单地归纳为：

1. 客户机通过TCP/IP协议建立到服务器的TCP连接
2. 客户端向服务器发送HTTP协议请求包，请求服务器里的资源文档
3. 服务器向客户机发送HTTP协议应答包，如果请求的资源包含有动态语言的内容，那么服务器会调用动态语言的解释引擎负责处理“动态内容”，并将处理得到的数据返回给客户端
4. 客户机与服务器断开。由客户端解释HTML文档，在客户端屏幕上渲染图形结果

## 2.Go如何使得Web工作

前面小节介绍了如何通过Go搭建一个Web服务，我们可以看到简单应用一个net/http包 就方便的搭建起来了。那么Go在底层到底是怎么做的呢？

### web工作方式的几个概念

- Request：用户请求的信息，用来解析用户的请求信息，包括post、get、cookie、url等信息
- Response：服务器需要反馈给客户端的信息
- Conn：用户的每次请求链接
- Handler：处理请求和生成返回信息的处理逻辑

## 分析 http包运行机制

如下图所示，是Go实现Web服务的工作模式的流程图 

![20180913084755391](..\image\20180913084755391.png)

创建Listen Socket, 监听指定的端口, 等待客户端请求到来。

Listen Socket接受客户端的请求, 得到Client Socket, 接下来通过Client Socket与 客户端通信。

处理客户端的请求, 首先从Client Socket读取HTTP请求的协议头, 如果是POST 方法, 还可能要读取客户端提交的数据,然后交给相应的handler处理请求, handler处理完 毕准备好客户端需要的数据, 通过Client Socket写给客户端。

## 实例

### 导入依赖包

```go
import (
    "net/http"
    "io/ioutil"
    _"net/http/pprof"
)
```

定义路由处理，并监听在指定端口

```go
func main() {

//定义对/hello路由的响应
http.HandleFunc(
    "/hello",
    func(responseWriter http.ResponseWriter, request *http.Request) {
        //向响应中写入客户端地址
        responseWriter.Write([]byte(request.RemoteAddr+";"))
        //向响应中写入客户端方法
        responseWriter.Write([]byte(request.Method))
        //向响应中写入内容
        responseWriter.Write([]byte("hello 祥哥 hello go"))
    })

//定义对/sister路由的响应
http.HandleFunc(
    "/sister",
    func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("你妹"))
    })

//定义对/fuck路由的响应
http.HandleFunc("/fuck", 
                func(writer http.ResponseWriter, request *http.Request) {
    //从本地html文件中读入HTML页面的原始字节
    contentBytes, _ := ioutil.ReadFile("/home/sirouyang/Desktop/demos/W3/day5/02HTTP/fuck.html")
    //向客户端写出响应
    writer.Write(contentBytes)
})

    //开启服务器并监听在8080端口
    http.ListenAndServe("127.0.0.1:8080", nil)

}
```


