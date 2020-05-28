# sync

## 同步调度概述

- Go语言的并发中，当使用go关键字开辟若干新的协程时，如果不加干涉，它们会完全并发地得到执行；
- 而所谓调度，就是在并发的局部植入一些串行和同步的操作，让某些协程有逻辑上的先后关系；
- 执行同步调度，可以通过管道读写阻塞和Go语言SDK提供的sync包下的API两种方式；
- Go语言的官方文档推荐我们使用管道的读写阻塞功能来执行这些同步操作，并称使用sync包为“低水平编程”

```
package sync

import "sync"

sync 包提供了基本的同步基元，如互斥锁。除了once和WaitGroup类型，大部分都是适用于低水平程序线程，高水平的同步使用channel通信更好一些
```

而事实上，sync包下的许多功能还是很好用的，而且有助于传统并发代码做业务上的无缝转换，例如使用Go语言重构Java工程；

## 本例需求摘要

实现一个可以边写边读的文件类

以指定大小的数据块作为基本读写单位，即一块一块地进行读写

读写成功时返回当前读写的数据块序列号

## 本例技术栈

- 等待组 sync.WaitGroup
- 条件变量 sync.Cond
- 原子操作 sync/atomic
- 读写锁 sync.RWMutex

## 实现思路摘要

需要通过roffset和woffset属性来记录最后一次读写的字节偏移量，以便下一次读写时进行定位；
通过【上一次读写的字节偏移量】/【数据块长度】，我们可以得到当前读写的数据块的序列号；
边写边读需要监听读取是否到已到文件末尾，如果已到文件末尾，则【读取协程】需要阻塞等待至【写入协程】向其发送【数据更新通知】，这里我们可以使用【条件变量的通知机制】，即cond.Wait()和cond.Signal()/Broadcast()加以实现；
文件只允许一写多读，我们使用读写锁实现；
更新文件的读写偏移量时，一定要杜绝并发操作，保证并发安全，我们使用原子操作来替代资源锁；
使用等待组使主程序恰好结束在读写都结束的时候；
完整实现

```go
import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

/*给原始字节切片起一个类别名*/
type Data []byte

type DataFile struct {
	//文件
	f *os.File

    //文件的读写锁
    fmutex sync.RWMutex

    //数据块长度
    datalen uint32

    //最后一次写入的偏移量（字节）
    woffset int64
    //最后一次读取的偏移量（字节）
    roffset int64

    //条件变量
    rcond *sync.Cond

}

/*
工厂方法
参数：
path string		文件路径

datalen uint32	指定数据块大小

返回值：
DataFile	数据文件对象
error		错误
*/

func NewDataFile(path string, datalen uint32) (*DataFile, error) {
	//得到一个可读可写（覆盖式写入）的文件
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

    /*数据块大小不能为0*/
    if datalen == 0 {
        return nil, errors.New("Invalid data legth!")
    }

    //创建指定的myDataFile对象，并指定IO文件和数据块大小
    df := &DataFile{f: f, datalen: datalen}

    //初始化df对象的条件变量
    df.rcond = sync.NewCond(df.fmutex.RLocker())
    return df, nil

}

/*返回下一次要读取的数据块的序列号*/
func (df *DataFile) Rsn() int64 {
	//同步加载最后一次读取的字节偏移量
	offset := atomic.LoadInt64(&df.roffset)

    //返回下一次读取的数据块的序列号
    return offset / int64(df.datalen)

}

/*返回下一次要写入的数据块的序列号*/
func (df *DataFile) Wsn() int64 {
	//同步获取当前df的写入的字节偏移量
	offset := atomic.LoadInt64(&df.woffset)

    // 返回下一次写入的数据块的序列号
    return offset / int64(df.datalen)

}

/*返回块文件的大小*/
func (df *DataFile) DataLen() uint32 {
	//return df.datalen

    //同步获取当前df的数据块大小并返回
    return atomic.LoadUint32(&df.datalen)

}

/*
可以读，返回值：
rsn 	read-serial-number 当前读取到的【数据块序列号】
data	原始字节数据
err 	错误
*/
func (df *DataFile) Read() (rsn int64, data Data, err error) {

    var offset int64
    //offset = df.roffset
    //加载最后一次读取的字节偏移量
    offset = atomic.LoadInt64(&df.roffset)

    //计算本次读取的数据块的序列号
    rsn = offset / int64(df.datalen) //第几块数据块

    //创建一个数据块大小的缓冲区
    buffer := make([]byte, df.datalen)

    //加读锁
    df.fmutex.RLock()
    fmt.Println("Read get lock")
    defer func() {
        //返回之前释放读锁
        df.fmutex.RUnlock()
        fmt.Println("Read release lock")
    }()

    for {
        //从指定的字节偏移量处进行读取
        _, err = df.f.ReadAt(buffer, offset)
        //fmt.Println("n,err=",n,e)
        if err != nil {

            //如果已经读到了文件末尾
            if err == io.EOF {
                fmt.Println("eof:Read release lock")
                //阻塞等待有新的内容写入
                df.rcond.Wait()
                fmt.Println("eof:Read get lock")
                continue
            }
            return
        }

        //正常地读到数据，并返回
        data = buffer

        //通过原子操作，让最后一次读取的字节偏移量+=df.datalen
        atomic.AddInt64(&df.roffset, int64(df.datalen))

        return
    }

}

func (df *DataFile) Write(data Data) (wsn int64, err error) {
	//获取并更新写偏移量
	var offset int64
	//offset = df.woffset
	//加载上次写入的字节偏移量
	offset = atomic.LoadInt64(&df.woffset)

    /*如果要写入的数据超过一个数据块的长度，就进行截取操作，否则直接使用*/
    var buffer []byte
    if len(data) > int(df.datalen) {
        buffer = data[0:df.datalen]
    } else {
        buffer = data
    }

    /*加写锁准备进行写入*/
    df.fmutex.Lock()
    fmt.Println("Write get lock")
    defer func() {
        //操作完毕释放写锁
        df.fmutex.Unlock()
        fmt.Println("Write release lock")
    }()

    /*进行数据写入*/
    //写入数据，并向读取协程发送通知信号
    _, err = df.f.WriteAt(buffer, offset)

    if err == nil{
        //本次写入导致最后一次写入的字节偏移量+=df.datalen
        atomic.AddInt64(&df.woffset, int64(df.datalen))

        //计算本次写入的块序列号
        wsn = offset / int64(df.datalen)

        //向读取协程发送数据更新通知
        fmt.Println("Write signal")
        df.rcond.Signal()
    }

    return

}
```

## 测试效果

这里我们开辟一读一写两条子协程
当读取协程读取不到数据时就进入阻塞等待，直到被写入协程的内容更新通知唤醒

```
func main() {
	var wg sync.WaitGroup

    df, e := NewDataFile(`D:\BJBlockChain1804\code\W5\day1\01昨日作业\示例文档.txt`, 15)
    fmt.Println(df, e)

    //fmt.Println("数据块大小=", df.DataLen())  //3
    //fmt.Println("下一次读取的块序列号=", df.Rsn()) //0
    //fmt.Println("下一次写入的块序列号=", df.Wsn()) //0

    wg.Add(2)
    /*写入3个数据块*/
    go func() {
        wsn, err := df.Write(Data("明月几时有"))
        fmt.Println("wsn, err=",wsn,err)

    	wsn, err = df.Write(Data("把酒问青天"))
    	fmt.Println("wsn, err=",wsn,err)

    	wsn, err = df.Write(Data("不知天上宫阙"))
    	fmt.Println("wsn, err=",wsn,err)

    	wg.Done()
    }()

    /*读取3个数据块*/
    go func() {
        rsn, data, err := df.Read()
        fmt.Println("rsn, data, err=",rsn,string(data),err)

    	rsn, data, err = df.Read()
    	fmt.Println("rsn, data, err=",rsn,string(data),err)

    	rsn, data, err = df.Read()
    	fmt.Println("rsn, data, err=",rsn,string(data),err)
    	wg.Done()
    }()

    wg.Wait()
    fmt.Println("main over")

}

```


