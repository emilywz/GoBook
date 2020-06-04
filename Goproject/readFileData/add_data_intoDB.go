package main

import (
	"bufio"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"runtime"
)

type KfPerson struct {
	Id     int    `db:"id"`
	//姓名
	Name   string `db:"name"`
	//身份证号码
	Idcard string `db:"idcard"`
	//性别
	Gender string `db:"gender"`
	//生日
	Birthday string `db:"birthday"`
	//地址
	Address string `db:"address"`
}

var (
	//全局数据管道
	chanData chan *KfPerson

	//管道是否已关闭
	readingFinished int
	chanDataClosed  bool
	wg sync.WaitGroup
)

/*全局错误处理*/
func DBHandleError(err error, why string) {
	if err != nil {
		fmt.Println("ERROR OCCURED!!!", err, why)
	}
}

/*文本大数据入库*/
func main() {

	//记录开始时间
	start := time.Now().Unix()

	//主协程数据库连接

	db, err := sqlx.Connect("mysql", "root:123456@tcp(127.0.0.1:3306)/kaifang")
	DBHandleError(err, "sqlx.Open")
	defer db.Close()

	//必要时先建表
	_, err = db.Exec(`create table if not exists kfperson(
  		id int primary key auto_increment,
  		name varchar(20),
  		idcard char(18),
  		gender char(1),
  		birthday char(8),
  		address varchar(100)
	);`)
	DBHandleError(err, "db.Exec create table")
	fmt.Println("数据表已创建")

	//创建全局数据管道
	chanData = make(chan *KfPerson)

	//开辟20条协程，执行数据插入任务，并注册在等待组
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			//在协程中执行数据插入任务
			DoInsertJob()
			wg.Done()
		}()
	}

	/*开辟10条读取协程，分别读取不同的文件*/
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			//打开大数据文本
			file, e := os.Open(`.\dataSet\kaifang_good_` + strconv.Itoa(index) + `.txt`)
			DBHandleError(e, "os.Open")
			defer file.Close()

			//创建缓冲读取器
			reader := bufio.NewReader(file)
			fmt.Println("大数据文本已打开")

			/*逐条读入开房者信息，丢入管道*/
			for {

				//读取一行数据
				lineBytes, _, err := reader.ReadLine()
				DBHandleError(err, "reader.ReadLine")

				//读到文件末尾时，关闭数据管道（通知其它协程停止对该管道的扫描）
				if err == io.EOF {
					//完毕的读取协程+1
					readingFinished ++

					//读取协程全部完毕时，关闭数据管道
					if readingFinished > 9 && !chanDataClosed {
						close(chanData)
						chanDataClosed = true
					}

					//退出读取
					break
				}

				//以逗号为定界符，将字符串数据炸碎为字段
				lineStr := string(lineBytes)
				fields := strings.Split(lineStr, ",")

				if len(fields) > 4 {
					//将合法字段封装为KfPerson，并丢入全局数据管道
					name, idcard, gender, birthday, address := fields[0], fields[1], fields[2], fields[3], fields[4]
					kfPerson := KfPerson{Name: name, Idcard: idcard, Gender: gender, Birthday: birthday, Address: address}
					chanData <- &kfPerson
				} else {
					//fmt.Println("脏数据：", lineStr)
				}
			}

			wg.Done()
		}(i)
	}

	//等待所有子协程完成任务
	wg.Wait()
	fmt.Println("main over!")

	end := time.Now().Unix()
	fmt.Printf("共用时%d秒\n", end-start)

}

/*执行数据插入任务*/
func DoInsertJob() {

	//创建当前协程的数据库连接
	db, err := sqlx.Connect("mysql", "root:123456@tcp(127.0.0.1:3306)/kaifang")
	DBHandleError(err, "sqlx.Open")
	defer db.Close()

	//创建KfPerson切片，长度达到阈值时，做一次数据库写入操作
	kfs := make([]*KfPerson, 0)

	//扫描管道数据（直到管道被关闭）
	for kfPerson := range chanData {

		//向切片中添加开房者
		kfs = append(kfs, kfPerson)

		//切片中的数据量每达到1000（或者管道已关闭），就执行一次数数据库写入操作
		if len(kfs) > 1000 || chanDataClosed {
			//执行数据库插入
			insertPersons2DB(db, kfs)

			//清空切片并重新创建
			CleanSlice(kfs)
			kfs = make([]*KfPerson, 0)
		}
	}

}

/*清空切片，回收内存，避免内存泄露*/
func CleanSlice(s []*KfPerson) {
	for i := 0; i < len(s); i++ {
		s[i] = nil
	}
	runtime.GC()
}

/*将切片中的数据一次性插入DB中*/
func insertPersons2DB(db *sqlx.DB, kps []*KfPerson) error {
	/*文本大数据中含有各种各样不合法的脏数据，做好异常的处理*/
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("!!!!!!!!!!!!!!!!!!!!", err, "!!!!!!!!!!!!!!!!!!!!")
		}
	}()

	//构建SQL语句
	sql := `insert into kfperson(name,idcard,gender,birthday,address) values`

	/*拼接每个开房者信息到SQL语句中*/
	for _, kp := range kps {
		/*
		根据做一些必要的数据处理
		//对姓名字段去头尾空格
		//滤掉姓名超长的数据
		//其它的数据处理
		*/
		kp.Name = strings.TrimSpace(kp.Name)
		if len(kp.Name) > 20 {
			continue
		}
		//kp.Address = strings.Replace(kp.Address, `"`, "*", -1)
		//kp.Address = strings.Replace(kp.Address, `''`, "*", -1)

		/*拼接名字为SQL语句*/
		personValue := `("` + kp.Name + `","` + kp.Idcard + `","` + kp.Gender + `","` + kp.Birthday + `","` + kp.Address + `"),`
		//insert into kfperson(name,idcard,gender,birthday,address) values("张三","123456199001011234","M","19900101","火星"),("李四","123456199001011234","M","19900101","火星"),("王五","123456199001011234","M","19900101","火星")
		sql += personValue
	}

	//去掉最后一个逗号再加分号，形成最终SQL语句
	sql = sql[:len(sql)-1] + ";"
	//fmt.Println(sql)

	//执行SQL语句
	result, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}

	//打印受影响的行数
	affected, err := result.RowsAffected()
	if err == nil {
		fmt.Println("执行成功，affected=", affected, err)
	} else {
		fmt.Println("!!!!!!!!!执行失败!!!!!!!!!!", err)
	}

	return nil

}
