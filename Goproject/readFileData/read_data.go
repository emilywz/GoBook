package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"github.com/axgle/mahonia"
	"strings"
)

//利用缓冲式读取
func main() {

	//打开文件
	file, _ := os.Open(`.\dataSet\kaifangX.txt`)
	defer file.Close()

	//创建文件的缓冲读取器
	reader := bufio.NewReader(file)

	for {

		//逐行读取
		lineBytes, _, err := reader.ReadLine()

		//读到文件末尾时退出循环
		if err == io.EOF {
			break
		}

		//输出当前行内容
		gbkStr := string(lineBytes)
		utfStr, _ := ConvertEncoding(gbkStr, "GBK", "UTF-8")
		fmt.Println(utfStr)
	}

}

func main021() {

	//打开原始文本
	srcFile, _ := os.Open(`.\dataSet\kaifangX.txt`)
	defer srcFile.Close()

	//准备一个优质数据文件
	goodFile, _ := os.OpenFile(`.\dataSet\kaifang_good.txt`, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer goodFile.Close()

	//准备一个不良数据文件
	badFile, _ := os.OpenFile(`.\dataSet\kaifang_bad.txt`, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer badFile.Close()

	//创建原始文件的缓冲式读取器
	reader := bufio.NewReader(srcFile)

	//循环读取
	for {

		//一次读取一行数据
		lineBytes, _, err := reader.ReadLine()

		//读到文件末尾时退出循环
		if err == io.EOF {
			break
		}

		//转码得到当前行的UTF-8字符串
		gbkStr := string(lineBytes)
		lineStr, _ := ConvertEncoding(gbkStr, "GBK", "UTF-8")
		//fmt.Println(lineStr)

		//使用逗号将当前行打散为子串
		fields := strings.Split(lineStr, ",")

		//如果第二个子串的长度为18位，我们认为它是一个合法的身份证信息
		if len(fields) > 1 && len(fields[1])==18{
			//摘取到另一个优质的数据文件中
			goodFile.WriteString(lineStr+"\n")
			fmt.Println("Good:",lineStr)
		}else{
			//不合法的写入不良数据文件中
			badFile.WriteString(lineStr+"\n")
			fmt.Println("Bad:",lineStr)
		}
	}

}


/*
将文本转换字符编码并返回
srcStr		待转码的原始字符串
srcEncoding	原始字符串的编码（字符集）
dstEncoding	目标编码（字符集）
*/
func ConvertEncoding(srcStr string, srcEncoding string, dstEncoding string) (dstStr string, err error) {

	//创建指定字符集的解码器
	srcDecoder := mahonia.NewDecoder(srcEncoding)
	dstDecoder := mahonia.NewDecoder(dstEncoding)

	//将内容转换为UTF-8字符串
	utfStr := srcDecoder.ConvertString(srcStr)

	//将UTF-8字节转换为目标字符集的字节
	_, dstBytes, err := dstDecoder.Translate([]byte(utfStr), true)
	if err != nil {
		return
	}

	//还原为字符串并返回
	dstStr = string(dstBytes)
	return

}
