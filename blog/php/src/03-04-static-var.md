<!--
author: 刘青
date: 2017-03-28
title: 静态变量
type: note
source: https://github.com/reeze/tipi/blob/master/book/chapt03/03-04-static-var.markdown
tags: 
category: php/src
status: publish 
summary: 
-->

> 静态变量：静态分配的变量，生命周期和程序的生命周期一致（程序退出时和自动释放是才回收）。

静态变量可以分为：
- 静态全局变量：PHP全局变量都可以理解为静态全局变量。因为除非明确调用`unset`释放，它在程序运行过程中始终存在
- 静态局部变量：在函数内定义的静态变量，函数在执行时对变量的操作会保持到下一次函数被调用
- 静态成员变量：在类中定义的静态变量，在所有实例中共享。

这里讨论静态局部变量。

有PHP代码：
```
[php]
function t() {}{
    stratic $i = 0;
    $i++;
    echo $i, ' ';
}

t();
t();
t();
```

`static` 是关键，涉及到词法分析，语法分析，opcode中间码生成和执行。编译原理，跳过。
