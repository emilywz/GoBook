# Go gorm用法

## 库安装

```go
go get -u github.com/jinzhu/gorm
```

## 数据库连接 

```go
import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
）

var db *gorm.DB

func init() {
    var err error
    db, err = gorm.Open("mysql", "<user>:<password>/<database>?charset=utf8&parseTime=True&loc=Local")
    if err != nil {
        panic(err)
    }
}
```

连接比较简单，直接调用 `gorm.Open` 传入数据库地址即可

`github.com/jinzhu/gorm/dialects/mysql` 是 golang 的 mysql 驱动，实际上就是 `github.com/go-sql-driver/mysql` 

mysql，实际上支持基本上所有主流的关系数据库，连接方式上略有不同