# 博客项目环境搭建和Model设计

## 准备工作：

a.下载goland，安装go环境，在mysql创建数据库beego_blog，然后导入beego_blog.sql
b.然后安装beego和开发工具

```
go get -u github.com/astaxie/beego
go get -u github.com/beego/bee
```

c.或许还需要安装go-mysql的驱动

```
go get github.com/go-sql-driver/mysql
```

d.然后在你的gopath的src目录下执行：

```
bee new beego_blog
```

然后看见自动创建的beego_blog文件夹，进入beego_blog文件夹，将默认的static文件夹和views文件夹替换成自己的static文件夹和views文件夹,并执行：

```
bee run
```


访问http://localhost:8080/出现welcome to beego画面表示准备工作完成，然后就可以开始我们的开发了。

## Mode设计

### 表结构分析

根据sql脚本分析可得出表之间的关系如下：

```go
tb_link:友情链接表
    id:主键，自增长
    sitename:网站名称
    url:网址
    rank:排序值

tb_mood：说说表
    id:主键，自增长
    content：说说内容
    cover:说说封面
    posttime：发表时间

tb_post:文章表
    id:主键，自增长
    userid:创建者id
    author:作者
    title：标题
    color:标题的颜色
    content:文章的内容
    tags:标签名称(所属分类)
    views:浏览量
    status：状态
    posttime：发布时间
    updated:更新时间
    istop:是否置顶
    cover：封面

tb_tag：标签表
    id:主键，自增长
    name:名称
    count:文章数量

tb_tag_post:标签文章表
    id：主键，自增长
    tagid:标签id
    postid:文章id
    poststatus:状态
    posttime:发表时间

tb_user:管理员表
    id:主键，自增长
    username:用户名
    password:密码
    email:邮箱
    logincount：登陆次数
    authkey:
    active:是否激活
```

tb_post表和tb_tag表之间的关系：

![20180905233613725](..\image\20180905233613725.png)

文章表和标签表之间是多对多的关系，因为某一篇文章(go语言数组)可以属于多个标签(go,数组),然而某个标签中可以存在多篇文章。
1.根据文章id查找其对应的一个分类，首先在标签文章表中找到该文章id所对应的标签id，然后去标签表中根据标签id查找对应的记录。
2.根据标签id查找该标签下的所有文章，首先在标签文章表中找到对应文章id(11,12,16,22,23),在文章表中根据这些文章id查找对应的记录(select * from tb_post where id in (11, 12,16,22,23))

## 设置配置文件

当我们在创建项目的时候，beego已经在app.conf文件中为我们配置了一些默认信息，但是这些信息对于我们这个博客项目来说远远不够，需要我们自己去配置相关的信息，例如：连接数据库的信息，网站头部信息等等。

```
appname = beego_blog
httpport = 8080
runmode = dev


dbprefix = tb_

dbhost = 127.0.0.1
dbport = 3306
dbuser = root
dbpassword = 111111
dbname = beego_blog

title = GoBlogBlock
keywords = studyForGo
description = goodgoodstudy,daydayup

version = 1.0

```

## 创建model

根据上面的表结构，在我们的项目中的models文件下创建项目的model,这些model和数据库中的表一一对应，在我们以后和数据的交互过程中扮演着重要角色，因为它们是我们项目和数据库交互的数据的载体，主要文件见下图：

其中base.go主要是获取配置文件中和数据库连接相关的信息，例如：登录密码，用户名等等，在获得这些信息之后注册模型，注册数据库，和数据连接。

```go
package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"crypto/md5"
	"fmt"
)

func init() {
	dbhost := beego.AppConfig.String("dbhost")
	dbport := beego.AppConfig.String("dbport")
	dbuser := beego.AppConfig.String("dbuser")
	dbpassword := beego.AppConfig.String("dbpassword")
	dbname := beego.AppConfig.String("dbname")
	dburl := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	// set default database
	orm.RegisterDataBase("default", "mysql", dburl, 30)
	// register model
	orm.RegisterModel(new(Link), new(Mood),  new(Post), new(Tag), new(TagPost), new(User))
}
```

Pager.go主要封装的是和分页相关的信息，例如：当前页码，每页显示的文章的数量等信息，其中ToString方法是核心，首先通过计算得出所要显示的页码，然后通过循环将页码和对应的rootpath写入a标签中，最后转换为字符串返回。

```go
package models

import (
	"bytes"
	"fmt"
)

type Pager struct {
	Page int //第几页
	Pagesize int //每页大小
	Totalnum int//总页数
	urlpath string //每页所对应的url
}

//创建Pager对象
func NewPager(page, pagesize, totalnum int, urlpath string) *Pager {
	pager := new(Pager)
	pager.Page = page
	pager.Pagesize = pagesize
	pager.Totalnum = totalnum
	pager.urlpath = urlpath
	return pager
}

//修改page
func (this *Pager) SetPage(page int) {
	this.Page = page
}

//修改pagesize
func (this *Pager) SetPagesize(pagesize int) {
	this.Pagesize = pagesize
}

//设置总数量
func (this *Pager) SetTotalnum(totalnum int) {
	this.Totalnum = totalnum
}

//设置rootpath
func (this *Pager) SetUrlpath(urlpath string) {
	this.urlpath = urlpath
}

func (this *Pager) url(page int) string {
	return fmt.Sprintf(this.urlpath, page)
}

func (this *Pager) ToString() string{//str2html
	//文章的总数量小于等于每页显示的文章的数量
	if this.Totalnum <= this.Pagesize {
		return ""
	}
	//偏移量
	offset := 5
	//显示10个页码
	linknum := 10
	var totalpage int
	var from int//从哪一页开始显示
	var to int //显示到哪一页
	if this.Totalnum % this.Pagesize != 0 {
		totalpage = this.Totalnum / this.Pagesize + 1
	}else {
		totalpage = this.Totalnum / this.Pagesize
	}
	//总的页码小于10，直接从第一页显示到最后一页
	if totalpage < linknum {
		from = 1
		to = totalpage
	}else {
		//计算开始显示的页码
		from = this.Page - offset
		//计算结束页码
		to = from + linknum
		//判断开始页码是否小于1
		if from < 1 {
			from = 1
			//结束页根据开始页来计算，1+10-1
			to = from + linknum - 1
		}else if to > totalpage {//结束页大于总页数
			to = totalpage
			from = to - linknum + 1//20 - 10 + 1 = 11(11-20)
		}
	}
	//开辟空间
	var buf bytes.Buffer
	buf.WriteString("<div class='page'>")
	//上一页
	if this.Page > 1 {
		buf.WriteString(fmt.Sprintf("<a href='%s'>&laquo;</a>", this.url(this.Page-1)))//<<
	}

	for i := from; i <= to; i++ {
		if i == this.Page {
			//<b> 标签规定粗体文本。
			buf.WriteString(fmt.Sprintf("<b>%d</b>", i))
		}else {
			buf.WriteString(fmt.Sprintf("<a href='%s'>%d</a>", this.url(i), i))
		}
	}
	
	//设置下一页标签
	if this.Page < totalpage {
		buf.WriteString(fmt.Sprintf("<a href='%s'>&raquo;</a>", this.url(this.Page+1)))//>>
	}
	buf.WriteString("</div>")
	str := buf.String()
	return str

}

```

其他的model中内容是和数据库意义对应的，完整代码可在github上获取。

