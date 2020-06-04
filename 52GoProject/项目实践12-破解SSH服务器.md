# 破解SSH服务器

## 需求和思路分析

现在的很多小伙伴们都拥有了自己的云服务器了，lots of them！
平时大家是怎么做服务器管理的呢？相信多数人都是通过SSH客户端连接过去的吧；
无论PUTTY还是XShell，我们只需要一个登陆密码，就能轻松地登陆到远程服务器终端，然后对我们的服务器做任何事情；
只需要一个密码就可以了！
Go语言有SSH连接的第三方库，参数自然是用户名、密码、远程IP和端口，而密码我们可以通过暴力枚举来进行破解；

## 建立靶机

安装并运行SSH服务，以Ubuntu为例

//安装SSH服务

```shell
sudo apt install openssh-server
```

//查看服务状态

```go
sudo systemctl status ssh
```

//启动服务

```shell
sudo systemctl start ssh
```

![20190221161620613](..\image\20190221161620613.png)

使用另一台虚拟机与之进行SSH连接

```go
ssh sriouyang@192.168.137.152
```

![20190221162030883](..\image\20190221162030883.png)

安装第三方库

```go
go get golang.org/x/crypto/ssh
```

访问golang.org需要翻墙，我们可以先手动在GOPATH/src下创建golang.org/x目录，cd到这个目录，然后执行

```go
git clone https://github.com/golang/crypto.git
```

## 连接SSH服务器

核心API如下：

```go
func Dial(network, addr string, config *ClientConfig) (*Client, error)
func (c *Client) NewSession() (*Session, error)
```

从核心API出发，一步一步按图索骥，我们就能够写出连接的完整代码了