# walk常用控件

## 认识几种常用控件

- 按钮
- 菜单
- 工具栏
- 界面布局
- 列表
- 显示图片

## 导入依赖

```go
import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"os"
	"strings"
	"io/ioutil"
	"fmt"
	"log"
)
```

## 定义全局数据

```go
// 全局应用的菜单项
var myAction *walk.Action
//自定义的主窗口
var myWindow *MyMainWindow

//自定义窗口
type MyMainWindow struct {
	*walk.MainWindow
	te *walk.TextEdit

	//listbox使用的数据
	model *EnvModel
	//listbox控件
	listBox *walk.ListBox

}
```

## 定义列表数据模型

```go
//环境变量条目数据模型
type EnvItem struct {
	//环境变量的名字和值
	name  string
	value string
}

//列表数据模型
type EnvModel struct {
	//继承ListModelBase
	walk.ListModelBase

	//环境变量数集合
	items []EnvItem

}

//列表数据模型的工厂方法
func NewEnvModel() *EnvModel {
	env := os.Environ()
	m := &EnvModel{items: make([]EnvItem, len(env))}
	for i, e := range env {
		j := strings.Index(e, "=")
		if j == 0 {
			continue
		}

		name := e[0:j]
		value := strings.Replace(e[j+1:], ";", "\r\n", -1)
	
		m.items[i] = EnvItem{name, value}
	}
	
	return m

}
```

## 定义列表相关的监听和系统回调接口

```go
//定义列表项目的单击监听
func (mw *MyMainWindow) lb_CurrentIndexChanged() {
	i := mw.listBox.CurrentIndex()
	item := &mw.model.items[i]

	mw.te.SetText(item.value)
	
	fmt.Println("CurrentIndex: ", i)
	fmt.Println("CurrentEnvVarName: ", item.name)

}

//定义列表项目的双击监听
func (mw *MyMainWindow) lb_ItemActivated() {
	value := mw.model.items[mw.listBox.CurrentIndex()].value

	walk.MsgBox(mw, "Value", value, walk.MsgBoxIconInformation)

}

//列表的系统回调方法：获得listbox的数据长度
func (m *EnvModel) ItemCount() int {
	return len(m.items)
}

//列表的系统回调方法：根据序号获得数据
func (m *EnvModel) Value(index int) interface{} {
	return m.items[index].name
}

```

## 定义普通事件回调

```go
//显示消息窗口
func ShowMsgBox(title, msg string) int {
	return walk.MsgBox(myWindow, "天降惊喜", "你老婆跟人跑了！", walk.MsgBoxOK)
}

//一个普通的事件回调函数
func TiggerFunc() {
	ShowMsgBox("天降惊喜", "你老婆跟人跑了！")
}

```

## 主调函数

```go
func main() {
	defer func() {
		if err := recover(); err != nil {
			errMsg := fmt.Sprintf("%#v", err)
			ioutil.WriteFile("fuck.log", []byte(errMsg), 0644)
		}
	}()

myWindow = &MyMainWindow{model: NewEnvModel()}

if _, err := (MainWindow{
	AssignTo: &myWindow.MainWindow,
	Title:    "常用控件",

	//窗口菜单
	MenuItems: []MenuItem{
		//主菜单一
		Menu{
			Text: "川菜",

			//菜单项
			Items: []MenuItem{
				//菜单项一
				Action{
					AssignTo: &myAction,
					Text:     "鱼香肉丝",
					//菜单图片
					Image: "img/open.png",

					//快捷键
					Shortcut: Shortcut{walk.ModControl, walk.KeyO},

					OnTriggered: func() {
						ShowMsgBox("已下单", "我是菜单项")
					},
				},

				//分隔线
				Separator{},

				//菜单项二
				Action{
					//文本
					Text: "水煮鱼",

					//响应函数
					OnTriggered: func() {
						ShowMsgBox("已下单", "您要的菜马上就去买")
					},
				},
			},
		},

		//主菜单二
		Menu{
			Text: "粤菜",

			//菜单项
			Items: []MenuItem{
				//菜单项一
				Action{

					Text: "鱼香肉丝",
					//菜单图片
					Image: "img/open.png",

					//快捷键
					Shortcut: Shortcut{walk.ModControl, walk.KeyO},

					OnTriggered: func() {
						ShowMsgBox("已下单", "我是菜单项")
					},
				},

				//分隔线
				Separator{},

				//菜单项二
				Action{
					//文本
					Text: "水煮鱼",

					//响应函数
					OnTriggered: func() {
						ShowMsgBox("已下单", "您要的菜马上就去买")
					},
				},
			},
		},
	},
	//工具栏
	ToolBar: ToolBar{
		//按钮风格：图片在字的前面
		ButtonStyle: ToolBarButtonImageBeforeText,

		//工具栏中的工具按钮
		Items: []MenuItem{
			//引用现成的Action
			ActionRef{&myAction},

			Separator{},
			//自带子菜单的工具按钮
			Menu{
				//工具按钮本身的图文和监听
				Text:        "工具按钮2",
				Image:       "img/document-properties.png",
				OnTriggered: TiggerFunc,

				//附带一个子菜单
				Items: []MenuItem{
					Action{
						Text:        "X",
						OnTriggered: TiggerFunc,
					},
					Action{
						Text:        "Y",
						OnTriggered: TiggerFunc,
					},
					Action{
						Text:        "Z",
						OnTriggered: TiggerFunc,
					},
				},
			},
			Separator{},

			//普通工具按钮
			Action{
				Text:  "工具按钮3",
				Image: "img/system-shutdown.png",
				OnTriggered: func() {
					ShowMsgBox("天降惊喜", "你老婆跟人跑了！")
				},
			},
		},
	},

	MinSize: Size{600, 400},
	Layout:  VBox{},

	//控件们
	Children: []Widget{
		//水平局部
		HSplitter{
			MinSize: Size{600, 300},
			Children: []Widget{
				ListBox{
					StretchFactor: 1,
					//赋值给myWindow.listBox
					AssignTo: &myWindow.listBox,
					//要显示的数据
					Model: myWindow.model,
					//单击监听
					OnCurrentIndexChanged: myWindow.lb_CurrentIndexChanged,
					//双击监听
					OnItemActivated: myWindow.lb_ItemActivated,
			},
				TextEdit{
					StretchFactor: 1,
					AssignTo:      &myWindow.te,
					ReadOnly:      true,
				},
			},
		},
		HSplitter{
			MaxSize:Size{600,50},
			Children: []Widget{
				//图像
				ImageView{
					Background: SolidColorBrush{Color: walk.RGB(255, 191, 0)},
					//图片文件位置
					Image: "img/open.png",
					//和四周的边距
					Margin: 5,
					//定义最大拉伸尺寸
					MinSize: Size{50, 50},
					//显示模式
					Mode: ImageViewModeZoom,
				},

				//按钮
				PushButton{
					StretchFactor:8,
					Text: "摸我有惊喜",
					OnClicked: func() {
						ShowMsgBox("天降惊喜", "你老婆跟人跑了！")
					},
				},
			},
		},
	},
}.Run()); err != nil {
	log.Fatal(err)
}

}
```

