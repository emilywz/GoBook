# 直接跳转goto

## 概述

goto SOMEWHERE就是：直接去到标记为SOMEWHERE地方，Go 语言的 goto 语句可以无条件地转移到过程中指定的行。

goto语句通常与条件语句配合使用。可用来实现条件转移， 构成循环，跳出循环体等功能。但是，在结构化程序设计中一般不主张使用goto语句， 以免造成程序流程的混乱，使理解和调试程序都产生困难。

## 直接跳转GAMEOVER

```go
func main() {

    //顺序结构
    fmt.Println("hello")
    fmt.Println("golang")

    //选择结构
    if time.Now().Hour()%2 == 1 {
        fmt.Println("情绪稳定")
    } else {
        fmt.Println("大姨夫蠢蠢欲动...")
    }

    //循环结构
    var i int
    for {
        if i > 10 {
            //去到GAMEOVER标记的地方
            goto GAMEOVER
        }
        fmt.Println(i)
        time.Sleep(500 * time.Millisecond)
        i++
    }

    //这里执行不到
    fmt.Println("此处免费领取靠海别野一套")
    fmt.Println("此处免费领取满汉全席一套")
    fmt.Println("此处免费领取我厂生产的女朋友一个")

GAMEOVER:
    fmt.Println("GAME OVER!")

}
```

