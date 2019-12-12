<!--
author: 刘青
date: 2017-1-30
title: 常用代码
type: translate
source: http://www.php-internals.com/book/?p=chapt01/01-03-comm-code-in-php-src
tags: 
category: php/src
status: publish 
summary: 
-->

PHP中有许多常见的宏非常高频，或许对于刚刚接触源码的读者非常难懂。但是每个模块基本上都有他们的身影。本小节提取一些进行说明。

## 1. "\#\#"和"\#"
宏是C/C++非常强大且高频的一个功能，有时用来实现类似函数内敛的效果，或者将复杂的代码进行简单封装，提供可读性或可移植性。PHP宏定义包括##和#。

宏的使用可以参考：[预处理器](/blog/clang/c_programming/14-pre-processor.html)

先看双井号

```
#define PHP_FUNCTION ZEND_FUNCTION
#define ZEND_FN(name) zif_##name
#define ZEND_FUNCTION(name) ZEND_NAMED_FUNCTION(ZEND_FN(name))
#define ZEND_NAMED_FUNCTION(name) void name(INTERNAL_FUNCTION_PARAMETERS)
#define INTERNAL_FUNCTION_PARAMETERS int ht, zval *return_value, zval **return_value_ptr, \
zval *this_ptr, int return_user_used TSRMLS_DC

//对于
PHP_FUNCTION(count);
//预处理器得到的代码是：
ZEND_FUNCTION(count);
ZEND_NAMED_FUNCTION(ZEND_FN(count));
ZEND_NAMED_FUNCTION(zif_count);
void zif_count(int ht, zval *return_value, zval **return_value_ptr, 
zval *this_ptr, int return_value_used TSRMLS_DC)
```

以连接的方式作为基础多次使用宏进行展开可以一定程度减少代码密度。

再看单引号

```
#define STR(x) #x

//对于
printf("%s\n", STR(ok));
//得到
printf("%s\n", #ok);
printf("%s\n", "ok");
```

### 2.宏定义中的do-while
有一种公认的宏编写方式，在宏中加上do-while(0)：

```
#define ALLOC_ZVAL(z)                               \
do {                                                \
    (z) = (zval*)emalloc(sizeof(zval_gc_info));     \
    GC_ZVAL_INTIT(z);                               \
} while(0)                  
```

将摆明了只会运行一次的代码放到循环体内是为了适应多语句的情况：

```
#defin TEST(a, b) a++;b++;

//对于
if(expr)
    TEST(a, b);
else
    do_else();
//得到
if(expr)
    a++;
    b++;        //程序挂了
else
    do_else();
```

### 3. \#line 预处理

```
#line 838 "Zend/zend_language_scanner.c"
```
是为了告诉编译器在编译过程中固定文件名，便于进行调试分析

### 4. PHP中的全局变量宏
PHP中经常能看到一下PG()之类的函数，其实是宏。
举个例子：

```
//src/main/php_globals.h

#ifdef ZTS          //编译时开启了线程安全则使用现场安全库
# define PG(v) TSRMG(core_globals_id, php_core_globals *, v)
extern PHPAPI int core_globals_id;
#else
# define PG(v) (core_globals.v)     //否则它就是一个普通的全局变量
extern ZEND_API struct _php_core_globals core_globals;
#endif
```

再简单说说PHP运行时的一些全局参数，都定义在结构体中，大部分和php.ini的配置对于。PHP在启动时会读取并赋值，使用int_get()和ini_set()可以更新。
```
//src/main/php_globals.h

struct _php_core_globals {
    //是否对虽然的GET/POST/Cookie数据使用自动字符串转义
	zend_bool magic_quotes_gpc;
    //是否对运行时从外部资源产生的数据使用字符串自动转义
	zend_bool magic_quotes_runtime;
    //是否采用Sybase形式的自动字符串转义
	zend_bool magic_quotes_sybase;

    //是否开启安全模式
	zend_bool safe_mode;

    //是否强迫在函数调用时按引用传递参数
	zend_bool allow_call_time_pass_reference;
    //是否要求PHP输出层在每个输出块之后自动刷新数据
	zend_bool implicit_flush;

    //输出缓存的大小（字节）
	long output_buffering;

    //在安全模式下，改组目录和子目录下的文件被包含是，将跳过UID/GID检查
	char *safe_mode_include_dir;
    //在安全模式下，默认在访问文件是会做UID检查
	zend_bool safe_mode_gid;
	zend_bool sql_safe_mode;
    //是否允许使用dl()函数。dl()仅在将PHP作为apache模块安装时才有效
	zend_bool enable_dl;

    //将所有脚本的输出重定向到一个输出处理函数
	char *output_handler;

    //如果解序列号处理器需要实例化一个未定义的类，这里指定的回调函数将以改未定义类的名字作为参数被unserialize()调用
	char *unserialize_callback_func;
    //将浮点型和双精度数据序列号存储是的精度
	long serialize_precision;

    //在安全模式下，只有该目录下的可执行程序才允许被执行系统程序的函数执行
	char *safe_mode_exec_dir;

    //一个脚本能够事情到的最大内存字节数（K|M作为单位）
	long memory_limit;
    //每个脚本解析虽然数据（POST,GET,upload）的最大允许时间（秒）
	long max_input_time;

    //是否在变量$php_errormsg中保持最近一个错误或警告消息
	zend_bool track_errors;
    //是否将错误信息作为输出的一部分显示
	zend_bool display_errors;
    //是否显示PHP启动时的错误
	zend_bool display_startup_errors;
    //是否在日志文件中记录错误
	zend_bool log_errors;
    //设置错误日志中附加的与错误信息相关的错误源的最大长度
	long      log_errors_max_len;
    //记录错误日志是是否忽略重复的错误
	zend_bool ignore_repeated_errors;
    //是否在忽略重复的错误信息是忽略重复的错误源
	zend_bool ignore_repeated_source;
    //是否报告内存泄漏
	zend_bool report_memleaks;
    //错误日志的位置
	char *error_log;

    //PHP的根路径
	char *doc_root;
    //告诉php在使用 /~username 打开脚本时到哪个目录下去找
	char *user_dir;
    //指定一组目录用于require(), include(), fopen_with_path()函数寻找文件
	char *include_path;
    //将PHP允许操作的所有文件（包括文件自身）都限制在此组目录列表下
	char *open_basedir;
    //吃饭扩展库（模块）的目录
	char *extension_dir;

    //文件上传是存放文件的临时目录
	char *upload_tmp_dir;
    //允许上传的文件的最大尺寸
	long upload_max_filesize;
	
    //用于错误信息后输出的字符串
	char *error_append_string;
    //用于错误信息前输出的字符串
	char *error_prepend_string;

    //指定在主文件之前自动解析的文件名
	char *auto_prepend_file;
    //指定在主文件之后自动解析的文件名
	char *auto_append_file;

    //PHP所产生的URL中很用力分隔参数的分隔符
	arg_separators arg_separator;

    //PHP注册 Environment, GET, POST, Cookie, Server 变量的顺序
	char *variables_order;

    //RFC1867保护的变量名，在main/rfc1867.c文件中有用到此变量
	HashTable rfc1867_protected_variables;

    //连接状态 = {正常，中断，超时}
	short connection_status;
    //是否即使在用户中止请求后也检查完成整个请求
	short ignore_user_abort;

    //是否头信息正在发送
	unsigned char header_is_being_sent;

    //仅在main目录下的php_ticks.c文件中有用到，此处定义的函数在register_tick_function等函数中有用到。
	zend_llist tick_functions;

    //存放 GET POST SERVER等信息
	zval *http_globals[6];

    //是否展示php的信息
	zend_bool expose_php;

    //是否想 E, G, P, C, S 变量注册为全局变量
	zend_bool register_globals;
    //是否启用就式的长式数组（HTTP_*_VARS）
	zend_bool register_long_arrays;
    //是否声明$argv和$argc全局变量(包含用GET方法的信息)。
	zend_bool register_argc_argv;
    //是否仅在使用到$_SERVER和$_ENV变量时才创建(而不是在脚本一启动时就自动创建)。
	zend_bool auto_globals_jit;

    //是否强制打开2000年适应(可能在非Y2K适应的浏览器中导致问题)。
	zend_bool y2k_compliance;

    //如果打开了html_errors指令，PHP将会在出错信息上线上超链接
	char *docref_root;
    //指定文件的扩展名（必须含有 ‘.’）
	char *docref_ext;

    //是否在出错信息信息中使用HTML标记
	zend_bool html_errors;
	zend_bool xmlrpc_errors;

	long xmlrpc_error_number;

	zend_bool activated_auto_globals[8];

    //是否已经激活模块
	zend_bool modules_activated;
    //是否允许HTTP文件上传
	zend_bool file_uploads;
	zend_bool during_request_startup;
    //是否允许打开远程文件
	zend_bool allow_url_fopen;
    //是否总是生成$HTTP_RAW_POST_DATA变量(原始POST数据)。
	zend_bool always_populate_raw_post_data;
    // 是否打开zend debug，仅在main/main.c文件中有使用。
	zend_bool report_zend_debug;

    //最后的错误类型
	int last_error_type;
    //最后的错误信息
	char *last_error_message;
    //最后的错误文件
	char *last_error_file;
    //最后错误行
	int  last_error_lineno;

    //该指令接受一个用逗号分隔的函数名列表，以禁用特定的函数。
	char *disable_functions;
    //该指令接受一个用逗号分隔的函数名列表，以禁用特定的类
	char *disable_classes;
    //是否允许include/require远程文件。
	zend_bool allow_url_include;
    //是否超时就退出
	zend_bool exit_on_timeout;
#ifdef PHP_WIN32
	zend_bool com_initialized;
#endif
    //最大的嵌套层数
	long max_input_nesting_level;
    //是否在用户包含空间
	zend_bool in_user_include;

    //用户的ini文件名
	char *user_ini_filename;
    //ini缓存过期现在
	long user_ini_cache_ttl;

    //优先级比variables_order高，在request变量生成时用到，个人觉得是历史遗留问题
	char *request_order;

    //仅在ext/standard/mail.c文件中使用
	zend_bool mail_x_header;
	char *mail_log;

	zend_bool in_error_log;

#ifdef PHP_WIN32
	zend_bool windows_show_crt_warning;
#endif

	long max_input_vars;
};
```
