<!--
author: 刘青
date: 2016-04-15
title: Web api 设计
tags: web-api-design
category: fundation
status: publish
summary: Web API的设计是让API开发者能更好更快捷的使用API进行应用开发，我们就这个目的的实现展开讨论。
-->
Web API的设计是让API开发者能更好更快捷的使用API进行应用开发，我们就这个目的的实现展开讨论。

本文章是 《Web API Design —Crafting Interfaces that Developers Love》一书的读书笔记。[下载地址](http://wenku.baidu.com/link?url=9SIe_g_yo-o8ung2ufy9iCLHENlpE0hfQ0nXo_WIaUWzUFDxo8ZwPWTtuODXpNZBHnUtjJR3CNQfe4na359egiCpvi1Eu559j-AGE17nitm)

###前言
一部分设计灵感来自于 REST。REST是一种架构风格，它没有严格的标准，极其灵活。

这本电子书是我们曾经和世界各地的合作过的API团队的设计实践的集合。
我们称我们的观点叫“实用宽松（pragmatic REST）”。因为它超越其他的任何设计原则，带领开发者走向成功。开发者是Web API的用户，衡量API设计成功与否是看开发人员多短的时间内能使用你的API。

####实用者 or RESTafarian
> RESTafarian：REST软件架构风格的热情支持者。

我们支持实用主义，不迷信REST。我们的观点是：从外而内的视角来接近API设计，也就是我们设计API从询问 “我们打算用这个API实现什么”开始。

API的工作就是让开发人员尽可能的获得成功。我们要从应用开发人员的角度来思考API的设计。
应用开发人员是这个API策略的关键。主要是API设计原则就是让应用开发人员生产力和成功最大化。这就是我们说的“实用宽松”。

我们要做对的设计，因为设计描述事物的使用方式。这个问题现在变成：让应用开发者获得最多好处的设计是什么？

**所有具体的点和最佳实践的指导原则是：用开发者的视角看问题。**

以下是具体的点：

---------------

###用名词，不用动词
> 实用宽松设计的第一原则：让简单的事情简单。

**保持URL简单明了**

基本的URL是API最重要的功能可见性的设计。简单明了的URL让API使用起来更容易。

功能可见性是设计的一个属性：不通过需求文档交代事物的使用方式。反例就是：门得告诉用户使用方式：推还是拉？

> Web API 设计的一个关键测试就是：一个资源只应该用两个基本URL表示。

比如：
- /dogs：获取狗的集合
- /dogs/1234：获取id为1234的狗的信息

这就要求我们：**将动词排除在基本URL中**

许多API使用请求方法来参与URL设计。请求方法本身就包含了动词。
> 使用HTTP动词

| Resource| POST <br>create | GET<br>read |  PUT<br>update |DELETE<br>delete |
| :-- : | :---:| :---: | :---:| :---: |
| /dogs |创建一个对象 |获取一个列表 | 更新一系列对象 |删除所有的狗|
| /dogs/1234 |Error |获取Bo | 更新Bo或者Error |删除Bo|

总结一下：
- 每个资源只使用两个基本URL
- URL中不适用动词
- 使用HTTP动词来操作集合|元素

###复数名词 & 具体名称
**使用复数名词而不是单数名词会让API更好的被理解**
这避免了有时候单数复数都可以使用的情况。

**具体名称要比抽象名称好**

比如 “/items” 比 “/dogs” 就不那么见名知意了

###低耦合——将复杂性隐藏在 「?」后面
资源和资源之间必然有关联。我们应该如何在API中处理这种关联呢？
一种解决方案是使用 /resource/identifier/resource。

比如：
/owners/345/dogs：获取主人id为345的狗的集合

但是，我们还要加上狗的颜色等等属性呢？所以这样子不行，一个主键一层，这会有太多层。
我们可以提出一个原则：当只有一个主键时，我们放到URL中，不然放到参数列表中。

比如：
GET /dogs?color=red&owner=345

###处理错误
异常处理是软件开发者不想处理但是又必须处理的部分。异常处理对于API设计者来说尤其重要，因为API对于应用开发人员来说是一个黑盒，异常信息就成了其分析错误的关键工具。

我们先看一系列的例子：
```json
//Facebook
// HTTP Status Code : 200
{
	"type" : "OauthException", 
	"message":"(#803) Some of the
aliases you requested do not exist: foo.bar"
}
//200返回码，错误信息包含在响应里面，包括一个#803的标志
```

```json
//Twilio
// HTTP Status Code : 401
{
	"status" : "401", "message":"Authenticate",
	"code": 20003, 
	"more
info": "http://www.twilio.com/docs/errors/20003"}
//添加了一个链接
```

```json
//SimpleGeo
// HTTP Status Code : 401
{
	"code" : 401, 
	"message": "Authentication Required"
}
//只是简单的描述
```
我们其实可以提出一系列最佳实践：
- 使用HTTP 状态码：将错误信息和HTTP标准的错误码进行映射。这里有超过70个错误码，这就要去用户去了解这些状态码的含义。我们最好使用它的一个精简集。
- 尽可能的返回详细的错误信息：包括错误码信息和描述信息·

###关于版本
**强制为API加上版本号，避免发布没有版本的API**

我们先看几个例子：
|-|-|
|-|-|
|Twilio |/2010-04-01/Accounts/|
|salesforce.com |/services/data/v20.0/sobjects/Account|
|Facebook |?v=1.0|

版本号的设计有一些建议：
- 指定版本号使用 ‘v’ 前缀在最前面：/v1/dogs
- 使用单数字，不要有使用 . 因为版本号是一个接口而不是一个实现
- 至少维护一个旧的版本，为开发者留出周期进行更新

那么，版本号应该放在哪呢？
- 如版本号会影响到响应的结果，将其放入URL使其更容易看见
- 要是不会影响到响应的结果，将其放入header

###分页和部分响应
部分响应允许应用开发者只得到他们想要的数据。

我们先看几个例子：
- LinkedIn：/people:(id,first-name,last-name,industry)
- Facebook：/joe.smith/friends?fields=id,name,picture
- Google：?fields=title,media:group(media:thumbnail)
他们都提供可选的字段让用户指定需要返回的字段。

**添加可选的逗号分隔列表的字段。**

返回数据库所有的记录也是一个不好的行为。

我们也看几个例子：
- LinkedIn：start 50 and  count 25
- Facebook：offset 50 and limit 25
- Twitter：page 3 and rpp(records per page) 25

他们都表示从50到75那25条数据。
**建议对于翻页数据加上总页数字段**
**翻页数据我们应该给默认值，比如：limit=10&offset=0**

###没有相关资源的响应
像计算、翻译、转换这些行为，它并不是指向某个资源，这是，我们要：**使用动词而非名词**

比如：/convert?from=EUR&to=CNY&amount=100

在API文档中要写明：“非资源”的场景比较特殊。

###支持多格式
能够支持多种格式推荐做的事情。

我们先看几个例子：
- Google：?alt=json
- Foursquare：/venue.json
- Digg*：Accept: application/json 或 ?type=json

建议使用 dog.json这种，因为大家对文件系统都比较熟悉。

应该有一个默认格式，json就不错。

###属性名
我们先看例子：
- Twitter："created_at": "Thu Nov 03 05:19;38 +0000 2011"
- Bing："DateTime": "2011-10-29T09:35:00Z"
- Foursquare："createdAt": 1320296464

孰优孰劣这个很难说。从应用开发者受益最大化给一些建议（很多接口都是JavaScript调用）：
- 默认使用JSON
- 跟随JavaScript的命名规范
	
###关于查找
一个简单的查询后面可能是一个复杂的建模。
**使用动词q（query）来表明查询而不是名词**

举例：
- 全局查询：/search?q=fluffy+fur
- 局部查询：/owners/5678/dogs?q=fluffy+fur

###加强同子域名下的API请求
我们先看几个例子：
- Facebook：graph.facebook.com api.facebook.com
- Foursquare：api.foursquare.com
- Twitter：stream.twitter.com api.twitter.com search.twitter.com

使用多个子域名我们都可以理解：配置DNS可以指向不同的集群。
但是从应用开发者的角度吃饭，我们还是建议：
**将所有的API请求都放在一个API 子域名下**

###处理异常行为
我们这里的异常行为指的是：API的客户端不能处理我们讨论的所有的事情。
####客户端拦截HTTP错误代码
这种情况在一些 Flash版本中发生：如果你响应不是200，Flash容器会拦截这个响应并展示错误给用户。应用开发者都没有机会去处理这个错误码。

Twitter有一个邮箱的解决方案：他们提供可选参数 suppress_response_codes，如果其值为true就始终返回200。注意这个参数非常的冗长。

我们提出几点建议：
- 使用 suppress_response_codes = true
- Http 码不再仅仅是错误码
- 将然后响应码放到响应信息中
- 
####客户端只支持有限的HTTP方法
GET 和 POST 被支持但是 PUT 和 DELETE 不被支持非常的普遍。我们要维持我们前面讲的原则，可以：
**让方法作为URL中的一个可选参数**

如：/dogs/1234?method=put&location=park
这个方法其实非常危险。

###鉴权
建议使用最新最好的 OAuth 2.0。这样app就不用暴露密码了。

###API 外观模式
我们如何遵守使用的操作规范，将内部服务和系统通过有效的途径暴露给应用开发者，并可以进行迭代和维护。
> 在你想要为复杂的子系统提供一个简单的接口是使用外观模式。子系统在演进的过程中通常越来越复杂。—— 设计模式。

实现API 外观模式的基本步骤：
- 设计完美的API：包括URL,请求的参数和响应的参数，header，查询参数等
- 实现该设计
- 将系统和外观集成

在内部系统和应用开发者之间加上一层，其对请求进行初步分发处理。

