# Go与Redis的交互

## Redis简介

Redis是一个开源的、使用C语言编写的、支持网络交互的、可基于内存也可持久化的Key-Value数据库。

特点

- 支持更多数据类型

和Memcached类似，它支持存储的value类型相对更多，包括string(字符串)、list(链表)、set(集合)、zset(sorted set 有序集合)和hash（哈希类型）。

- 支持复杂操作

这些数据类型都支持push/pop、add/remove及取交集并集和差集及更丰富的操作，而且这些操作都是原子性的。在此基础上，Redis支持各种不同方式的排序。

- 支持主从同步

与memcached一样，为了保证效率，数据都是缓存在内存中。区别的是Redis会周期性的把更新的数据写入磁盘或者把修改操作写入追加的记录文件，并且在此基础上实现了master-slave(主从)同步。数据可以从主服务器向任意数量的从服务器上同步，从服务器可以是关联其他从服务器的主服务器。这使得Redis可执行单层树复制。从盘可以有意无意的对数据进行写操作。由于完全实现了发布/订阅机制，使得从数据库在任何地方同步树时，可订阅一个频道并接收主服务器完整的消息发布记录。同步对读取操作的可扩展性和数据冗余很有帮助。

Redis的出现，很大程度补偿了memcached这类key/value存储的不足，在部 分场合可以对关系数据库起到很好的补充作用。它提供了Java，C/C++，C#，PHP，JavaScript，Perl，Object-C，Python，Ruby，Erlang等客户端，使用很方便。Redis的官网地址，非常好记，是redis.io。目前，Vmware在资助着Redis项目的开发和维护。

## 导入包

```
import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"sync"
	"time"
)
```

## 连接数据库

连接数据库

```go
func main11() {

    //连接本地redis
    conn, err := redis.Dial("tcp", "localhost:6379")
    if err != nil {
        fmt.Println("conn redis failed,", err)
        return
    }
    fmt.Println("connect ok!")

    //关闭连接
    defer conn.Close()

}
```

## 键值对操作set,get

```go
func main12() {
	//连接本地redis
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("conn redis failed,", err)
		return
	}

    //延时关闭redis连接
    defer conn.Close()

    //设置键值abc=100，等同set abc 100
    reply, err := conn.Do("Set", "abc", 100)
    fmt.Printf("reply=%#v,err=%v\n", reply, err)
    if err != nil {
        fmt.Println(err)
        return
    }

    //执行get abc，将结果转化为int
    reply, err = conn.Do("Get", "abc")
    fmt.Printf("reply=%#v,err=%v\n", reply, err)
    r, err := redis.Int(reply, err)
    if err != nil {
        fmt.Println("get abc failed,", err)
        return
    }

    //打印结果100
    fmt.Println(r)

}
```

## 批量存取mset和mget

```go
func main13() {

    //连接本地redis
    conn, err := redis.Dial("tcp", "localhost:6379")
    if err != nil {
        fmt.Println("conn redis failed,", err)
        return
    }

    //延时关闭连接
    defer conn.Close()

    //存入两个键值
    cdoReply, err := conn.Do("MSet", "abc", 100, "efg", 300)
    fmt.Printf("reply=%#v,err=%v\n", cdoReply, err)
    if err != nil {
        fmt.Println(err)
        return
    }
    //cdoReply= OK
    fmt.Println("cdoReply=", cdoReply)

    //获取存入的键值
    cdoReply, err = conn.Do("MGet", "abc", "efg")
    fmt.Printf("reply=%#v,err=%v\n", cdoReply, err)
    retInts, err := redis.Ints(conn.Do("MGet", "abc", "efg"))
    if err != nil {
        fmt.Println("get abc failed,", err)
        return
    }

    //打印获取到的结果
    for i, v := range retInts {
        fmt.Println(i, "=", v)
    }

}
```

## 设置数据过期时间

```go
func main14() {

    //连接本地redis
    c, err := redis.Dial("tcp", "localhost:6379")
    if err != nil {
        fmt.Println("conn redis failed,", err)
        return
    }

    //延时关闭
    defer c.Close()

    //设置abc过期
    cdoReply, err := c.Do("expire", "abc", 1000)
    fmt.Printf("reply=%#v,err=%v\n", cdoReply, err)
    if err != nil {
        fmt.Println(err)
        return
    } else {
        fmt.Println(cdoReply)
        fmt.Println("abc将于1000秒后过期！")
    }

}
```

## 哈希（对象）操作

```go
func main15() {

    //连接本地redis
    conn, err := redis.Dial("tcp", "localhost:6379")
    if err != nil {
        fmt.Println("conn redis failed,", err)
        return
    }

    //延时关闭
    defer conn.Close()

    //设置books下的键值:abc=100
    _, err = conn.Do("HSet", "books", "price", 56.78)
    if err != nil {
        fmt.Println(err)
        return
    }

    //获取books下的键值abc，并转化为int
    //hget books abc
    r, err := redis.Float64(conn.Do("HGet", "books", "price"))
    if err != nil {
        fmt.Println("get abc failed,", err)
        return
    }

    //打印结果
    fmt.Println(r)

}
```

## list相关操作

```go
func main16() {
	//连接本地redis
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("conn redis failed,", err)
		return
	}

    //延时关闭
    defer c.Close()

    //向一个键为books_list的列表中追加3个元素
    //lpush book_list abc ceg 300
    _, err = c.Do("lpush", "book_list", "abc", "ceg", 300)
    if err != nil {
        fmt.Println(err)
        return
    }

    //弹出books_list中最左侧的元素：300，转化为字符串
    //lpop book_list
    r, err := redis.String(c.Do("lpop", "book_list"))
    if err != nil {
        fmt.Println("get abc failed,", err)
        return
    }

    //输出结果
    fmt.Println(r)

}
```

## sort-set相关操作

```go
func main17() {
	conn, e := redis.Dial("tcp", "localhost:6379")
	if e != nil {
		fmt.Println("e=", e)
		return
	}
	defer conn.Close()

reply, err := conn.Do("zadd", "gozset", 10, "bill", 5, "jack", 7, "jobs")
fmt.Printf("reply=%#v,err=%v\n", reply, err)

}
```

## 复杂命令举例：zset求交集

```go
func main18() {
	conn, e := redis.Dial("tcp", "localhost:6379")
	if e != nil {
		fmt.Println("e=", e)
		return
	}
	defer conn.Close()

    reply, err := conn.Do("zadd", "comezset", 10, "bill", 5, "jack", 7, "mark")
    fmt.Printf("reply=%#v,err=%v\n", reply, err)

    reply, err = conn.Do("ZINTERSTORE", "interZset", 2, "gozset", "comezset")
    fmt.Printf("reply=%#v,err=%v\n", reply, err)

}

```

## 连接池

```go
var wg sync.WaitGroup
func ConnSetName(pool *redis.Pool, i int) {
	conn := pool.Get()
	defer conn.Close()
	reply, _ := conn.Do("set", "name"+strconv.Itoa(i), "123")
	fmt.Printf("conn=%v reply=%v\n", conn, reply)
	<-time.After(3 * time.Second)
	wg.Done()
}

func main() {

    //定义redis连接池指针poot
    //var pool *redis.Pool

    //新建连接池对象，指针丢给pool
    pool := &redis.Pool{
        //最大闲置连接数
        MaxIdle: 16,

    	//最大活动连接数，0=无限
    	MaxActive: 0,

    	//闲置超时时间：300秒，超时自动断开连接
    	IdleTimeout: 300,

    	//连接函数
    	Dial: func() (redis.Conn, error) {
    		//返回连接对象(或错误信息)
    		return redis.Dial("tcp", "localhost:6379")
    	},
    }
    defer pool.Close()

    for i := 0; i < 50; i++ {
        wg.Add(1)
        go ConnSetName(pool, i)
    }
    wg.Wait()

}

```

