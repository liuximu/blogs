<!--
author: 刘青
date: 2017-1-26
title: PHP编码标准
type: translate
source: php-src_5.3/CODING_STANDARDS
tags: 
category: php/src
status: publish 
summary: 
-->

======================

    PHP 编码标准
======================

本文件列出了任何程序员在PHP中添加或者修改代码都一个遵守的标准。这个文件在PHP V3.0 晚期添加的，代码至今也没有完全遵守它，但是已经在哪个方向上努力了。而V5.*版本发布，许多章节都用这些规则。

编码实现
--------------------

0. 在源文件和手册中为代码写文档。[tm]
1. 不释放赋予资源指针的函数。

  比如：``function int mail(char *to, char *from)`` 不应该释放to|from。

  例外：
  - 函数设计用来释放资源。E.g. efree()
  - 函数有一个boolean参数用来控制函数是否释放
  - 为了使用最小的内存副本而和token cache 和bison代码高耦合在一起的底层的解析器例程。
2. 和同模块其他函数高耦合的提供特定功能的函数应该使用 `static` 注释和定义。这种函数应该尽可能的避免出现。
3. 尽可能的使用定义和宏，这样常量就有有意义的名字，也更容易使用。唯一的例外是0和1分别代表false和true。其他的任何使用数字常量来指定不同的行为都应该通过 #define 来定义。
4. 当编写和字符串相关的函数时，要记住PHP对每个字符串都维护了长度思想，不要使用strlen()来计算。这样既能提高效率，也能保证二进制安全。函数在改变和获取字符串新的长度时应该返回新的函数，这样就不需要再次通过strlen()函数重复计算了。
5. 绝不使用 strncat()。如果你真的知道你在做什么，在看一眼手册，然后在考虑是否使用它。即便这样，还是要避免使用它。
6. 在PHP代码中使用 ``PHP_*`` 宏，在Zend的源代码中使用 ``ZEND_*``。虽然 ``PHP_*``宏大多数情况下是``ZEND_*``宏的别名。这样做会让被调用的宏的种类更清晰。
7. 当使用 #if 注释代码，不用只使用0，而是用 "<svn username here>_0"。比如， #if FOO_0,FOO是你的svn用户。这样就人员追踪为什么代码被注释掉，尤其是在库链接时。
8. 不要定义不可用的函数。比如，如果类库没有这个函数，就不要在该PHP版本中定义这个函数，也不要提高函数不存在的运行时错误的风险。最终用户应该使用 function_exists() 来产生函数是否存在。
9. 倾向于使用 emalloc(), efree(), estrdup() 等标准C库的类似的函数。这些函数实现了内部的“安全-网”机制确保然后为释放的内存在请求结束时进行释放。他们还在debug模式下提供有用的分配和溢出信息。

在绝大多数情况下内存的分配必须使用emalloc().
当第三方类型需要控制或者释放内存，或者当出问题的内存需要在多个请求中释放的情形下malloc()的使用需要限制。


命名惯例
-----------
1. 用户等级的函数名应该在PHP_FUNCTION()宏中。他们应该是小写的，用下划线分隔单词，使用最少的单词数。专有名词会降低函数名本身的可读性，不应该使用::
```
  好的：
  'mcrypt_enc_self_test'
  'mysql_list_fields'

  可接受的：
  'mcrypt_module_get_algo_supported_key_sizes'
  (可以是：'mcrypt_mod_get_algo_sup_key_sizes'?)
  'get_html_translation_table'
  (可以是：'html_get_trans_table'?)

  不好的：
  'hw_GetObjectByQueryCollObj'
  'pg_setclientencoding'
  'jf_n_s_i'
```
2. 如果函数是'父级'的一部分，父级应该被包含在用户函数的名称中，清楚的关联其两者的关系。形式应该是：``parent_*``::
```
  一个'foo'家族的函数，
  好的：
  'foo_select_bar'
  'foo_insert_baz'
  'foo_delete_baz'

  不好的：
  'fooselect_bar'
  'fooinsertbaz'
  'delete_foo_baz'
```
3. 用户函数的名称应该用 ``_php_``前缀，后接下划线分隔的小写单词列表进行描述。如果可以要使用'static'进行定义。
4. 变量名必须有含义。单单词的变量名必须避免，除了那些真的没有意义的变量名或者不重要的含义(e.g. for(i=0; i<100; i++) ...)。
5. 变量名应该小写，使用下划线解析分隔。
6. 函数名使用驼峰命名法是要使用最少的单词，第一个单词应该是小写的，每个单词的首字母大写::
```
  好的：
  'content()'
  'getData()'
  'buildSomeWidget()'

  不好的：
  'get_Data()'
  'buildsomewidget'
  'getI()'
```
7. 类的名称要有意义，尽量避免专有名词。名字中每个单词首字母都大写，不需要下划线分隔（第一个单词大写的驼峰命名法）。类名应该加上父集的前缀（e.g. 扩展名）::
```
  好的：
  'Curl'
  'FooBar'

  不好的：
  'foobar'
  'foo_bar'
```

  语法和缩进
  --------------

1. 绝不使用C++风格的注释(i.e. // comment)。总是使用C风格的注释。PHP使用C编写，致力于兼容任何ANSI-C编译器。虽然许多编译器都支持在C的代码中有C++风格的注释，但是依旧要确保其他编译器可以编译你的代码。Win32平台的代码是例外，因为Win32部分的代码使用 MS-Vsiual C++ 编译器，它接受C++风格的注释。
2. 使用K&R风格。当然，我们不能强制别人使用他不习惯的风格，但在写PHP核心代码或者其每个标准模块，请保持K&R风格。将其应用到任何事情上：缩进、注释、函数定义。参见[Indentstyle_](http://www.catb.org/~esr/jargon/html/I/indent-style.html).
3. 尽可能多的使用空白字符和大括号。在变量定义章节和语句块之间保持一个空号。在两个函数之间至少有一个空号，两个更好。
```
  总是：
  if (foo) {
    bar;    
  }
  而不是：
  if(foo)bar;
```
4. 使用tab字符缩进。一个tab代表四个空格。维护缩进的一致性付出重要，这样在定义、注释和流程控制结构上正确的排版。
5. 预处理语句(#if 等)必须在第一列开始。对预处理指令你应该将#放在行首，然后任意多个空白字符。


测试
---------------
1. 扩展应该可以使用 *.phpt 轻易的进行测试。细节阅读 README.TESTING。


文档和折叠
--------------
每个用户级函数在其前都应该有自己的一句话简短描述函数用途的描述用户级函数原型以确保在线文档和代码一致。大概是这样的：

```
  /* {{{ proto int abs(int numbe)
   Returns the absolute value of the number */
  PHP_FUNCTION(abs) 
  {
    ... 
  }
  /* }}} */
```

{{{ 符号是 Emacs 和 vim 折叠模式的默认的折叠符。折叠在大文件中非常有用，你可以只展开你想看的函数。 }}} 是折叠的结束符，应该另起一行排版。

"proto"关键字是 doc/genfuncsummary 脚本的一个命令，用来生成完整函数的概括。有了这个关键字，我们可以将折叠放到其他地方而不用担心函数概括的丢失。
更多可选的参数包括：

```
  /* {{{ proto object imap_header(int stream_id, int msg_no [, int from_lengt [, int subject_length [, string default_host]]])
   Returns a header object with the difined parameters */
```

即便很长，也要把原型写作一行。

新的和实验性函数
------------------
为了减少和第一次发布的新的函数的实现相关的问题，我们建议在函数目录中包含一个用 'EXPERIMENTAL'标记的文件。函数在初次实现时遵守前缀惯例。
使用'EXPERIMENTAL'标记的文件应该包括的信息有：
  然后创作信息（已知bugs，模块的未来方向）。
  不方便写在SVN注释中的状态。

别名 & 历史文档
---------------------
你也许有一些减少重复名称而废弃了的别名，比如： somedb_select_result 和 somedb_selectresult。出于归档的木丁，只有当前的名称会被记录，而别名会在父函数中列出来。为了间断两种的关系，用户级有和别名完全不一样的名称的函数要单独归档。原型应该仍然包括、描述函数的别名。

向后兼容函数和名称应该作为代码库的一部分一直维护。更多信息查看 /phpdoc/README。
