# math

```go
package main

import (
    "math/big"
    "fmt"
)

func main() {

    //创建大数（值可以突破int64）
    bigInt1 := big.NewInt(123)
    bigInt2 := new(big.Int)
    bigInt2.SetString("314159265358979323846264338327950288419716939937510582097494459", 10)
    fmt.Printf("type=%T,value=%v\n",bigInt1,bigInt1)
    fmt.Printf("type=%T,value=%v\n",bigInt2,bigInt2)

    //大数的计算
    bigx := big.NewInt(1)
    big1 := big.NewInt(11)
    big2 := big.NewInt(3)
    fmt.Println(bigx.Add(big1, big2),bigx)//14
    fmt.Println(bigx.Sub(big1, big2),bigx)//8
    fmt.Println(bigx.Mul(big1, big2),bigx)//33
    fmt.Println(bigx.Div(big1, big2),bigx)//3
    fmt.Println(bigx.Mod(big1, big2),bigx)//2
    fmt.Println(bigx.And(big1, big2),bigx)//3
    fmt.Println(bigx.Or(big1, big2),bigx)//11
    fmt.Println(bigx.Xor(big1, big2),bigx)//8

    //每一步的结果都重新给bigx赋值，所以事实上得到的是最后一步的结果
    fmt.Println(bigx.Add(big1,big2).Sub(big1,big2).Mul(big1,big2).Div(big1,big2))

    //(11+3)*2%3=1
    fmt.Println(bigx.Mod((bigx.Mul((bigx.Add(big1, big2)),big.NewInt(2))),big.NewInt(3)))
}

```

