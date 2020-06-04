# Beego框架的安装和升级

## beego 的安装

beego 的安装是典型的 Go 安装包的形式：

```shell
go get github.com/astaxie/beego
```

### 常见问题：

git 没有安装，请自行安装不同平台的 git，如何安装请自行搜索。
git https 无法获取，请配置本地的 git，关闭 https 验证：

```shell
git config --global http.sslVerify false
```

无法上网怎么安装 beego，目前没有好的办法，接下来我们会整理一个全包下载，每次发布正式版本都会提供这个全包下载，包含依赖包。

## beego 的升级

beego 升级分为 go 方式升级和源码下载升级：

Go 升级,通过该方式用户可以升级 beego 框架，强烈推荐该方式：

```shell
go get -u github.com/astaxie/beego
```

源码下载升级，用户访问 https://github.com/astaxie/beego ,下载源码，然后覆盖到 **$GOPATH/src/github.com/astaxie/beego** 目录，然后通过本地执行安装就可以升级了：

```shell
go install  github.com/astaxie/beego
```

## bee 工具简介

bee 工具是一个为了协助快速开发 beego 项目而创建的项目，通过 bee 您可以很容易的进行 beego 项目的创建、热编译、开发、测试、和部署。

### bee 工具的安装

您可以通过如下的方式安装 bee 工具：

```go
go get github.com/beego/bee
```


安装完之后，bee 可执行文件默认存放在 GOPATH/bin∗∗里面，所以您需要把∗∗ GOPATH/bin** 里面，所以您需要把 **GOPATH/bin∗∗里面，所以您需要把∗∗GOPATH/bin 添加到您的环境变量中，才可以进行下一步。

如何添加环境变量，请自行搜索
如果你本机设置了 GOBIN，那么上面的命令就会安装到 GOBIN 下，请添加 GOBIN 到你的环境变量中。

### bee 工具命令详解

我们在命令行输入 bee，可以看到如下的信息：

```shell
C:\Users\TF>bee
Bee is a Fast and Flexible tool for managing your Beego Web Application.

USAGE
    bee command [arguments]

AVAILABLE COMMANDS

    version     Prints the current Bee version
    migrate     Runs database migrations
    api         Creates a Beego API application
    bale        Transforms non-Go files to Go source files
    fix         Fixes your application by making it compatible with newer versions of Beego
    dlv         Start a debugging session using Delve
    dockerize   Generates a Dockerfile for your Beego application
    generate    Source code generator
    hprose      Creates an RPC application based on Hprose and Beego frameworks
    new         Creates a Beego application
    pack        Compresses a Beego application into a single file
    rs          Run customized scripts
    run         Run the application by starting a local development server
    server      serving static content over HTTP on port

Use bee help [command] for more information about a command.

ADDITIONAL HELP TOPICS



Use bee help [topic] for more information about that topic.

```

## new 命令

new 命令是新建一个 Web 项目，我们在命令行下执行 bee new <项目名> 就可以创建一个新的项目。但是注意该命令必须在 GOPATH/src∗∗下执行。最后会在∗∗ GOPATH/src** 下执行。最后会在 **GOPATH/src∗∗下执行。最后会在∗∗GOPATH/src 相应目录下生成如下目录结构的项目：

```go
C:\Users\TF>bee new MyBeego

______

| ___ \
| |_/ /  ___   ___
| ___ \ / _ \ / _ \
| |_/ /|  __/|  __/
\____/  \___| \___| v1.10.0
2018/09/03 21:31:30 WARN     ▶ 0001 You current workdir is not inside $GOPATH/src.
2018/09/03 21:31:30 INFO     ▶ 0002 Creating application...
        create   C:\Go\GOPATH\src\MyBeego\
        create   C:\Go\GOPATH\src\MyBeego\conf\
        create   C:\Go\GOPATH\src\MyBeego\controllers\
        create   C:\Go\GOPATH\src\MyBeego\models\
        create   C:\Go\GOPATH\src\MyBeego\routers\
        create   C:\Go\GOPATH\src\MyBeego\tests\
        create   C:\Go\GOPATH\src\MyBeego\static\
        create   C:\Go\GOPATH\src\MyBeego\static\js\
        create   C:\Go\GOPATH\src\MyBeego\static\css\
        create   C:\Go\GOPATH\src\MyBeego\static\img\
        create   C:\Go\GOPATH\src\MyBeego\views\
        create   C:\Go\GOPATH\src\MyBeego\conf\app.conf
        create   C:\Go\GOPATH\src\MyBeego\controllers\default.go
        create   C:\Go\GOPATH\src\MyBeego\views\index.tpl
        create   C:\Go\GOPATH\src\MyBeego\routers\router.go
        create   C:\Go\GOPATH\src\MyBeego\tests\default_test.go
        create   C:\Go\GOPATH\src\MyBeego\main.go
2018/09/03 21:31:30 SUCCESS  ▶ 0003 New application successfully created!

```

## run 命令

我们在开发 Go 项目的时候最大的问题是经常需要自己手动去编译再运行，bee run 命令是监控 beego 的项目，通过 fsnotify监控文件系统。但是注意该命令必须在 $GOPATH/src/appname 下执行。
这样我们在开发过程中就可以实时的看到项目修改之后的效果：

```shell
C:\Users\TF>bee run MyBeego

______

| ___ \
| |_/ /  ___   ___
| ___ \ / _ \ / _ \
| |_/ /|  __/|  __/
\____/  \___| \___| v1.10.0
2018/09/03 21:34:38 INFO     ▶ 0001 Using 'MyBeego' as 'appname'
2018/09/03 21:34:38 INFO     ▶ 0002 Initializing watcher...
MyBeego/controllers
MyBeego/routers
MyBeego
2018/09/03 21:34:48 SUCCESS  ▶ 0003 Built Successfully!
2018/09/03 21:34:48 INFO     ▶ 0004 Restarting 'MyBeego.exe'...
2018/09/03 21:34:48 SUCCESS  ▶ 0005 './MyBeego.exe' is running...
2018/09/03 21:34:48.946 [I] [asm_amd64.s:2361]  http server Running on http://:8080

```

我们打开浏览器就可以看到效果 http://localhost:8080/:

![20180903213919927](..\image\20180903213919927.png)

## version 命令

这个命令是动态获取 bee、beego 和 Go 的版本，这样一旦用户出现错误，可以通过该命令来查看当前的版本:

```shell
C:\Go\GOPATH\src>bee version

______

| ___ \
| |_/ /  ___   ___
| ___ \ / _ \ / _ \
| |_/ /|  __/|  __/
\____/  \___| \___| v1.10.0

├── Beego     : 1.10.1
├── GoVersion : go1.10
├── GOOS      : windows
├── GOARCH    : amd64
├── NumCPU    : 12
├── GOPATH    : C:\Go\GOPATH
├── GOROOT    : C:\Go\
├── Compiler  : gc
└── Date      : Monday, 3 Sep 2018

```

学院Go语言视频主页
https://edu.csdn.net/lecturer/1928