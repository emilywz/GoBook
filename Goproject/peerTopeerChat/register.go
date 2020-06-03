package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

/*
负责注册节点名称：节点运行端口
*/
func RHandleErr(err error, when string)  {
	if err != nil{
		fmt.Println("注册器err=",err,when)
		os.Exit(1)
	}
}

//所有节点【名称-监听地址】映射表
var peerNameListeningAddressMap map[string]string

/*注册机的主业务*/
func main() {

	//初始化注册表
	peerNameListeningAddressMap = make(map[string]string)

	//开启注册服务
	listener, e := net.Listen("tcp", ":8888")
	RHandleErr(e,"net.Listen")

	buffer := make([]byte, 1024)

	/*循环接受节点的注册和查询服务*/
	for  {
		conn, e := listener.Accept()
		RHandleErr(e,"listener.Accept()")

		/*
		接收节点消息
		reg pa 1234 	注册：节点名称pa，监听端口1234
		get pa 			查询：节点名称为pa的节点的监听地址
		*/
		n, e := conn.Read(buffer)
		RHandleErr(e,"conn.Read(buffer)")
		msg := string(buffer[:n])

		//将消息炸碎为字符串，获取命令和节点名称
		strs := strings.Split(msg, " ")
		cmd := strs[0]
		peerName := strs[1]

		if cmd == "reg"{
			//将节点名称和【节点-地址】写入全局映射表 reg pa 1234
			runningAddress := conn.RemoteAddr().String()

			//拼接节点的IP和监听端口,得到节点的监听地址
			peerIP := strings.Split(runningAddress, ":")[0]
			peerListeningPort := strs[2]
			listenAddress := peerIP +":" + peerListeningPort

			//节点名称为键，监听地址为值，写入map（下次就可以供别人查询了）
			peerNameListeningAddressMap[peerName] = listenAddress

			//将节点的名称和监听地址写入全局映射表
			conn.Write([]byte(listenAddress))

		} else if cmd=="get"{

			//根据节点名称查询节点监听地址
			listeningAddress := peerNameListeningAddressMap[peerName]
			conn.Write([]byte(listeningAddress))

		}

		//断开会话，继续接受其他节点的注册或查询请求
		conn.Close()
	}

}
