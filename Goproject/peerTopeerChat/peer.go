package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

/*
·用一个可执行程序实现相互聊天
·实现注册节点名称，并通过名称发起会话
·实现群发消息
*/

/*
思路概要：
·节点同时具备服务端和客户端的职能
·服务端只负责接收其它节点主动发送过来的消息
·客户端只负责主动向其它节点发送消息
·通信都用短连接，服务端收完消息/客户端发完消息都断开conn——一方面是节约IO资源，另一方面是为了使逻辑清晰
·节点名称注册到【注册服务器】（很像DNS），以便根据节点名称访问节点而不是监听端口
*/

/*
节点注册服务器地址
提供节点注册和查询功能
*/
const registerAddress = "127.0.0.1:8888"

/*
节点的主业务逻辑
*/
func main() {

	//初始化缓存表
	cacheMap = make(map[string]string)

	/*从命令行接收监听端口和节点名称*/
	//从命令行上接收一个用于监听的端口:peer pa 1234
	peerName = os.Args[1]
	peerListeningPort = os.Args[2]
	fmt.Println(peerListeningPort, peerName)

	/*
	向注册器注册自己
	reg pa 1234
	*/
	peerAddress := RegOrGetPeerListeningAddress("reg " + peerName + " " + peerListeningPort)
	fmt.Println("节点注册成功", peerName, peerAddress)

	/*在独立并发任务中接收其它节点的消息*/
	go StartServe()

	/*在独立并发任务中向其它节点发送消息*/
	go StartRequest()

	//不能主协程会在此挂掉（如果主协程挂掉，子协程就跟着挂掉了）
	for {
		time.Sleep(1 * time.Second)
	}
	fmt.Println("GAME OVER")

}

/*错误处理*/
func HandleErr(err error, when string) {
	if err != nil {
		fmt.Println("err=", err, when)
		os.Exit(1)
	}
}

//节点监听端口，节点名称
var peerListeningPort, peerName string

//缓存其它通信节点的监听地址（如果已经查询过一回，就没必要每次都查询）
var cacheMap map[string]string

/*
向【注册器】注册/获取【节点的监听地址】
request 请求命令：
	reg pa 1234 向注册机注册名为pa的节点，监听在1234端口
	get pa		向注册机获取名为pa的节点的监听地址
返回值	pa节点的监听地址
*/
func RegOrGetPeerListeningAddress(request string) string {

	//拨号【注册服务器】
	conn, e := net.Dial("tcp", registerAddress)
	HandleErr(e, "RegOrGetPeerListeningAddress")

	//发送注册/查询命令
	conn.Write([]byte(request))

	//得到要注册/查询的节点的监听地址
	buffer := make([]byte, 1024)
	n, e := conn.Read(buffer)
	HandleErr(e, "RegGetPeerAddressconn.Read(buffer)")
	peerAddress := string(buffer[:n])

	//返回这个监听地址
	return peerAddress

}

/*
监听并接收其它节点发送过来的消息
这是节点【服务端】的一面
*/
func StartServe() {
	//在配置和注册好的端口建立TCP监听
	listener, e := net.Listen("tcp", ":"+peerListeningPort)
	HandleErr(e, "net.Listen")

	/*循环接入其它节点*/
	for {
		conn, e := listener.Accept()
		HandleErr(e, "listener.Accept()")

		//接收远程节点的消息
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		HandleErr(err, "conn.Read(buffer)")
		msg := string(buffer[:n])
		fmt.Println(conn.RemoteAddr(), ":", msg)

		//接收完毕立即断开
		conn.Close()
	}

}

/*
主动向其它节点发起会话
这是节点【客户端】的一面
*/
func StartRequest() {

	//目标节点名称，要发送的消息
	var targetName, msg string

	for {

		//从控制台输入信息
		fmt.Println("请输入对方名称：消息内容")
		fmt.Scan(&targetName, &msg)

		//看看缓存中是否有节点信息
		var targetAddress string
		if temp, ok := cacheMap[targetName]; !ok {
			//向注册器查询节点的监听地址
			fmt.Println("从注册服务器获得节点监听地址")
			targetAddress = RegOrGetPeerListeningAddress("get " + targetName)

			//将查询结果写入缓存
			cacheMap[targetName] = targetAddress

		} else {
			//使用缓存中的监听地址
			fmt.Println("从缓存获得节点监听地址")
			targetAddress = temp
		}

		//向目标地址发送消息
		conn, e := net.Dial("tcp", targetAddress)
		HandleErr(e, "net.Dial")
		conn.Write([]byte(msg))

		//消息发送完毕，断开连接
		conn.Close()
	}

}


