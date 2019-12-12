<!--
author: 刘青
date: 2017-1-26
title: 如何阅读PHP源码
type: note
source: http://www.php-internals.com/book/?p=chapt01/01-02-code-structure
tags: 
category: php/src
status: publish 
summary: 
-->

### PHP源码的目录结构

```
. 根目录 有很多文件都值得拜读，包括那些README。           
├── autom4te.cache 
├── build 源码编译相关的文件
├── ext 官方扩展，包括绝大多数PHP的函数的定义和实现
├── include
├── libs
├── main  PHP核心文件，重要实现PHP的基础设施
├── modules
├── netware
├── pear PHP扩展与应用仓库（PHP Extension and Application Repository） 
├── sapi 服务器抽象层的代码
├── scripts
├── tests PHP的测试脚本集合，包含PHP各项功能的测试文件。
├── TSRM PHP的线程安全是构建在TSRM库之上的，PHP实现中的常见的*G宏通常是对TSRM的封装。
├── win32 包括Windows平台相关的实现。
└── Zend Zend引擎实现
```

有必要阅读一下READM* 和 CODING_STANDARDS。
