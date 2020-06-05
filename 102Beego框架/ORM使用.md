# ORM的使用

用户结构体，和订单是一对多的关系

```go
type User struct {
	Id int //beego中默认Id为主键，且自增长
	Name string //姓名
	Age int //年龄
	Nickname string  //昵称
	Userorders []*Userorder `orm:"reverse(many)"`
}

```

订单结构体，和用户是多对多一的关系

```go
type Userorder struct {
	Id int
	Orderdata string
	User *User `orm:"rel(fk)"`
	Products []*Product `orm:"rel(m2m)"`//ManyToMany
}
```

## 创建orm对象，插入用户：

```go
func (this *InsertUserController) Get() {
	this.TplName = "insertuser.html"
}
func (this *InsertUserController) Post() {
	name := this.GetString("name")
	age, err := this.GetInt("age")
	if err != nil {
		this.Ctx.WriteString("插入失败!")
	}
	nickname := this.GetString("nickname")
	user := models.User{Name:name, Age:age, Nickname:nickname}
	orm := orm.NewOrm()
	//INSERT INTO USER (NAME, age, nickname)  VALUE('Alice', 24, 'bb')
	n, err := orm.Insert(&user)
	if err == nil && n > 0 {
		this.Ctx.WriteString("插入成功!")
	}else {
		this.Ctx.WriteString("插入失败!")
	}
}

```

insertuser.html

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>插入用户</title>
</head>
<body>
    <form action="/insertuser" method="post">
        <table>
            <thead>插入用户</thead>
            <tbody>
                <tr>
                    <td>姓名:</td>
                    <td>
                        <input id="name" name="name" type="text">
                    </td>
                </tr>
                <tr>
                    <td>年龄:</td>
                    <td>
                        <input id="age" name="age" type="text">
                    </td>
                </tr>
                <tr>
                    <td>昵称:</td>
                    <td>
                        <input id="nickname" name="nickname" type="text">
                    </td>
                </tr>
                <tr>
                    <td>
                        <input type="submit" value="插入"/>
                    </td>
                    <td>
                        <input type="reset" value="重置"/>
                    </td>
                </tr>
            </tbody>
        </table>
    </form>
</body>
</html>

## 删除用户

```go
func (this *DeleteUserController) Get1() {
	//通过id删除
	//http://localhost:8080/deleteuser?id=2
	//获取id
	id, err := this.GetInt("id")
	if err != nil {
		this.Ctx.WriteString("删除失败!")
		return
	}
	//创建user对象
	user := models.User{Id:id}
	orm := orm.NewOrm()
	//delete函数默认根据id进行删除，也可以根据指定的字段进行删除
	//DELETE FROM USER WHERE id = 2;
	n, err := orm.Delete(&user)
	if n > 0 && err == nil {
		this.Ctx.WriteString("删除成功!")
		return
	}else {
		this.Ctx.WriteString("删除失败!")
	}
}

func (this *DeleteUserController) Get() {

	//通过姓名删除
	//http://localhost:8080/deleteuser?name=admin
	name := this.GetString("name")
	user := models.User{Name:name}
	
	//创建orm对象
	orm := orm.NewOrm()
	//delete函数默认根据id进行删除，也可以根据指定的字段进行删除
	//DELETE FROM USER WHERE NAME='Alice';
	n, err := orm.Delete(&user, "name")
	if n > 0 && err == nil {
		this.Ctx.WriteString("删除成功!")
		return
	}else {
		this.Ctx.WriteString("删除失败!")
	}

}

```

## 查询用户

```go
func (this *QueryUserController) Get1() {
	//获取id查询用户
	//http://localhost:8080/queryuser?id=2
	id, err := this.GetInt("id")
	if err != nil {
		this.Ctx.WriteString("查询出错!")
		return
	}
	user := models.User{Id:id}

	//创建orm
	orm := orm.NewOrm()
	//默认根据id查询
	//SELECT * FROM USER WHERE id = 4;
	err = orm.Read(&user)
	if err != nil {
		this.Ctx.WriteString("查询出错!")
		return
	}else {
		this.Ctx.WriteString("id = " + strconv.Itoa(user.Id) + "\nname = " +
			user.Name + "\nage = " + strconv.Itoa(user.Age) + "\nnickname = " + user.Nickname)
	}

}


func (this *QueryUserController) Get() {
	//http://localhost:8080/queryuser?name=admin
	//获取用户输入的姓名
	name := this.GetString("name")
	//构造user
	user := models.User{Name:name}
	//创建orm
	orm := orm.NewOrm()
	//根据name查询用户，Read函数如果不指定参数
	//SELECT * FROM USER WHERE NAME = '吴俏祥';
	err := orm.Read(&user, "name")
	if err != nil {
		this.Ctx.WriteString("查询出错!")
		return
	}else {
		this.Ctx.WriteString("id = " + strconv.Itoa(user.Id) + "\nname = " +
			user.Name + "\nage = " + strconv.Itoa(user.Age) + "\nnickname = " + user.Nickname)
	}
}

```

## 更新用户

```go
//http://localhost:8080/updateuser?id=1&name=tom&nickname=234&age=23
func (this *UpdateUserController) Get1() {
	id, err := this.GetInt("id")//获取id
	if err != nil {
		this.Ctx.WriteString("更新失败!")
		return
	}
	name := this.GetString("name")//获取姓名
	nickname := this.GetString("nickname")//获取昵称
	age, err := this.GetInt("age")
	if err != nil {
		this.Ctx.WriteString("更新失败!")
		return
	}
	user := models.User{Id:id, Name:name, Nickname:nickname, Age:age}
	//UPDATE USER SET NAME = '杜紫维' WHERE id = 1;
	//user := models.User{Id:id, Name:name, Nickname:nickname}
	orm := orm.NewOrm()
	//update函数根据对象的id进行更新,必须指定id
	n, err := orm.Update(&user)
	if n > 0 && err == nil {
		this.Ctx.WriteString("更新成功!")
	}else {
		this.Ctx.WriteString("更新失败!")
	}
}



//http://localhost:8080/updateuser?id=1&name=tom&nickname=234&age=23
func (this *UpdateUserController) Get() {
	id, err := this.GetInt("id")//获取id
	if err != nil {
		this.Ctx.WriteString("更新失败!")
		return
	}
	name := this.GetString("name")//获取姓名
	nickname := this.GetString("nickname")//获取昵称
	age, err := this.GetInt("age")
	if err != nil {
		this.Ctx.WriteString("更新失败!")
		return
	}
	user := models.User{Id:id, Name:name, Nickname:nickname, Age:age}
	//UPDATE USER SET NAME = '杜紫维' WHERE id = 1;
	//user := models.User{Id:id, Name:name, Nickname:nickname}
	orm := orm.NewOrm()
	//update函数根据指定字段进行更新,且必须指定id
	n, err := orm.Update(&user, "name")
	if n > 0 && err == nil {
		this.Ctx.WriteString("更新成功!")
	}else {
		this.Ctx.WriteString("更新失败!")
	}
}

```

## 插入订单

```go
func (this *InserOrderController) Get(){
	//创建orm
	orm := orm.NewOrm()
	//创建订单
	order := models.Userorder{}
	//对订单初始化
	order.Orderdata = "this is order123"
	order.User = &models.User{Id:6}
	//如果不指定Id，可以插入成功，但是外键user_id为0
	//order.User = &models.User{Name:"吴俏祥"}
	//插入订单
	n, err := orm.Insert(&order)
	//判断是否插入成功
	if n > 0 && err != nil {
		this.Ctx.WriteString("插入失败!")
		return
	}
	this.Ctx.WriteString("插入成功!")
}

```

## 查询订单

```go
func (this *QueryOrderController) Get(){
	orm := orm.NewOrm()
	//可以传递表名(大小写无关)和对象名
	qs := orm.QueryTable("Userorder")
	var orders []*models.Userorder
	//SELECT * FROM userorder WHERE user_id = 1;
	order_num, err := qs.Filter("user__id", 1).All(&orders)//user__id=user.id
	if err != nil {
		this.Ctx.WriteString("query order fail！")
		return
	}
	fmt.Println("order_num = ", order_num)
	for _, order := range orders {
		fmt.Println("order = ", order)
	}
	this.Ctx.WriteString("查询成功!")
}

```