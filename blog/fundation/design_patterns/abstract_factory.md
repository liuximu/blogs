<!--
author: 刘青
date: 2016-04-02
title: 抽象工厂模式
tags: 抽象工厂 对象创建型模式
category: fundation/design_patterns
status: publish
summary: 提供一个创建一系列相关|相互依赖的对象的接口，而无需指定它们具体的类。
-->
###名称
abstract factory | 抽象工厂 | Kit

###分类
对象 && 创建型

###意图
提供一个创建一系列相关|相互依赖的对象的接口，而无需指定它们具体的类。

###场景
我们有一个支持多种视感标准的用户界面工具包，不同的视感标准的窗口组件的外观|行为不一样。
- bad实践：我们对某个特定的视感硬编码它的窗口组件，在整个应用中实例化特定视感风格的窗口组件类将使得以后很难改变视感风格，可移植性低。
- good实践：我们为每一个基本的窗口组件都有一个基类。再定义一个抽象的组件工厂类 WidgetFactory，其声明了创建每一个基本窗口组件的接口，每个子类负责创建特定的视感风格的基本窗口组件，并返回实例。客户得到这些实例，但是不知道具体的类，实现了与视感风格的解耦。



![bad|center](http://7nliuximu.liuximu.com/design_patterns_abstract_factory_bad.png)

不好的实践的如上图。
client依赖具体的视感标准，某个视感标准由各个窗口组件组成，依赖于具体的窗口组件类。
当客户程序需要具体的视感标准A时，直接使用A。而A的窗口组件是被硬编码的，当要进行修改时很麻烦。

![good|center](http://7nliuximu.liuximu.com/design_patterns_abstract_factory_good.png)
使用抽象工厂的实践如上图。
client依赖的是视感标准工厂类和各个窗口组件的基类，不和具体的实现耦合。具体的窗口组件由视感标准工厂类提供。
当客户程序需要使用具体的视感标准A是，会有一个A标准的窗口创建工厂对象B，B负责创建A标准的窗口组件。当某个窗口组件需要修改时，只需要找到B进行修改。

###参与对象：
- AbstractFactory：
	- 图中的WidgetFactory
	- 声明一个创建抽象产品对象的操作接口
- ConcreteFactory：
	- 图中的MotifWidgetFactory | PmWidgetFactory
	- 创建具体的产品对象
- AbstractProduct：
	- 图中的Window | ScrollBar
	- 为一类产品对象声明一个接口
- ConcreteProduct：
	- 图中的MotifWindow | MotifScrollBar | PmWindow | PmScrollBar
	- AbstractProduct的一个实现
- client：
	- 图中的client
	- 使用有Abstract 和 AbstarctProduct类声明的接口

###UML

![uml|center](http://7nliuximu.liuximu.com/design_patterns_abstract_factory_uml.png)

###相关模式
单例模式
