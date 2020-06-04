package utils

import (
	"fmt"
	/*MySQL*/
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/garyburd/redigo/redis"
	"errors"
	"sync"
)

var(
	//数据库读写锁
	dbMutext sync.RWMutex
)

/*
通用Mysql查询工具
tableName	要查询的表名
argsMap		查询条件集合
dest		查询结果存储地址
*/
func QueryFromMysql(tableName string,argsMap map[string]interface{},dest interface{}) (err error) {
	fmt.Println("QueryScoreFromMysql...")

	//写入期间不能进行数据库读访问
	dbMutext.RLock()
	db, err := sqlx.Connect("mysql", "root:123456@tcp(localhost:3306)/driving_exam")
	HandlerError(err, `sqlx.Connect("mysql", "root:123456@tcp(localhost:3306)/driving_exam")`)
	defer db.Close()

	selection := ""
	values := make([]interface{}, 0)
	for col,value := range argsMap{
		selection += (" and "+col+"=?")
		values = append(values, value)
	}
	selection = selection[4:]
	sql := "select * from "+tableName+" where "+selection;

	err = db.Select(dest, sql, values...)
	if err != nil {
		fmt.Println(err, `db.Select(&examScores, "select * from score where name=?;", name)`)
		return
	}

	dbMutext.RUnlock()
	return nil

}

/*将全员考试成绩单写入MySQL数据库*/
func WriteScore2Mysql(scoreMap map[string]int) {
	//锁定为写模式，写入期间不允许读访问
	dbMutext.Lock()
	db, err := sqlx.Connect("mysql", "root:123456@tcp(localhost:3306)/driving_exam")
	HandlerError(err, `sqlx.Connect("mysql", "root:123456@tcp(localhost:3306)/driving_exam")`)
	defer db.Close()

	for name, score := range scoreMap {

		_, err := db.Exec("insert into score(name,score) values(?,?);", name, score)
		HandlerError(err, `db.Exec("insert into score(name,score) values(？，？);", name, score)`)
		fmt.Println("插入成功！")
	}
	fmt.Println("成绩录入完毕！")

	//解锁数据库，开放查询
	dbMutext.Unlock()

}

/*根据姓名从Redis缓存查询分数*/
func QueryScoreFromRedis(name string) (score int, err error) {
	fmt.Println("QueryScoreFromRedis...")
	conn, err := redis.Dial("tcp", "localhost:6379")
	HandlerError(err, `redis.Dial("tcp", "local:6379")`)
	defer conn.Close()

	reply, e := conn.Do("get", name)
	if reply != nil {
		score, e = redis.Int(reply, e)
		//fmt.Println("!!!!!!!!!!!!", score, e)
	} else {
		return 0, errors.New("未能从Redis中查到数据")
	}

	if err != nil {
		fmt.Println(err, `conn.Do("get", name)或者redis.Int(reply, err)`)
		return 0, e
	}

	return score, nil

}

/*将姓名与分数写入Redis缓存*/
func WriteScore2Redis(name string, score int) error {
	conn, err := redis.Dial("tcp", "localhost:6379")
	HandlerError(err, `redis.Dial("tcp", "local:6379")`)
	defer conn.Close()

	_, err = conn.Do("set", name, score)
	fmt.Println("Redis写入成功！")
	return err

}
