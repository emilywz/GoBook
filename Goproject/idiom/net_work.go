package idiom

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//获取模糊查询的url
func GetAmbiguousUrl(keyword string, page string) (url string) {
	return "http://route.showapi.com/1196-1?showapi_appid=19988&showapi_sign=968ad4fcc2144e41b5c366838d1b0ec4&keyword=" + keyword + "&page=" + page + "&rows=20"
}

//获取精确查询的url
func GetAccurateUrl(keyword string) (url string) {
	return "http://route.showapi.com/1196-2?showapi_appid=19988&showapi_sign=968ad4fcc2144e41b5c366838d1b0ec4&keyword=" + keyword
}

//从url拿到json数据
func GetJson(url string) (jsonStr string, err error) {

	//获得网络数据
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http请求失败,err=", err)
		return
	}
	//延时关闭网络IO资源
	defer resp.Body.Close()

	//resp.Body实现了Reader接口，对其进行数据读入
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取网络数据失败,err=", err)
		return
	}

	//将网络数据转化为字符串输出
	jsonStr = string(bytes)
	//fmt.Println(jsonStr)

	return

}

//模糊查询
func DoAmbiguousQuery(keyword string, page string, chanAccurate chan<- string) {
	//先拿到json
	url := GetAmbiguousUrl(keyword, page)
	jsonStr, _ := GetJson(url)

	//将json转化为成语集合
	idiomsMap := ParseJson2Idioms(jsonStr)

	//将成语集合写入内存数据
	for title, idiom := range idiomsMap {
		dbDataMap[title] = idiom
	}

	//将成语的名字写入精确管道
	for title, _ := range idiomsMap {
		chanAccurate <- title
	}

	/*	chanAccurate<- "大鹏展翅"
		chanAccurate<- "隔壁老王"
		chanAccurate<- "龟派气功"
		chanAccurate<- "我很牛逼"
		chanAccurate<- "来咬我呀"
		fmt.Println("DoAmbiguousQuery",keyword,page)
	*/
}

//精确查询
func DoAccurateQuery(keyword string) {
	//fmt.Println("DoAccurateQuery",keyword)

	//拿到json
	url := GetAccurateUrl(keyword)
	jsonStr, _ := GetJson(url)

	//将json转化为一个Idiom对象
	idiom := ParseJson2Idiom(jsonStr)

	//将Idiom对象存入总集合，覆盖原来的粗糙对象
	dbDataMap[idiom.Title] = idiom

}
