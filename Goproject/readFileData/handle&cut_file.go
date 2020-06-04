package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	//"sync"
)

var (
	//全局数据管道
	ch chan string

)

func main() {

	/*主协程读入数据，将不同省份的记录丢入对应的管道*/
	file, _ := os.Open(`..\dataSet\kaifang_good.txt`)
	defer file.Close()

	//初始化数据管道
	//ch = make(chan []byte)
	ch = make(chan string)

	/*并发10个文件写入协程*/
	for i := 0; i < 10; i++ {
		wg.Add(1)

		/*协程任务：从管道中拉取数据并写入到文件中*/
		go func(indx int) {
			f, err := os.OpenFile(`..\dataSet\kaifang_good_`+strconv.Itoa(indx)+`.txt`, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			HandleError(err, `os.OpenFile`)
			defer f.Close()

			totalLines := 0
			for lineStr := range ch {
				//向文件中写出UTF-8字符串
				f.WriteString(lineStr)

				//统计写出数量
				totalLines ++
				log.Printf("协程%d写入：%d", indx, totalLines)
			}
			wg.Done()
		}(i)
	}

	//创建缓冲读取器
	reader := bufio.NewReader(file)
	totalLines := 0
	for {
		//读取一行字符串（编码为UTF-8）
		lineStr, err := reader.ReadString('\n')
		totalLines ++
		println("读取数据：", totalLines)

		//读取完毕时，关闭所有数据管道，并退出读取
		if err == io.EOF {
			fmt.Println("已经读到文件末尾！")
			close(ch)
			break
		}

		ch <- lineStr
	}

	//阻塞等待所有协程结束任务
	wg.Wait()
	fmt.Println("main over!")

}

/*处理错误*/
func HandleError(err error, when string) {
	if err != nil {
		log.Fatal(err, when)
	}
}
