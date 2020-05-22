# 通用数据库接口sql.DB

从`*gorm.DB`连接获取通用数据库接口[*sql.DB](http://golang.org/pkg/database/sql/#DB)

```go
// 获取通用数据库对象`*sql.DB`以使用其函数
db.DB()

// Ping
db.DB().Ping()
```

## 连接池

```go
db.DB().SetMaxIdleConns(10)
db.DB().SetMaxOpenConns(100)
```