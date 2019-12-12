<!--
author: 刘青
date: 2016-04-02
title: Github的Webhooks服务
tags: Webhooks
category: tool/git
status: publish
summary:  当我们在git上做是一些事情，比如推了一些代码，我们可能希望它可以自动的调起一些事件，比如持续集成之类的。Webhooks就是用来做这个事情的。
-->

当我们在git上做是一些事情，比如推了一些代码，我们可能希望它可以自动的调起一些事件，比如持续集成之类的。Webhooks就是用来做这个事情的。

###什么是Webhooks
> Webhooks是一个对回调地址的一个有效负载。我们可以订阅一个git的确定事件，当事件发生时github.com就会回调我们配置的地址，我们得到请求后可以做我们想做的事情（持续集成...）。
>
> more:
> - 一个webhook是和一个组织|确定的代码库关联的。
> - 每个webhook最多有20个事件。

###支持的事件
webhook提供了很多的事件供选择，我们只需要使用我们想要的事件。
当我们配置了webhook的事件类型后，只有指定的事件发生时才会触发回调。每一个事件的回调的数据格式都不一样。
具体的事件请查看 [官方文档](https://developer.github.com/v3/activity/events/types/)。

###[创建一个Webhooks](https://developer.github.com/webhooks/creating/)
创建一个webhooks分两步：
- github上配置相应的信息：我们可以通过Web界面（github上某个repository中的setting页可以看到按钮）或者[API接口](https://developer.github.com/v3/repos/hooks/)来添加一个Webhooks。
- 对应的服务器配置：第一步添加的payload URL对应的服务器要得到请求后进行处理。

界面添加如图：

![githook add|center|700*0](http://7nliuximu.liuximu.com/git_githook_add.jpg)

####查看 Recent Deliveries
在如图的页面可以看到最近的github的回调，点进去可以查看具体的请求和相应，可用于调试。

###安全设置
实际上，服务器收到一个请求，我们并不能知道是谁发送的。但是我们应该只对github.com的请求进行本地的其他操作。
在添加节目，有一个Secret的，我们使用其可以实现用户鉴权。
具体实现请参考[官方文档](https://developer.github.com/webhooks/securing/)

###参考文档
[github.com/webhooks](https://developer.github.com/webhooks/)
