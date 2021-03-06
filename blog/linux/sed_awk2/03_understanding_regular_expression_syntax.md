<!--
author: 刘青
date: 2017-03-12
title: 了解正则表达式语法
tags: 
category: linux/sed_awk2
status: publish
type: note
summary:
-->

grep、sed和awk都使用正则表达式，但是不能完全使用正则表达式语法中的所有原字符。
- 基本的元字符集是由ed编辑器引入的，在grep中可用，sed送相同集合。
- egrep 程序提供了一个扩展的元字符集，awk使用这个扩展了的元字符集。

### 表达式
我们见过数学表达式：可计算的字面字符。

> 正则表达式：一种模式或字符序列的描述。

正则表达式的组成元素为：
- 以一个字面值或变量表示的值
- 一个操作符

除了元字符外，字符都被解析为自身的字面值。元字符有特定含义：
- .：匹配除换行符以外的任意单个字符。awk中也可以匹配换行符；
- *：匹配任意（≥0）个在它前面的字符；
- [...]：匹配[]中的字符集中的任意一个。如果[]中的第一个字符为^，表示不匹配[]中的字符集的元素；
- ^：如果作为正则表达式的第一个字符，表示匹配行首；
- $：如果作为正则表达式的最后一个字符，表示匹配行的结尾；
- \{n, m\}：匹配前面字符出现的次数（n ≤ x ≤ m ）。m可以省略。
- \ ：转义特殊字符

扩展元字符（egrep和awk）
- +：匹配前面的正则表达式一次或多次；
- ?：匹配前面的正则表达式零次或一次；
- |：指定可以匹配其前面或后面的正则表达式；
- ()：对正则表达式分组；
- {n, m}： 和\{n, m\}一样；

eg:
```
# 无处不在的\：转义
# 匹配字符点
\.
# 匹配 \f
\\f
#匹配一个a或一个]或一个1
[a\]1]

# 通配符 .
# 匹配 88188，88277 88.88 等
88.88

# [] 表示一个集合
# 匹配What 和 what
[Ww]hat

# [a-b]表示一个范围
# 匹配所有的大写字母：
[A-Z]
# 匹配一个数字
[0-9]

# 排除符 ^
#排除元音
[^aeiou]

# 匹配零到多次 *
# 匹配1或5开头的，后面紧随至少一个0的数字
[15]00*
# 匹配任意字符
.*

# 匹配零次或者一次

# 定位开头：^


# 定位结尾： $
# 查找空行
^$
^ *$

# 字符跨度 \{n, m\}
# 匹配 1001 10001 100001
10\{2, 4\}1

# 选择性查找 | (扩展集)
# UNIX 或 LINUX 或 NETBSD
UNIX|LINUX|NETBSD

# 分组() （扩展集）
# 匹配 BigOne 或 BigOneComputer
BigOne(Computer)?
# 匹配单复数
compan(y|ies)
```
