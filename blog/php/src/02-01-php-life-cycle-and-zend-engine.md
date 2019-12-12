<!--
author: 刘青
date: 2017-03-15
title: 生命周期和Zend引擎
type: note
source:
https://github.com/reeze/tipi/blob/master/book/chapt02/02-01-php-life-cycle-and-zend-engine.markdown
tags: 
category: php/src
status: publish 
summary: 
-->

### 一切的开始：SAPI接口
> SAPI:(Server Application Programming Interface) PHP具体应用的编程接口。

PHP的执行方法有很多，只要符合SAPI规范，PHP就能正常运行。方式有：
- 通过web服务器
- 命令行下通过PHP解释器
- 嵌入在其他程序。

它们看起来好像不一样，但是工作流程实际上是一样的。脚本的执行都是以SAPI接口实现开始的，然后返回数据。

### 开始和结束
PHP开始执行后会经过两个重要的阶段：
- 处理请求之前的开始阶段，包括：
    - 模块初始化（MINIT）：在整个SAPI生命周期（eg Apache启动以后的整个生命周期内或命令行程序整个执行过程）内执行一次；
    - 模块激活（RINIT）：发生在请求阶段（eg 通过url请求某个页面）
- 处理请求之后的结束阶段，一般脚本执行到末尾或者通过exit()|die()，包括：
    - RSHUTDOWN：对应RINT
    - MSHUTDOWN：对应MINIT，在Web服务器退出或者命令行脚本执行完毕退出是关闭模块

PHP会在MINIT阶段调用所有扩展模块的MINIT函数，具体实现跨越通过如下宏来实现这些函数的回调：
```
PHP_MINIT_FUNCTION(myphpextension)
{
    //注册常量或者类等初始化操作
    return SUCCESS;
}
```

上面完成以后PHP旧初始化执行脚本的基本环境了。接着就是调用所有模块的RINT函数
```
PHP_RINT_FUNCTION(myphpextension)
{
    //eg 记录时间，在后面就能算出请求花费的时间
    return SUCCESS;
}
```

请求处理完成就进入结束阶段，PHP提供钩子函数：
```
PHP_RSHUTDOWN_FUNCTION(myphpextendsion)
{
    //eg 记录请求结束时间，写日志
    return SUCCESS;
}
```

#### 单进程SAPI生命周期
> 单进程SAPI模式：CLI/CGI模式。请求在处理异常后就关闭。

生命周期为：
- $ php -f test.php
- 调用每个扩展的MINIT
- \>\>\> 独立请求开始
- 请求test.php文件
- 调用每个扩展的RINIT
- 执行test.php
- 调用每个扩展的RSHUTDOWN
- 清除test.php
- <<< 独立请求结束
- 调用每个扩展的MSHUTDOWN
- 中断PHP

进一步补充：
- 启动阶段：在调用每个模块的初始化之前会：
    - 初始化若干全局变量：大多数情况是将其设置为NULL，少数例外；
    - 初始化若干变量：要么硬编码（PHP_VERSION），要么卸载配置头文件中(PEAR_EXTENSION_DIR);
    - 初始化Zend引擎和核心组件：zend_startup()函数。内存管理初始化、全局使用的函数指针初始化，对PHP源文件进行词法分析、语法分析、中间代码执行的函数指针赋值，初始化若干HashTable，为ini文件解析做准备，注册内置函数，注册标准常量(E_ALL, TRUE, NULL...)，注册GLOBALS全局变量等。
    - 解析php.ini：php_ini_config()函数。设置配置参数，加载zend扩展并注册PHP扩展函数。
        - CLI模式下会做如下初始化：
        ```
        INI_DEFAULT("report_zend_debug", "0");
        INI_DEFAULT("display_errors", "1");
        ```
        - 判断是否有php_ini_path_override
        - 如果没有php_ini_path_override,判断php_ini_ignore是否为空
        - 如果不忽略ini配置，开始处理php_ini_search_path（查找ini文件的路径）
        - 如果ini路径可以打开，打开ini
    - 全局操作函数的初始化：php_startup_auto_globals()函数。初始化在用户空间所使用频率很高的一些全局变量，eg：$_GET等。只是初始化。php_startup_sapi_content_types()函数用来初始化SAPI对于不同类型内容的处理函数。
    - 初始化静态构建的模块和共享模块(MINIT)：php_register_internal_extensions_func()用来注册静态构建的模块，也就是没人加载的模块，也就是内置模块。
    - 静止函数和类：将php.ini中指定的disable_functions通过调用zend_disable_functions()供CG(function_table)中删除，disable_classes通过调用zend_disable_class()从GC(class_table)中删除。
- Activation：在处理了文件相关的内容，PHP会调用php_request_startup()做请求初始化操作：
    - 调用每个模块的RINIT函数
    - 激活Zend引擎：gc_reset()用来重置垃圾收集机制；init_compiler()用来初始化编译器；init_executor()用来初始化中间代码执行过程。
    - 激活SAPI：sapi_activate()用来初始化SG(sapi_headers)和SG(request_info)
    - 环境初始化：用户空间需要用到的的一下环境变量的初始化，包括服务器环境，请求数据环境
    - 模块请求初始化：zend_activate_modules()实现模块的请求初始化，也就是调用每个扩展的RINIT
- 运行：
    - php_execute_script()包含了运行脚本的全部过程。
    - zend_compile_file()做词法分析、语法分析和中间代码生成。
    - zend_excute()执行中间代码。
    - 所有操作处理完成，EG()返回结果
- Deactivation：php_request_shutdown()
    - 调用所有通过register_shutdown_function()注册的函数
    - 执行所有可用的__destruct()
    - 将所有输出刷出去
    - 发送HTTP应答头
    - 变量每个模块的 RSHUTDOWN
    - 销毁全局变量表PG(http_globals)的变量
    - 调用zend_deactivate()关闭词法分析器，语法分析器和中间代码执行器
    - 调用每个扩展的post-RSHUTDOWN()
    - 关闭SAPI，调用sapi_deactivate()销毁SG(sapi_headers),SG(request_info)等内容
    - 关闭流的包装器，关闭流的过滤器
    - 关闭内存管理
    - 重新设置最大执行时间。
- 结束
    - flush：调用sapi_module.flush()将最后的内容刷新出去
    - 关闭Zend引擎：zend_shutdown()

#### 多进程SAPI生命周期
通常PHP是编译为apacher的一个模块来处理PHP请求。Apache一般会采用多进程模式，在启动后会fork出多个子进程，每个进程的内存空间独立，都会经过开始和结束阶段。不同的是开始阶段正在进程fork出来以后进行，只有在Apache关闭后进程被结束之后才会进行关闭阶段，中间会处理多个请求，并行重复请求开始-请求结束，

### Zend引擎
Zend引擎是PHP实现的核心，提供了语言实现上的基础设施，而PHP提供请求处理和SAPI。

目前PHP的实现和Zend引擎中间的关系非常紧密，甚至有些过于紧密了。比如现在很多PHP扩展使用Zend API，而PHP只是使用Zend这个内核来国建PHP语言，这就导致PHP很多扩展和Zend引擎耦合了。

很多脚本语言会有语言扩展机制，PHP通过Pear库或原始扩展。
