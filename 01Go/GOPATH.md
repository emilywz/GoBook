# GOPATH

```shell
File->Settings->Go->GOPATH
```

## 作用

- 存放SDK以外的第三方类库

- 可以是下载的第三方类库

- 也可以是自己收藏的可复用代码


## 目录结构

- GOPATH目录可以指定多个
- 每一个GOPATH目录下必须有一个src目录
- src目录下的文件夹名称就是引用时的包名

eg:

```go
import fuck.shit
shit.EatSome(5)
```
上述实例成功运行的前提 
GOPATH目录之一下存在路径：src/fuck/shit/ 
shit目录下的某个go源文件中有函数定义：func EatSome(kg int) 

### 全局GOPATH

- 所有工程可用
- 全局GOPATH目录也可以以GOPATH环境变量的方式配置

### 工程GOPATH

- 只有当前工程可用
  