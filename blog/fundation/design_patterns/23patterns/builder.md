<!--
author: 刘青
date: 2016-04-06
title: 生成器模式
tags: 生成器 对象创建型模式
category: fundation/design_patterns
status: publish
summary: 将一个复杂对象的构建和它的表示分离，使得同样的构建过程可以创建不同的表示。
-->
###名称
builder | 生成器

###分类
对象 && 创建型

###意图
将一个复杂对象的构建和它的表示分离，使得同样的构建过程可以创建不同的表示。

###场景
比如，我们有一个RTF(Rich Text Format)文档交换格式的阅读器，能将ETF转化为多种正文格式。
- bad实践：直接从RTF转换到目标格式。因为目标格式的数量是无限的，当实现新的转换时，势必会影响到阅读器本身。
- good实践：用一个TextConverter对象去配置阅读器，其可以将RTF转换为另一种正文表示。无论何时阅读器识别了一个RTF标记，它都发送给一个请求给TextConverter，TextConverter转换数据后用特定的格式表示该标记。

使用生成器模式的如图：

![good|center](http://7nliuximu.liuximu.com/design_patterns_builder_good.png)
TextConverter的子类对不同转换和不同格式进行特殊处理，所有的实现都被隐藏在抽象接口后面。
每个子类都是一个生成器（builder），而阅读器为导向器（director）。
Builder模式将分析文本算法和创建|表示一个转换后格式的算法分离开，使得文本分析算法可以复用。


###参与对象：
- Builder：
	- 图中的TextConverter
	- 为创建一个Product 对象的各个部件指定抽象接口
- ConcreteBuilder：
	- 图中的 ASCIIConvert | TeXConverter | TextWidgetConverter
	- 实现Builder的接口以构建和装配该产品的各个部件，定义并明确它所创建的表示，提供一个检索产品的接口。
- Product：
	- 图中的ASCIIText | TeXText | TextWidget
	- 被构建的复杂对象
- Director：
	- 图中的RTFReader
	- 构建一个使用Builder接口的对象

###UML

![uml|center](http://7nliuximu.liuximu.com/design_patterns_builder_uml.png)

