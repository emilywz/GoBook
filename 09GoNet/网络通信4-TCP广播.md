# TCP广播

## 服务端实现

```go
import (
    "fmt"
    "net"
        "os"
    "strings"
)

//存储客户端连接， key,ip端口，value 链接对象
var onlineConnsMap = make(map[string]net.Conn)

//消息队列，缓冲区
var messageQueue = make(chan string, 1000)

//消息，处理程序退出
var quitchan = make(chan bool)

func CheckErrorS(err error) {
    if err != nil {
        fmt.Println("网络错误", err.Error())
        os.Exit(1)
    }
}

func Processinfo(conn net.Conn) {
    buf := make([]byte, 1024) //开创缓冲区
    defer func() {
        //关闭连接
        conn.Close()
        fmt.Println("连接已关闭:",conn.RemoteAddr())
    }()

    for {
        n, err := conn.Read(buf) //读取数据
        if err != nil {
            fmt.Println("读取客户端消息出错,err=",err)
            break
        }

        if n != 0 {
            //消息处理，
            message := string(buf[0:n])
            fmt.Printf("受到客户端%v消息：%s\n",conn.RemoteAddr(),message)

            //消息队列存储消息
            messageQueue <- message
            fmt.Println("该条消息已加入队列...")

        }else{
            fmt.Println("读取客户端消息长度为0")
        }
    }

}

//消息的协程
func consumeMessage() {
    for {
        select {
        case message := <-messageQueue: //取出消息
            fmt.Println("已从队列中取出消息")
            strs := strings.Split(message, "#") //字符串切割
            if len(strs) > 1 {

                //127.0.0.1:12345#去死

                //截取地址并裁减头尾空白
                addr := strs[0]
                addr = strings.TrimSpace(addr)

                //截取内容
                msg := strs[1]

                //给所有客户端群发此消息
                /*for addr, conn := range onlineConnsMap { //循环在线列表
                    _, err := conn.Write([]byte(msg))
                    fmt.Println("服务器发送消息", msg, "给", addr)
                    if err != nil {
                        fmt.Println("在线链接发送失败")
                    }
                }*/

                //@某人，发消息
                if conn, ok := onlineConnsMap[addr]; ok {
                    _, err := conn.Write([]byte(msg))
                    fmt.Println("服务器发送消息", msg)
                    if err != nil {
                        fmt.Println("在线链接发送失败")
                    }
                }


            }
        case <-quitchan: //处理退出
            os.Exit(0)
        }
    }
}

func main() {
    //建立TCP服务器
    listener, err := net.Listen("tcp", "127.0.0.1:8898")
    CheckErrorS(err)
    defer listener.Close() //关闭网络

    fmt.Println("服务器正在等待")

    go consumeMessage()

    for {
        conn, err := listener.Accept() //新的客户端连接
        CheckErrorS(err)

        //处理每一个客户端
        addr := fmt.Sprint(conn.RemoteAddr()) //取出地址
        onlineConnsMap[addr] = conn

        fmt.Println("客户端列表")
        fmt.Println("-------------------")
        for key := range onlineConnsMap { //循环每一个链接
            fmt.Println(key)
        }

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
    }
}

func MessageSend(conn net.Conn) {
    var input string
    for {
        reader := bufio.NewReader(os.Stdin) //读取键盘输入
        data, _, _ := reader.ReadLine()     //读取一行
        input = string(data)                //键盘输入转化为字符串

        if input == "exit" {
            conn.Close()
            fmt.Println("客户端关闭")
            break
        }

        _, err := conn.Write([]byte(input)) //输入写入字符串
        fmt.Println("发送消息", input)
        if err != nil {
            conn.Close()
            fmt.Println("客户端关闭")
            break
        }

    }
}

func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:8898") //建立TCP服务器
    defer conn.Close()                             //延迟关闭网络连接
    CheckErrorC(err)                               //检查错误

    go MessageSend(conn)                           //开启一个协程,处理发送
    //conn.Write([]byte("hello  nimei"))

    //协程,负责收取消息
    buf := make([]byte, 1024)
    for {
        numOfBytes, err := conn.Read(buf)
        CheckErrorC(err)
        fmt.Println("收到服务器消息", string(buf[:numOfBytes]))
    }

    fmt.Println("客户端关闭")

}
```

