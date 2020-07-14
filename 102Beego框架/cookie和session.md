# cookie和session

## Cookie

### 什么是Cookie

HTTP是无状态的协议，服务器不能记录浏览器的访问状态，也就是说服务器不能区分两次请求是否是一个客户端，这样的设计严重的阻碍了web程序的设计。

Cookie是解决HTTP协议无状态的方案之一。

Cookie实际上就是服务器保存在浏览器上的一段信息，浏览器有了Cookie之后，每次向服务器发送请求时都会同时将该信息发送给服务器，服务器在收到请求之后，就可以根据该信息处理请求。

Cookie由服务器创建，并发送给浏览器，最终由浏览器保存。

### Cookie的用途

电商网站的购物车

保持用户登录状态

### Cookie的缺点

Cookie作为请求或者响应的报文进行发送，无形中增加了网络流量

Cookie使用明文传送安全性差

各个浏览器对Cookie有限制，每个网站大概只能保存20个Cookie

可以在浏览器禁止Cookie

### Cookie的使用

```
this.Ctx.SetCookie("username", "admin")
```

1.设置Cookie,参数一为Cookie名称,参数二为Cookie值

2.设置Cookie是不能存在中文，如果存在中文虽然可以设置成功，但是不能取出.

```
取出Cookie：
this.Ctx.GetCookie("username")
```

```
this.Ctx.SetCookie("username", "admin", 10)
```

1.参数三为Cookie的存活时间，单位秒。

2.如果不设置时间，Cookie只在本次回话有效，Beego默认存活3600秒。

```
this.Ctx.SetCookie("username", "admin", 60*60*60*60, "/demo2")
```

1.Cooie的路径是指告诉浏览器访问那些网址时需要带上该Cookie.

2.浏览器会保存很多网站的Cookie，比如百度的Cookie,新浪的Cookie,腾讯的Cookie,浏览器在访问百度时不可能发送新浪的Cookie，也不可能访问每个网站都带上所有的Cookie,所以我们需要为Cookie设置一个Path属性，用于告诉浏览器何时需要带上该Cookie。

3.我们访问的网站必须是Cookie保存路径的子集，在访问服务器时才会带上该Cookie。

4.如果是想设置Cookie的保存路径而不设置存活时间是不行的，因为SetCookie方法默认第三个参数是设置存活时间，第四个参数是设置保存路径。

## Session

### Session是什么

Session是一段保存在服务器上的信息，当客户端第一次访问服务器时创建Session，同时也创建一个名为beegosessionID,值为创建的Seesion的id的Cookie。

这个beegoseesionID对应服务器中的一个Session对象，通过它就可以获取到保存用户信息的Session。

### Seeion的必要性

使用Cookie有一个非常大的局限，如果Cookie很多，则无形的增加了客户端与服务端的数据传输。

Cookie安全性差

由于大部分浏览器对Cookie大小和数量有限制，所以不能在Cookie中保存过多的信息。

### Session的使用

Beego默认关闭Session,如果想要使用Session,需要在主函数中或配置文件中开启Session,分别为

```
Beego.BConfig.WebConfig.Session.SessionOn = true,session = true
```

```
设置Session: 

this.SetSession("username", "admin")
```

this.GetSession("username"),如果没有获取到Session和GetCookie不一样，GetCookie会返回空字符串，但是GetSession返回nil。

通过this.DelSession("username")和this.DestroySession()均可删除Session，其区别在于DelSession删除指定Session,而DestroySession删除全session。

设置Session的时候可以存在中文。Beego中的Session在当前回话默认存活3600秒，可以在主函数或配置文件中设置Session的存活时间，分别为beego.BConfig.WebConfig.Session.SessionGCMaxLifetime，sessiongcmaxlifetime
可以通过beego.BConfig.WebConfig.Session.SessionName = "hellobeegoid123" 和sessionname=hellobeegoid来设置session的名称，如果同时在配置文件和主函数中设置，主函数优先，应为beego先读取配置文件后执行主函数，所以主函数优先。
