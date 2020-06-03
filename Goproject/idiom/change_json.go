package idiom

import (
	"encoding/json"
	"fmt"
	"os"
)

//将模糊查询的json转化为go数据
func ParseJson2Idioms(jsonStr string) (idiomsMap map[string]Idiom) {
	idiomsMap = make(map[string]Idiom)

	//将json转换为go数据
	tempMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &tempMap)
	if err != nil {
		fmt.Println("解析json数据错误，%v", err)
	}
	//fmt.Println(tempMap)
	dataSlice := tempMap["showapi_res_body"].(map[string]interface{})["data"].([]interface{})
	//fmt.Printf("type=%T,value=%v",dataSlice,dataSlice)
	for _, v := range dataSlice {
		title := v.(map[string]interface{})["title"].(string)
		idiom := Idiom{Title: title}
		idiomsMap[title] = idiom
	}
	return

}

//将精确查询的json转化为go数据
func ParseJson2Idiom(jsonStr string) Idiom {
	idiom := Idiom{}
	tempMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &tempMap)
	if err != nil {
		fmt.Println("解析json数据错误，%v", err)
	}
	dataMap := tempMap["showapi_res_body"].(map[string]interface{})["data"].(map[string]interface{})
	for k, v := range dataMap {
		valueStr := v.(string)
		switch k {
		case "title":
			idiom.Title = valueStr
		case "spell":
			idiom.Spell = valueStr
		case "samples":
			idiom.Sample = valueStr
		case "derivation":
			idiom.Derivation = valueStr
		case "content":
			idiom.Content = valueStr
		}
	}
	return idiom
}

//将go数据写出到json文件
func WriteIdioms2File(idiomsMap map[string]Idiom, path string) {
	fmt.Println("WriteIdioms2File")
	dstFile, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer dstFile.Close()

	encoder := json.NewEncoder(dstFile)
	err := encoder.Encode(idiomsMap)
	if err != nil {
		fmt.Println("写出json文件失败，err=", err)
		return
	}
	fmt.Println("写出json文件成功！")

}

//读入json文件为go数据
func ReadIdiomsFromFile(dstPath string) (idiomsMap map[string]Idiom, err error) {
	idiomsMap = make(map[string]Idiom)

	//读入json文件数据
	dstFile, _ := os.OpenFile(dstPath, os.O_RDONLY|os.O_CREATE, 0666)
	defer dstFile.Close()
	decoder := json.NewDecoder(dstFile)
	err = decoder.Decode(&idiomsMap)
	if err != nil {
		fmt.Println("加载数据失败！err=", err)
	} else {
		fmt.Println("成功加载数据！")
		fmt.Println("idiomsMap=", idiomsMap)
	}
	return

}
