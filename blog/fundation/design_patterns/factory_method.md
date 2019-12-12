<!--
author: 刘青
date: 2016-04-06
title: 工厂方法模式
tags: 生成器 对象创建型模式
category: fundation/design_patterns
status: publish
summary: 定义一个用于创建对象的接口，让子类决定实例化哪个类。Facttory Metnod使一个雷的实例化延迟到其子类。
-->
###名称
factory method | virtual constructor

###分类
对象 && 创建型

###意图
定义一个用于创建对象的接口，让子类决定实例化哪个类。
Facttory Metnod使一个雷的实例化延迟到其子类。

###场景
一个可以向用户显示多个文档的应用。有Application 和 Document 两个抽象类，Application负责管理Document。实例化Document的子类由Application的子类负责，但是Application类只知道何时（用户点击打开）时实例化一个Document，但是不知道实例化哪个Document的子类。为了解决这个尴尬的问题，将Document的实例化放到Application具体的子类中。


![good|center](http://7nliuximu.liuximu.com/design_patterns_factory_method_good.png)

具体是Document的哪个子类被实例化由Application的子类决定。我们称CreateDocument是一个工厂方法，因为它负责生成对象。

###参与对象：
- Product：
	- 图中的Docuemnt
	- 定义工厂方法所创建的对象的接口。
- ConcreteDocument：
	- 图中的 MyDocument
	- 实现Product的接口以。
- Creator：
	- 图中的Application
	- 声明工厂方法，其返回一个Product类型的对象。
- ConcerteCreator：
	- 图中的 MyApplition
	- 重定义工厂方法以返回ConcreteDocument实例。

###UML

![uml|center](http://7nliuximu.liuximu.com/design_patterns_factory_method_uml.png)
