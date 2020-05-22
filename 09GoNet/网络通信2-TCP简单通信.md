# TCP简单通信

## 首部格式

图释: 

![20180912234014108](..\image\20180912234014108.png)

各个段位说明:

- 源端口和目的端口:各占 2 字节.端口是传输层与应用层的服务接口.传输层的复用和分用功能都要通过端口才能实现
- 序号:占 4 字节.TCP 连接中传送的数据流中的每一个字节都编上一个序号.序号字段的值则指的是本报文段所发送的数据的第一个字节的序号
- 确认号:占 4 字节,是期望收到对方的下一个报文段的数据的第一个字节的序号
- 数据偏移/首部长度:占 4 位,它指出 TCP 报文段的数据起始处距离 TCP 报文段的起始处有多远.“数据偏移”的单位是 32位字(以 4 字节为计算单位)
- 保留:占 6 位,保留为今后使用,但目前应置为 0
- 紧急URG:当 URG=1 时,表明紧急指针字段有效.它告诉系统此报文段中有紧急数据,应尽快传送(相当于高优先级的数据)
- 确认ACK:只有当 ACK=1 时确认号字段才有效.当 ACK=0 时,确认号无效
- PSH(PuSH):接收 TCP 收到 PSH = 1 的报文段,就尽快地交付接收应用进程,而不再等到整个缓存都填满了后再向上交付
- RST (ReSeT):当 RST=1 时,表明 TCP连接中出现严重差错（如由于主机崩溃或其他原因）,必须释放连接,然后再重新建立运输连接
- 同步 SYN:同步 SYN = 1 表示这是一个连接请求或连接接受报文
- 终止 FIN:用来释放一个连接.FIN=1 表明此报文段的发送端的数据已发送完毕,并要求释放运输连接
- 检验和:占 2 字节.检验和字段检验的范围包括首部和数据这两部分.在计算检验和时,要在 TCP 报文段的前面加上 12 字节的伪首部
- 紧急指针:占 16 位,指出在本报文段中紧急数据共有多少个字节（紧急数据放在本报文段数据的最前面）
- 选项:长度可变.TCP最初只规定了一种选项,即最大报文段长度 MSS.MSS 告诉对方TCP：“我的缓存所能接收的报文段的数据字段的最大长度是 MSS 个字节.” [MSS(Maximum Segment Size)是TCP 报文段中的数据字段最大长度.数据字段加上 TCP 首部才等于整个的 TCP 报文段]
- 填充:这是为了使整个首部长度是 4 字节的整数倍
- 其他选项: 
  - (1).窗口扩大:占 3 字节,其中有一个字节表示移位值 S.新的窗口值等于TCP 首部中的窗口位数增大到(16 + S),相当于把窗口值向左移动 S 位后获得实际的窗口大小 
  - (2).时间戳:占10 字节,其中最主要的字段时间戳值字段(4字节)和时间戳回送回答字段(4字节) 
  - (3).选择确认:接收方收到了和前面的字节流不连续的两2字节.如果这些字节的序号都在接收窗口之内,那么接收方就先收下这些数据,但要把这些信息准确地告诉发送方,使发送方不要再重复发送这些已收到的数据
    数据单位

**TCP 传送的数据单位协议是 TCP 报文段(segment)**

特点

TCP 是面向连接的传输层协议 每一条 TCP 连接只能有两个端点(endpoint),每一条 TCP 连接只能是点对点的（一对一） TCP 提供可靠交付的服务 TCP 提供全双工通信 面向字节流

## 实例

### 服务端实现

```go
import (
    "fmt"
    "net"
    "os"
)

func CheckErrorS(err error){
    if  err !=nil{
        fmt.Println("网络错误",err.Error())
        os.Exit(1)
    }
}

func  Processinfo(conn net.Conn){
    buf :=make([]byte,1024)//开创缓冲区
    defer conn.Close() //关闭连接

    //与客户端源源不断地IO
    for{
        n,err:=conn.Read(buf) //读取数据
        if err !=nil{
            break
        }

        if n !=0{
            msg := string(buf[:n])
            fmt.Println("收到消息",msg)
            conn.Write([]byte("已阅："+msg))

            if msg=="分手吧"{
                break
            }
        }
	}

}

func main() {
    //建立TCP服务器
    listener,err:=net.Listen("tcp","127.0.0.1:8898")
    defer listener.Close() //关闭网络
    CheckErrorS(err)

    //循环接收，来者不拒
    for  {
        //接入一个客户端
        conn,err:= listener.Accept() //新的客户端连接
        CheckErrorS(err)

        //为每一个客户端开一条独立的协程与其IO
        go  Processinfo(conn)
    }

}
```

### 客户端实现

```go
import (
    "fmt"
    "net"
    "os"
    "bufio"
)

func CheckErrorC(err error){
    if  err !=nil{
        fmt.Println("网络错误",err.Error())
        os.Exit(1)
    }
}

func main() {
    conn,err :=net.Dial("tcp","127.0.0.1:8898")//建立TCP服务器
    defer conn.Close() //延迟关闭网络连接
    CheckErrorC(err)//检查错误

    //创建一个黑窗口（标准输入）的读取器
    reader := bufio.NewReader(os.Stdin)
    buffer := make([]byte, 1024)

    for{
        lineBytes, _, err := reader.ReadLine()
        CheckErrorC(err)
        conn.Write(lineBytes)
        fmt.Println("发送消息")

        /*接收消息*/
        n, err := conn.Read(buffer)
        CheckErrorC(err)

        msg := string(buffer[:n])
        fmt.Println("收到服务端消息：",msg)
    }


}
```


