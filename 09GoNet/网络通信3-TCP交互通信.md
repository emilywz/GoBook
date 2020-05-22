# TCP交互通信

## 服务端实现

```go
import (
    "fmt"
    "net"
    "os"
    "strings"
)

func CheckErrorS(err error) {
    if err != nil {
        fmt.Println("网络错误", err.Error())
        os.Exit(1)
    }
}

func Processinfo(conn net.Conn) {
    buffer := make([]byte, 1024) //开创缓冲区
    defer conn.Close()           //关闭连接

    for {
        n, err := conn.Read(buffer) //读取数据
        CheckErrorS(err)

        if n != 0 {
            //拿到客户端地址
            remoteAddr := conn.RemoteAddr()
            msg := string(buffer[:n])
            fmt.Println("收到消息",msg, "来自", remoteAddr)

            if strings.Contains(msg,"钱") {
                conn.Write([]byte("fuckoff"))
                break
            }
            conn.Write([]byte("已阅："+msg))
        }
    }

}

func main() {
    //建立TCP服务器
    listener, err := net.Listen("tcp", "127.0.0.1:8898")
    CheckErrorS(err)
    defer listener.Close() //关闭网络

    fmt.Println("服务器正在等待")

    for {
        conn, err := listener.Accept() //新的客户端连接
        CheckErrorS(err)

        //处理每一个客户端
        go Processinfo(conn)
    }

}
```



## 客户端实现

```go
import (
    "fmt"
    "net"
    "bufio"
    "os"
)

func CheckErrorC(err error) {
    if err != nil {
        fmt.Println("网络错误", err.Error())
        os.Exit(1)
    }
}

func MessageSend(conn net.Conn) {
    var msg string
    reader := bufio.NewReader(os.Stdin) //读取键盘输入

    for {
        lineBytes, _, _ := reader.ReadLine() //读取一行
        msg = string(lineBytes)              //键盘输入转化为字符串

        if msg == "exit" {
            conn.Close()
            fmt.Println("客户端关闭")
            break
        }

        _, err := conn.Write([]byte(msg)) //输入写入字符串
        if err != nil {
            conn.Close()
            fmt.Println("客户端关闭")
            break
        }

    }

}

func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:8898") //建立TCP服务器
    CheckErrorC(err)                               //检查错误
    defer conn.Close()

    //发送消息中有阻塞读取标准输入的代码
    //为了避免阻塞住消息的接收，所以把它独立的协程中
    go MessageSend(conn)

    buffer := make([]byte, 1024)
    for {
        n, err := conn.Read(buffer)
        CheckErrorC(err)

        msg := string(buffer[:n])
        fmt.Println("收到服务器消息", msg)

        if msg=="fuckoff"{
            break
        }

    }

    fmt.Println("连接已断开")

}
```

