package idiom

import (
	"fmt"
	"os"
	"time"
)

const DB_PATH = "d:/temp/idioms-v2.0.json"

var (
	//数据管道
	chanAmbiguous = make(chan string, 20)
	chanAccurate  = make(chan string, 20)
	chanQuit      = make(chan string, 0)

	//全局内存数据
	dbDataMap = make(map[string]Idiom)
)

func main0() {

	//读入命令行参数
	//idiom.exe -cmd start -poem 大王派我来巡山
	cmdInfo := [3]interface{}{"cmd", "未知命令", "你打算干什么"}
	poemInfo := [3]interface{}{"poem", "绞尽果汁想不出", "用于启动的一行诗句"}
	retValuesMap := GetCmdlineArgs(cmdInfo, poemInfo)
	cmd := retValuesMap["cmd"].(string)
	poem := retValuesMap["poem"].(string)
	fmt.Println(cmd, poem)

	//将读入的诗句打碎丢入模糊管道
	for _, v := range poem {
		keyword := fmt.Sprintf("%c", v)
		chanAmbiguous <- keyword
	}

	//三选一读入管道数据，周期性执行
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			select {
			case keyword := <-chanAmbiguous:
				go DoAmbiguousQuery(keyword, "1", chanAccurate)
			case keyword := <-chanAccurate:
				go DoAccurateQuery(keyword)

			case <-chanQuit:
				WriteIdioms2File(dbDataMap, DB_PATH)
				os.Exit(0)
			}
		}
	}()

	//定时20秒结束主程序
	timer := time.NewTimer(20 * time.Second)
	<-timer.C
	chanQuit <- "OVER"

}
