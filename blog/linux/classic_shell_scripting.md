:!--
author: 刘青
date: 2016-06-01
title: Shell脚本学习指南笔记
tags: linux
category: linux
status: publish
summary:
-->
##准备工作

###软件工具的原则
- 一次只做好一件事：一个小的脚本就做好一件事情；
- 处理文本行，不要处理二进制数据；
- 使用正则表达式：文本处理强有力助手；
- 默认使用标准输入|输出：用于数据过滤器的使用；
- 避免不需要的输出：同上。另外一个原因，UNIX工具程序的使用者应该清楚会发送什么；
- 输出格式必须和可接受的输入格式一致：将结果交给其他程序处理；
- 让工具去做困难的事情：UNIX可能没有完全符合你需求的程序，但是现有工具可能完成了绝大多数工作，不要从头再来；

###为什么使用Shell脚本
- 简单：Shell是高级语言，虽然可能低效；
- 可移植

###第一行的 #!
Shell执行程序是会让UNIX内核启动一个新的进程，程序在其中执行。内核知道如何为编译型的程序做这件事，但是对于非编译型，需要告诉内核。
系统会默认使用/bin/sh来执行程序。但是现在的UNIX系统都有好几个Shell，这就需要指定了。
如果一个文件的开头是
#!，内核会扫描改行其余的部分，看是否存在可用来执行程序的解释器的完整路径和一个可选的参数。

###Shell的基本元素
- 命令和参数：
    - 命令
    - 参数
    - 可选参数：-字母
    - ; 可以分隔各个命令行：命令将依次执行
    - & 分隔各个命令行：命令将在后台执行前面的命令（Shell不用等到该命令完成就可以继续执行下一个命令）；
- 变量：变量名=变量值
- 简单的echo输出： echo [string ...]
- 华丽的printf输出：echo处理复杂输出存在兼容性问题 printf format-string [arguments ...] C++语法
- 基本的I|O重定向：
    - <：改变标准输入， tr -d '\r' < a.txt
    - >：改变标准输出， tr -d '\r' < a.txt > b.txt
    - >>：追加， tr -d '\r' < a.txt >> b.txt
    - |：管道，将前者的标准输出修改为后者的标准输入
- 特殊文件：/dev/null 与 /dev/tty
    - 所有传递到/dev/null的数据都会被系统丢弃掉。也就是说，当程序将数据写到此文件时，会任务它已经完成写入数据的操作，但实际上什么都没做。当需要的是命令的退出状态而非它的输出，此功能会很有用：
    ``