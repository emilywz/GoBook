# Go语言介绍

## Go特征

- 简单、可靠、高效
- 并发，有趣，开源
- 内存管理，数组安全，编译迅速

## Go语言优势

- 拥有全局自动垃圾回收机制
- 拥有全局错误处理机制
- 支持丰富的类型和接口
- 支持百万级高并发
- 底层语言，权限直达系统底层
- 拥有丰富的内置类型
- 函数支持多个返回值
- 支持匿名函数和闭包函数
- 支持反射

## Go核心思想

-  组合（composition） 

## Go语言从C和C++简化的功能

- 规范的语法（不需要符号表来解析）
- 垃圾回收（独有）
- 无头文件
- 明确的依赖
- 无循环依赖
- 常量只能是数字
- int和int32是两种类型
- 字母大小写设置可见性（letter case sets visibility）
- 任何类型（type）都有方法（不是类型）
- 没有子类型继承（不是子类）
- 包级别初始化以及明确的初始化顺序
- 文件被编译到一个包里
- 包package-level globals presented in any order
- 没有数值类型转换（常量起辅助作用）
- 接口隐式实现（没有“implement”声明）
- 嵌入（不会提升到超类）
- 方法按照函数声明（没有特别的位置要求）
- 方法即函数
- 接口只有方法（没有数据）
- 方法通过名字匹配（而非类型）
- 没有构造函数和析构函数
- postincrement（如++i）是状态，不是表达式
- 没有preincrement(i++)和predecrement
- 赋值不是表达式
- 明确赋值和函数调用中的计算顺序（没有“sequence point”）
- 没有指针运算
- 内存一直以零值初始化
- 局部变量取值合法
- 方法中没有“this”
- 分段的堆栈
- 没有静态和其它类型的注释
- 没有模板
- 没有异常
- 内建string、slice和map
- 数组边界检查

## Go语言设计

-  Go语言一个重要的安全设计就是禁止隐式的类型转换 

