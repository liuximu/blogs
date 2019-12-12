<!--
author: 刘青
date: 2017-03-17
title: SAPI
type: note
source: https://github.com/reeze/tipi/blob/master/book/chapt02/02-02-00-overview.markdown 
tags: 
category: php/src
status: publish 
summary: 
-->

PHP的生命周期的各个阶段与服务相关的操作都通过SAPI接口实现。代码在src/sapi目录下，包括一些常用的实现：
```
|—— aolserver
├── apache
├── apache2filter
├── apache2handler
├── apache_hooks
├── caudium
├── cgi
├── cli
├── continuity
├── embed
├── fpm
├── isapi
├── litespeed
├── milter
├── nsapi
├── phttpd
├── pi3web
├── roxen
├── tests
├── thttpd
├── tux
└── webjames
```

SAPI的调用流程为：

![SAPI的调用流程](https://raw.githubusercontent.com/reeze/tipi/master/book/images/chapt02/02-02-01-sapi.png)

在各个服务器抽象层中间遵守着相同的约定，就是我们说的SAPI接口。每个SAPI实现都是一个`_sapi_module_struct`结构体变量。服务器相关信息的调用全部通过调用SAPI接口中对应的方法实现，每个服务器抽象层会有各个方法的自己的实现。

我们来查看cgi的SAPI。
```
//启动方法
cgi_sapi_module.startup(&cgi_sapi_module)
```

cgi_sapi_module的变量类型是sapi_module_struct，对应的类型是_sapi_module_struct：
```
// main/SAPI.h

struct _sapi_module_struct {
    //名字（标识用）
	char *name;
    //更好理解的名字
	char *pretty_name;

    //启动函数指针
	int (*startup)(struct _sapi_module_struct *sapi_module);
    //关闭函数指针
	int (*shutdown)(struct _sapi_module_struct *sapi_module);

    //激活
	int (*activate)(TSRMLS_D);
    //停用
	int (*deactivate)(TSRMLS_D);

    //不写缓存的操作 unbuffered write
	int (*ub_write)(const char *str, unsigned int str_length TSRMLS_DC);
    //flush
	void (*flush)(void *server_context);
    
	struct stat *(*get_stat)(TSRMLS_D);
	char *(*getenv)(char *name, size_t name_len TSRMLS_DC);

    //错误处理者
	void (*sapi_error)(int type, const char *error_msg, ...);

    //header处理者
	int (*header_handler)(sapi_header_struct *sapi_header, sapi_header_op_enum op, sapi_headers_struct *sapi_headers TSRMLS_DC);
	int (*send_headers)(sapi_headers_struct *sapi_headers TSRMLS_DC);
	void (*send_header)(sapi_header_struct *sapi_header, void *server_context TSRMLS_DC);

	int (*read_post)(char *buffer, uint count_bytes TSRMLS_DC);
	char *(*read_cookies)(TSRMLS_D);

	void (*register_server_variables)(zval *track_vars_array TSRMLS_DC);
	void (*log_message)(char *message);
	time_t (*get_request_time)(TSRMLS_D);
	void (*terminate_process)(TSRMLS_D);

    //覆盖的ini路径
	char *php_ini_path_override;

	void (*block_interruptions)(void);
	void (*unblock_interruptions)(void);

	void (*default_post_reader)(TSRMLS_D);
	void (*treat_data)(int arg, char *str, zval *destArray TSRMLS_DC);
	char *executable_location;

	int php_ini_ignore;

	int (*get_fd)(int *fd TSRMLS_DC);

	int (*force_http_10)(TSRMLS_D);

	int (*get_target_uid)(uid_t * TSRMLS_DC);
	int (*get_target_gid)(gid_t * TSRMLS_DC);

	unsigned int (*input_filter)(int arg, char *var, char **val, unsigned int val_len, unsigned int *new_val_len TSRMLS_DC);
	
	void (*ini_defaults)(HashTable *configuration_hash);
	int phpinfo_as_text;

	char *ini_entries;
	const zend_function_entry *additional_functions;
	unsigned int (*input_filter_init)(TSRMLS_D);
};
```

这个SAPI类似于一个面向对象中的模板方法模式的应用。SAPI.c和SAPI.h文件所包含的一些函数就是模板方法模式中的抽象模板，刚刚服务器对于sapi_module的定义及相关实现则是一个个具体的模板。

### PHP执行方式之作为子模块：Apache模块
Apache支持许多特性，大部分通过模块扩展实现。一些通用的语言支持以Apache模块的方式与Apache集成，PHP是其中之一。

当PHP需要在Apache服务器下运行时，一般是以mod_php5模块的形式集成。该模块接收Apache传递过来的PHP文件请求，处理，然后返回给Apache。

Apache加载PHP的方式有两种：
- 在Apache启动前进行配置，Apache在启动时会启动该模块接收PHP的请求。
- 在Apache运行时动态加载。想服务器发送HUP|AP_SIG_GRACEFUL信号让去重新载入模块。因为加载需要使用mod_so模块，所以不能动态加载mod_so。动态加载的是动态链接库。

我们看看Apache模块的mod_php5的实现：
```
//sapi/apache2handler/mod_php5.c

AP_MODULE_DECLARE_DATA module php5_module = {
	STANDARD20_MODULE_STUFF,
	create_php_config,		/* create per-directory config structure */
	merge_php_config,		/* merge per-directory config structures */
	NULL,					/* create per-server config structure */
	NULL,					/* merge per-server config structures */
	php_dir_cmds,			/* command apr_table_t */   //这个是指令集合
    //注册钩子，此函数通过ap_hoo_开头的函数在一次请求处理过程中对于指定的步骤注册钩子
	php_ap2_register_hook	/* register hooks */
};
```

变量php_dir_cmds是mod_php5模块定义的指令表。当Apache遇到指令是将逐一遍历各个模块中的指令表。如果有能够处理该指令的模块就调用相应的函数。如果一直没找到就报错。可见该模块就能处理五个指令。
```
const command_rec php_dir_cmds[] =
{
	AP_INIT_TAKE2("php_value", php_apache_value_handler, NULL, OR_OPTIONS, "PHP Value Modifier"),
	AP_INIT_TAKE2("php_flag", php_apache_flag_handler, NULL, OR_OPTIONS, "PHP Flag Modifier"),
	AP_INIT_TAKE2("php_admin_value", php_apache_admin_value_handler, NULL, ACCESS_CONF|RSRC_CONF, "PHP Value Modifier (Admin)"),
	AP_INIT_TAKE2("php_admin_flag", php_apache_admin_flag_handler, NULL, ACCESS_CONF|RSRC_CONF, "PHP Flag Modifier (Admin)"),
	AP_INIT_TAKE1("PHPINIDir", php_apache_phpini_set, NULL, RSRC_CONF, "Directory containing the php.ini file"),
	{NULL}
};
```

变量php_ap2_register_hook：
```
void php_ap2_register_hook(apr_pool_t *p)
{
	ap_hook_pre_config(php_pre_config, NULL, NULL, APR_HOOK_MIDDLE);
	ap_hook_post_config(php_apache_server_startup, NULL, NULL, APR_HOOK_MIDDLE);
	ap_hook_handler(php_handler, NULL, NULL, APR_HOOK_MIDDLE);
	ap_hook_child_init(php_apache_child_init, NULL, NULL, APR_HOOK_MIDDLE);
}
```

在服务器启动时，pre_config,post_config,child_init被调用。其中post_config中将启动php，它通过php_apache_server_startup()调用sapi_startup()启动sapi来实现。

至此，我们知道了Apache加载mod_php5模块的整个过程。我们再看看Apache和SAPI的关系。
在mod_php5模块中定义了属于Apache的sapi_module_struct结构：

```
static sapi_module_struct apache2_sapi_module = {
	"apache2handler",
	"Apache 2.0 Handler",

	php_apache2_startup,				    /* startup */
	php_module_shutdown_wrapper,			/* shutdown */

	NULL,						            /* activate */
	NULL,						            /* deactivate */

	php_apache_sapi_ub_write,			    /* unbuffered write */
	php_apache_sapi_flush,				    /* flush */
	php_apache_sapi_get_stat,			    /* get uid */
	php_apache_sapi_getenv,				    /* getenv */

	php_error,					            /* error handler */

	php_apache_sapi_header_handler,			/* header handler */
	php_apache_sapi_send_headers,			/* send headers handler */
	NULL,						            /* send header handler */

	php_apache_sapi_read_post,			    /* read POST data */
	php_apache_sapi_read_cookies,			/* read Cookies */

	php_apache_sapi_register_variables,
	php_apache_sapi_log_message,			/* Log message */
	php_apache_sapi_get_request_time,		/* Request Time */
	NULL,						            /* Child Terminate */

	STANDARD_SAPI_MODULE_PROPERTIES
};
```

里面定义了各种行为，最终会调Apache对应的函数，其实现依赖于Apache的实现。
比如`php_apache_sapi_flush`，最终会调用：
```
// main/SAPI.c

SAPI_API int sapi_flush(TSRMLS_D)
{
	if (sapi_module.flush) {
		sapi_module.flush(SG(server_context));
		return SUCCESS;
	} else {
		return FAILURE;
	}
}
```

#### Apache的运行过程
- 启动阶段：以单进程单线程完成启动，包括：配置文件解析、模块加载、系统资源初始化等，使用获取到的最大权限。
- 运行阶段：Apache主要是处理用户的服务请求。使用普通权限。对HTTP请求分为三大阶段：连接，处理和断开。细化为11个小阶段：
    - Post-Read-Request
    - URI Translation
    - Header Parsing
    - Access Control
    - Authentication
    - Authorization
    - MIME Type Checking
    - FixUp
    - Response
    - Logging
    - CleanUp

#### Apache Hook机制
Apache允许模块（内部|外部）将自定义的函数注入到请求循环中。也就是说模块可以在Apache的任何一个处理阶段字段接(Hook)上自己的处理函数。


### PHP执行方式之作为嵌入式
嵌入式类似CLI。一般情况下她的一个请求的生命周期也会和其他的SAPI一样：模块初始化->请求初始化->处理请求->关闭请求->关闭模块。

对于嵌入式PHP我们了解比较少，或者说根本用不到。很多游戏中使用Lu语言作为粘合语言，或者作为扩展游戏的脚本语言，而JavaScript则是嵌入在浏览器中的语言。但是很少会将PHP嵌入到哪里，PHP主要还是用在Web开发。

不详细说了。sapi/embed/ 下有例子。

### PHP执行方式之作为独立程序：FastCGI
> CGI：Common Gateway Interface。通用网管接口。描述了客户端和服务器之间传输数据的标准。

CGI的运行原理为：
- 客户端访问某个URL地址之后，通过HTTP协议想Web服务器发出请求；
- 服务器端的HTTP Daemon启动一个子进程，将HTTP请求描述的信息通过标准输入stdin和环境变量传递给URL指定的CGI程序，启动此应用程序进行处理，将结果通过stdout返回给HTTP Daemon子进程（fork-and-execute）；
- HTTP Daemon 子进程通过HTTP协议返回给客户端。

![](https://raw.githubusercontent.com/reeze/tipi/master/book/images/chapt02/02-02-03-cgi.png)

> FastCGI： CGI的改进方案，常驻型（long-lived）的CGI。它一直执行，将CGI解析器进程保持在内存中，在请求到达时不会花费时间去fork一个进程来处理。

FastCGI的工作流程为：
- FastCGI进程管理器自身初始化，启动多个CGI解析器进程，并等待来着Web Server的连接；
- Web Server与FastCGI进程管理器进行Socket通信，通过FastCGI协议发送CGI环境变量和stdin数据给FastCGI创建的CGI解析器子进程；
- CGI解析器进程完成处理后将标准输出和错误信息从同一连接返回Web Server；
- CGI解析器进程等待下一个Web Server的连接；

![](https://raw.githubusercontent.com/reeze/tipi/master/book/images/chapt02/02-02-03-fastcgi-demo.png)


我们看看PHP的FastCGI的代码：
#### FastCGI消息类型
FastCGI将传输的消息做了很多类型的划分，结构体定义为：
```
// main/sapi/cgi/fastcgi.h
typedef enum _fcgi_request_type {
	FCGI_BEGIN_REQUEST		=  1, /* [in]                              */
	FCGI_ABORT_REQUEST		=  2, /* [in]  (not supported)             */
	FCGI_END_REQUEST		=  3, /* [out]                             */
	FCGI_PARAMS				=  4, /* [in]  environment variables       */
	FCGI_STDIN				=  5, /* [in]  post data                   */
	FCGI_STDOUT				=  6, /* [out] response                    */
	FCGI_STDERR				=  7, /* [out] errors                      */
	FCGI_DATA				=  8, /* [in]  filter data (not supported) */
	FCGI_GET_VALUES			=  9, /* [in]                              */
	FCGI_GET_VALUES_RESULT	= 10  /* [out]                             */
} fcgi_request_type;
```

#### 消息的发送顺序

![](https://raw.githubusercontent.com/reeze/tipi/master/book/images/chapt02/02-02-03-fastcgi-data.png)

- FCGI_BEGIN_REQUEST：标志请求开始
- FCGI_PARAMS：二进制形式，每次最多65535字节反复发送
- FCGI_STDIN：同上
- FCGI_STDOUT：同上
- FCGI_STDERR：同上
- FCGI_END_REQUEST：请求结束

举个例子：
```
{FCGI_BEGIN_REQUEST,   1, {FCGI_RESPONDER, 0}}
{FCGI_PARAMS,          1, "\013\002SERVER_PORT80\013\016SERVER_ADDR199.170.183.42 ... "}
{FCGI_STDIN,           1, "quantity=100&item=3047936"}
{FCGI_STDOUT,          1, "Content-type: text/html\r\n\r\n<html>\n<head> ... "}
{FCGI_END_REQUEST,     1, {0, FCGI_REQUEST_COMPLETE}}
```

FCGI_BEGIN_REQUEST 和 FCGI_END_REQUEST两种消息类型的消息时协议的一部分，有对应的数据结构。

```
typedef enum _fcgi_role {
	FCGI_RESPONDER	= 1,
	FCGI_AUTHORIZER	= 2,
	FCGI_FILTER		= 3
} fcgi_role; //Web服务器期望应用扮演的角色

typedef struct _fcgi_begin_request {
	unsigned char roleB1;       // 如上
	unsigned char roleB0;
	unsigned char flags;
	unsigned char reserved[5];
} fcgi_begin_request;


typedef enum _fcgi_protocol_status {
	FCGI_REQUEST_COMPLETE	= 0,    // 正常结束
	FCGI_CANT_MPX_CONN		= 1,    // 拒绝新请求。这发生在Web服务器通过一条线路向应用发送并发的请求时，后者被设计为每条线路每次处理一个请求。
	FCGI_OVERLOADED			= 2,    // 拒绝新请求。这发生在应用用完某些资源时，例如数据库连接。
	FCGI_UNKNOWN_ROLE		= 3     // 拒绝新请求。这发生在Web服务器指定了一个应用不能识别的角色时。
} dcgi_protocol_status;

typedef struct _fcgi_end_request {
    unsigned char appStatusB3;  // 应用级别的状态码
    unsigned char appStatusB2;
    unsigned char appStatusB1;
    unsigned char appStatusB0;
    unsigned char protocolStatus;   //协议级别的状态码
    unsigned char reserved[3];
} fcgi_end_request;
```

#### FastCGI消息头
消息是二进制连续传递的，必须定义一个统一的结构的消息头来方便消息体的读取和切割。
```
// main/sapi/cgi/fastcgi.h

typedef struct _fcgi_header {
	unsigned char version;              // FastCGI的协议版本
	unsigned char type;                 
	unsigned char requestIdB1;          // 请求标记
	unsigned char requestIdB0;
	unsigned char contentLengthB1;      // 内容的字节数
	unsigned char contentLengthB0;
	unsigned char paddingLength;
	unsigned char reserved;
} fcgi_header;
```

#### PHP中的CGI实现
PHP的CGI实现了FastCGI协议，是一个TCP|UDP协议的服务器，通过socket监听来自Web服务器的请求，然后进入PHPd-生命周期进行处理。

以TCP为例，大致的步骤为：
- 调用socket函数创建一个TCP用的流式套接字S
- 调用bind函数将服务器的本地地址和S绑定
- 调用listen函数将S作为监听，等待客户端发起请求
- 服务器进程调用accept函数进入阻塞状态，知道有客户进程调用connect函数而建立起一个连接
- 服务器调用read_stream函数读取客户请求
- 处理数据
- 服务器调用write函数向客户端发送应答

php的实现代码在： main/sapi/cgi_main.c
