# strings

## 概述

- strings库定义了对字符串的常用处理操作
- 常用的包括：检索子串、格式化、比较大小、裁剪、拼接、炸碎等

```go
func main() {

    //判断是否包含子串
    fmt.Println(strings.Contains("Hello", "shit"))        //false
    fmt.Println(strings.Contains("Hello", "lo"))          //true
    fmt.Println(strings.ContainsAny("Hello shit", "ae"))  //true 因为包含ae中的e
    fmt.Println(strings.ContainsAny("Hello shit", "azx")) //false 母串中不包含azx中的任何一个字母
    fmt.Println(strings.ContainsRune("Hello", 'o'))       //true 判断是否包含某字符
    fmt.Println(strings.ContainsRune("Hello", 'x'))       //false

    //返回子串索引
    fmt.Println(strings.Index("Hello shit", "shit"))   //6
    fmt.Println(strings.Index("Hello shit", "shiter")) //-1

    //格式化
    fmt.Println(strings.Title("hello shit")) //Hello Shit  每个单词的首字母为【标题格式】
    fmt.Println(strings.ToTitle("ShIt"))     //SHIT 单词的每个字母为【标题格式】
    fmt.Println(strings.ToUpper("ShIt"))     //SHIT 全大写
    fmt.Println(strings.ToLower("ShIt"))     //shit 全小写

    //比较大小
    fmt.Println(strings.Compare("azzz", "zaaa")) //-1
    fmt.Println(strings.Compare("z", "a"))       //1
    fmt.Println(strings.Compare("a", "a"))       //0

    //裁剪
    fmt.Println(strings.Trim("   hello shit ", " ")) //hello shit 收尾的空白都没有了
    fmt.Println(strings.Trim("   hello shit .,", " ,"))
    //去除【头尾】的所有【逗号和空白】，【头尾】是由cutset中的字符连续排列的区域
    fmt.Println(strings.Trim(", ,,  ,hello shit       .,   ,  ,  ", ", "))
    fmt.Println(strings.TrimLeft(", ,,  ,hello shit       .,   ,  ,  ", ", "))  //左侧修剪的干干净净，右侧一屁股翔
    fmt.Println(strings.TrimRight(", ,,  ,hello shit       .,   ,  ,  ", ", ")) //右侧修剪的干干净净，左侧一屁股翔
    //最后一个参数func允许你自定义裁剪规则，返回true的将被视为可以裁剪的字符
    fmt.Println(strings.TrimFunc(", ,,  ,hello shit       .,   ,  ,  ", func(r rune) bool {
        if r == ' ' || r == ',' || r == '.' {
            return true
        }
        return false
    }))
    fmt.Println(strings.TrimPrefix("-hello shit       .,   ,  ,  -", "-")) //修剪前缀
    fmt.Println(strings.TrimSuffix("-hello shit       .,   ,  ,  -", "-")) //修剪前缀
    fmt.Println(strings.TrimSpace("  hello shit     "))

    //判断是否相等
    fmt.Println(strings.EqualFold("HelloShit", "heLLoSHIT")) //true 忽略大小写和标题的差异后，是否依然相同

    //炸碎
    fmt.Println(strings.Fields("hello shit eater"))          //[hello shit eater] 以空白为分隔符，将字符串炸开为子串的切片
    fmt.Println(strings.Fields("hello,shit eater,devourer")) //[hello,shit eater,devourer]
    //在func中自定义分隔符的规则
    fmt.Println(strings.FieldsFunc("hello,shit eater,devourer", func(r rune) bool {
        if r == ' ' || r == ',' {
            return true
        }
        return false
    })) //[hello shit eater devourer]

    //炸碎
    /*
    func Split(s, sep string) []string
    func SplitN(s, sep string, n int) []string
    func SplitAfter(s, sep string) []string
    func SplitAfterN(s, sep string, n int) []string
    */
    fmt.Println(strings.Split("hello,world,hello,shitEater",","))//[hello world hello shitEater]
    fmt.Println(strings.SplitN("hello,world,hello,shitEater",",",2))//[hello world,hello,shitEater]
    fmt.Println(strings.SplitAfter("hello,world,hello,shitEater",","))//[hello, world, hello, shitEater]
    fmt.Println(strings.SplitAfterN("hello,world,hello,shitEater",",",2))//[hello, world,hello,shitEater]

    //拼接
    fmt.Println(strings.Join([]string{"hello","world","hello","shitEater"},"$"))//hello$world$hello$shitEater

    fmt.Println()

}
```

