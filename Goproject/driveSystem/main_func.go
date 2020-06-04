package main

import (
	"fmt"
	"sync"
	"time"
	"utils"
)

var (
	chNames = make(chan string, 100)
	examers = make([]string, 0)

	//信号量，只有5条车道
	chLanes = make(chan int, 5)
	//违纪者
	chFouls = make(chan string, 100)

	//考试成绩
	scoreMap = make(map[string]int)
)

/*巡考逻辑*/
func Patrol() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		//fmt.Println("战狼正在巡考...")
		select {
		case name := <-chFouls:
			fmt.Println(name, "考试违纪!!!!! ")
		default:
			fmt.Println("考场秩序良好")
		}
		<-ticker.C
	}
}

/*考试逻辑*/
func TakeExam(name string) {
	chLanes <- 123
	fmt.Println(name, "正在考试...")

	//记录参与考试的考生姓名
	examers = append(examers, name)

	//生成考试成绩
	score := utils.GetRandomInt(0, 100)
	//fmt.Println(score)
	scoreMap[name] = score
	if score < 10 {
		score = 0
		chFouls <- name
		//fmt.Println(name, "考试违纪！！！", score)
	}

	//考试持续5秒
	<-time.After(400 * time.Millisecond)

	<-chLanes
	//wg.Done()

}

/*二级缓存查询成绩*/
func QueryScore(name string) {
	score, err := utils.QueryScoreFromRedis(name)
	if err != nil {
		fmt.Println(err)
		//score, _ = utils.QueryScoreFromMysql(name)

		scores := make([]utils.ExamScore, 0)
		argsMap := make(map[string]interface{})
		argsMap["name"] = name
		//argsMap["score"] = 50
		err = utils.QueryFromMysql("score", argsMap, &scores)
		utils.HandlerError(err, `utils.QueryFromMysql("score", argsMap, &scores)`)
		fmt.Println("Mysql成绩：", name, ":", scores[0].Score)

		/*将数据写入Redis*/
		err := utils.WriteScore2Redis(name, scores[0].Score)
		if err != nil {
			fmt.Println("error message:%v", err)
		}

	} else {
		fmt.Println("Redis成绩：", name, ":", score)
	}
	//wg.Done()

}

/*
考场签到，名字丢入管道；
只有5个车道，最多供5个人同时考试；
考生按签到顺序依次考试，给予考生10%的违规几率；
每3秒钟巡视一次，发现违规的清出考场，否则输出考场时序良好；
所有考试者考完后，向MySQL数据库录入考试成绩；
成绩录入完毕通知考生，考生查阅自己的成绩；
当前目录下的成绩录入MySQL数据库,数据库允许一写多读；
再次查询成绩使用Redis缓存（二级缓存）；
整理优化代码，提高复用程度；
*/

var (
	//同步
	//创建任务使用wg.Add()，任务完成使用wg.Done()来任务减1，使用wg.Wait()来阻塞等待所有任务完成
	wg sync.WaitGroup
)

/*主程序*/
func main() {
	//获得随机名字
	for i := 0; i < 20; i++ {
		chNames <- utils.GetRandomName()
	}
	close(chNames)

	/*巡考*/
	go Patrol()

	/*考生并发考试*/
	for name := range chNames {
		wg.Add(1)
		go func(name string) {
			TakeExam(name)
			wg.Done()
		}(name)
	}

	wg.Wait()
	fmt.Println("考试完毕！")

	/*录入成绩*/
	wg.Add(1)
	go func() {
		utils.WriteScore2Mysql(scoreMap)
		wg.Done()
	}()
	//故意给一个时间间隔，确保WriteScore2DB先抢到数据库的读写锁
	<-time.After(1 * time.Second)

	/*考生查询成绩*/
	for _, name := range examers {
		wg.Add(1)
		go func(name string) {
			QueryScore(name)
			wg.Done()
		}(name)
	}
	<-time.After(1 * time.Second)
	for _, name := range examers {
		wg.Add(1)
		go func(name string) {
			QueryScore(name)
			wg.Done()
		}(name)
	}

	wg.Wait()
	fmt.Println("END")

}
