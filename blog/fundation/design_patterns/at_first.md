<!--
author: 刘青
date: 2016-03-30
title: 设计模式讲在前面
tags: 设计模式
category: fundation/design_patterns
status: publish
summary:  
-->
我们说常用的设计模式有23个。所以这个目录至少有24篇文章。第1篇讲一些概念、推荐一些资料。

我们先看一张思维导图。[查看原图](http://7nliuximu.liuximu.com/design_patterns_DesignPatterns.png)
![设计模式|center](http://7nliuximu.liuximu.com/design_patterns_DesignPatterns.png-700)

它阐述了几个问题：
> $什么是设计模式$：每个模式描述了一个在我们周围不断重复发生的问题，以及该问题的解决方案的核心。这样，你就能一次又一次的使用该方案二不必做重复劳动。
> 其实就是面向对象软件设计的经验。这个是非常的笼统，但是十分的形象。

> $设计模式的基本要素$：
> - 名字：标识；
> - 问题：使用该模式的处境；
> - 方案：描述模式的各个部分之间的职责和协助方式；
> - 效果：模式应用的效果及使用模式应权衡的问题

> $模式的分类$：有很多种分类类型，
> 按照目的分有：
> - 创建型：和对象的创建有关；
> - 结构型：处理类|对象间的关系；
> - 行为型：描述对类|对象交互方式和和职责分配
> 
> 按照范围，有：
> - 类：通过继承建立关系，静态的；
> - 对象：处理对象的关系，动态的

|-|-|-|目的|-|
|-|-|-|-|-|
|-|-|创建型|结构型|行为型|
|范围|类|Factory Method|Adapter(类)|Interpreter<br>Template Method|
|-|对象|Abstract Factory<br>Builder<br>Prototype<br>Singleton|Adapter(对象)<br>Bridge<br>Composite<br>Decorator<br>Facade<br>Flyweight<br>Proxy|Chain of Responsibility<br>Command<br>Iterator<br>Mediator<br>Memento<br>Observers<br>State<br>Strategy<br>Visitor|

设计模式的描述使用的是UML，常用的UML有：
![uml|center](http://7nliuximu.liuximu.com/design_patterns_UML.jpg)


> 技术没有难与不难，只有会与不会。而且学习也是一个迭代的过程，每次学习都有长进。不要畏惧，我们一个一个看。


推荐一些资料吧：

 [设计模式：可复用面向对象软件的基础](http://pan.baidu.com/s/1hsNKkxY)

 [大话设计模式](http://pan.baidu.com/s/1gfDCSd9)

 [设计模式教学视频](http://pan.baidu.com/s/1i5Mwo77)
