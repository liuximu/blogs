<!--
author: 刘青
date: 2017-03-27
title: 预定义变量
type: note
source: https://github.com/reeze/tipi/blob/master/book/chapt03/03-03-pre-defined-variable.markdown
tags: 
category: php/src
status: publish 
summary: 
-->
在PHP脚本执行的时候，用户全局变量（在用户控件显式定义的变量）会保存在一个数据类型HasTable的符号表（symbol_table）中。

> 预定义变量：PHP在脚本运行前就加入到符号表的变量。

比如： $_GET, $_POST, $_SERVER, $_FILES。

### $GLOBALS的初始化
以cgi模式为例说明$GLOBALS的初始化。
```
// @file: Zend/zend_execute_API.c

// [main() -> php_request_startup() -> zend_activate() -> init_executor() ]

void init_executor(TSRMLS_D) /* {{{ */
{
    // ...

	zend_hash_init(&EG(symbol_table), 50, NULL, ZVAL_PTR_DTOR, 0);
	{
		zval *globals;

		ALLOC_ZVAL(globals);
		Z_SET_REFCOUNT_P(globals, 1);
		Z_SET_ISREF_P(globals);
		Z_TYPE_P(globals) = IS_ARRAY;
		Z_ARRVAL_P(globals) = &EG(symbol_table);
		zend_hash_update(&EG(symbol_table), "GLOBALS", sizeof("GLOBALS"), &globals, sizeof(zval *), NULL);
	}

    // ...
```

EG(symbol_table)是一个HashTable，用来存放顶层作用域的变量。 zend_hash_update函数会将其$GLOBALS注册进去。
php_request_startup函数会在PHP的请求初始化生命周期阶段执行，对于每个用户请求$GLOBAL都会重新载入。

### $GET|$POST等变量的初始化

要讨论的变量有：
- $_GET
- $_COOKIE
- $_SERVCER
- $_ENV
- $_FILES
- $_POST

函数调用链为：main() -> php_request_startup() -> php_hash_environment()
```
// @file:main/php_variables.c
/* {{{ php_hash_environment
 */
int php_hash_environment(TSRMLS_D)
{
	char *p;
	unsigned char _gpc_flags[5] = {0, 0, 0, 0, 0};
	zend_bool jit_initialization = (PG(auto_globals_jit) && !PG(register_globals) && !PG(register_long_arrays));
	struct auto_global_record {
		char *name;
		uint name_len;
		char *long_name;
		uint long_name_len;
		zend_bool jit_initialization;
	} auto_global_records[] = {
		{ "_POST", sizeof("_POST"), "HTTP_POST_VARS", sizeof("HTTP_POST_VARS"), 0 },
		{ "_GET", sizeof("_GET"), "HTTP_GET_VARS", sizeof("HTTP_GET_VARS"), 0 },
		{ "_COOKIE", sizeof("_COOKIE"), "HTTP_COOKIE_VARS", sizeof("HTTP_COOKIE_VARS"), 0 },
		{ "_SERVER", sizeof("_SERVER"), "HTTP_SERVER_VARS", sizeof("HTTP_SERVER_VARS"), 1 },
		{ "_ENV", sizeof("_ENV"), "HTTP_ENV_VARS", sizeof("HTTP_ENV_VARS"), 1 },
		{ "_FILES", sizeof("_FILES"), "HTTP_POST_FILES", sizeof("HTTP_POST_FILES"), 0 },
	};
	size_t num_track_vars = sizeof(auto_global_records)/sizeof(struct auto_global_record);
	size_t i;

	/* jit_initialization = 0; */
	for (i=0; i<num_track_vars; i++) {
		PG(http_globals)[i] = NULL;
	} 
	for (p=PG(variables_order); p && *p; p++) {
		switch(*p) {
			case 'p':
			case 'P':
				if (!_gpc_flags[0] && !SG(headers_sent) && SG(request_info).request_method && !strcasecmp(SG(request_info).request_method, "POST")) {
					sapi_module.treat_data(PARSE_POST, NULL, NULL TSRMLS_CC);	/* POST Data */
					_gpc_flags[0] = 1;
					if (PG(register_globals)) {
						php_autoglobal_merge(&EG(symbol_table), Z_ARRVAL_P(PG(http_globals)[TRACK_VARS_POST]) TSRMLS_CC);
					}
				}
				break;
			case 'c':
			case 'C':
				if (!_gpc_flags[1]) {
					sapi_module.treat_data(PARSE_COOKIE, NULL, NULL TSRMLS_CC);	/* Cookie Data */
					_gpc_flags[1] = 1;
					if (PG(register_globals)) {
						php_autoglobal_merge(&EG(symbol_table), Z_ARRVAL_P(PG(http_globals)[TRACK_VARS_COOKIE]) TSRMLS_CC);
					}
				}
				break;
			case 'g':
			case 'G':
				if (!_gpc_flags[2]) {
					sapi_module.treat_data(PARSE_GET, NULL, NULL TSRMLS_CC);	/* GET Data */
					_gpc_flags[2] = 1;
					if (PG(register_globals)) {
						php_autoglobal_merge(&EG(symbol_table), Z_ARRVAL_P(PG(http_globals)[TRACK_VARS_GET]) TSRMLS_CC);
					}
				}
				break;
			case 'e':
			case 'E':
				if (!jit_initialization && !_gpc_flags[3]) {
					zend_auto_global_disable_jit("_ENV", sizeof("_ENV")-1 TSRMLS_CC);
					php_auto_globals_create_env("_ENV", sizeof("_ENV")-1 TSRMLS_CC);
					_gpc_flags[3] = 1;
					if (PG(register_globals)) {
						php_autoglobal_merge(&EG(symbol_table), Z_ARRVAL_P(PG(http_globals)[TRACK_VARS_ENV]) TSRMLS_CC);
					}
				}
				break;
			case 's':
			case 'S':
				if (!jit_initialization && !_gpc_flags[4]) {
					zend_auto_global_disable_jit("_SERVER", sizeof("_SERVER")-1 TSRMLS_CC);
					php_register_server_variables(TSRMLS_C);
					_gpc_flags[4] = 1;
					if (PG(register_globals)) {
						php_autoglobal_merge(&EG(symbol_table), Z_ARRVAL_P(PG(http_globals)[TRACK_VARS_SERVER]) TSRMLS_CC);
					}
				}
				break;
		}
	}

	/* argv/argc support */
	if (PG(register_argc_argv)) {
		php_build_argv(SG(request_info).query_string, PG(http_globals)[TRACK_VARS_SERVER] TSRMLS_CC);
	}

	for (i=0; i<num_track_vars; i++) {
		if (jit_initialization && auto_global_records[i].jit_initialization) {
			continue;
		}
		if (!PG(http_globals)[i]) {
			ALLOC_ZVAL(PG(http_globals)[i]);
			array_init(PG(http_globals)[i]);
			INIT_PZVAL(PG(http_globals)[i]);
		}

		Z_ADDREF_P(PG(http_globals)[i]);
		zend_hash_update(&EG(symbol_table), auto_global_records[i].name, auto_global_records[i].name_len, &PG(http_globals)[i], sizeof(zval *), NULL);
		if (PG(register_long_arrays)) {
			zend_hash_update(&EG(symbol_table), auto_global_records[i].long_name, auto_global_records[i].long_name_len, &PG(http_globals)[i], sizeof(zval *), NULL);
			Z_ADDREF_P(PG(http_globals)[i]);
		}
	}

	/* Create _REQUEST */
	if (!jit_initialization) {
		zend_auto_global_disable_jit("_REQUEST", sizeof("_REQUEST")-1 TSRMLS_CC);
		php_auto_globals_create_request("_REQUEST", sizeof("_REQUEST")-1 TSRMLS_CC);
	}

	return SUCCESS;
}
/* }}} */
```
- 首先以 `auto_global_record` 数组形式定义好将要初始化的变量的相关信息
- 然后按照PG(variables_orders)指定的顺序（php.ini中指定），调用 `sapi_module.treat_data`处理数据
- 如果打开了 `PG(register_globals)`，将数据同时写入符号表
- 将数据更新到 `symbol_table` 中

以 $_POST为例：
当客户端发情请求是，Apache会将收到的内容转交给mod_php5模块。当PHP收到请求后先调用 `sapi_activate` 来根据请求的方法处理数据：
```
//@file: main/SAPI.c

/*
 * Called from php_request_startup() for every request.
 */
SAPI_API void sapi_activate(TSRMLS_D)
{
	zend_llist_init(&SG(sapi_headers).headers, sizeof(sapi_header_struct), (void (*)(void *)) sapi_free_header, 0);
	SG(sapi_headers).send_default_content_type = 1;

	/*
	SG(sapi_headers).http_response_code = 200;
	*/
	SG(sapi_headers).http_status_line = NULL;
	SG(sapi_headers).mimetype = NULL;
	SG(headers_sent) = 0;
	SG(read_post_bytes) = 0;
	SG(request_info).post_data = NULL;
	SG(request_info).raw_post_data = NULL;
	SG(request_info).current_user = NULL;
	SG(request_info).current_user_length = 0;
	SG(request_info).no_headers = 0;
	SG(request_info).post_entry = NULL;
	SG(request_info).proto_num = 1000; /* Default to HTTP 1.0 */
	SG(global_request_time) = 0;

	/* It's possible to override this general case in the activate() callback, if
	 * necessary.
	 */
	if (SG(request_info).request_method && !strcmp(SG(request_info).request_method, "HEAD")) {
		SG(request_info).headers_only = 1;
	} else {
		SG(request_info).headers_only = 0;
	}
	SG(rfc1867_uploaded_files) = NULL;

	/* handle request mehtod */
	if (SG(server_context)) {
		if ( SG(request_info).request_method) {
			if(!strcmp(SG(request_info).request_method, "POST")
			   && (SG(request_info).content_type)) {
				/* HTTP POST -> may contain form data to be read into variables
				   depending on content type given
				*/
				sapi_read_post_data(TSRMLS_C);
			} else {
				/* any other method with content payload will fill 
				   $HTTP_RAW_POST_DATA if enabled by always_populate_raw_post_data 
				   it is up to the webserver to decide whether to allow a method or not
				*/
				SG(request_info).content_type_dup = NULL;
				if(sapi_module.default_post_reader) {
					sapi_module.default_post_reader(TSRMLS_C);
				}
			}
		} else {
			SG(request_info).content_type_dup = NULL;
		}

		/* Cookies */
		SG(request_info).cookie_data = sapi_module.read_cookies(TSRMLS_C);
		if (sapi_module.activate) {
			sapi_module.activate(TSRMLS_C);
		}
	}
	if (sapi_module.input_filter_init ) {
		sapi_module.input_filter_init(TSRMLS_C);
	}
}

static void sapi_read_post_data(TSRMLS_D)
{
	sapi_post_entry *post_entry;
	uint content_type_length = strlen(SG(request_info).content_type);
	char *content_type = estrndup(SG(request_info).content_type, content_type_length);
	char *p;
	char oldchar=0;
	void (*post_reader_func)(TSRMLS_D) = NULL;


	/* dedicated implementation for increased performance:
	 * - Make the content type lowercase
	 * - Trim descriptive data, stay with the content-type only
	 */
	for (p=content_type; p<content_type+content_type_length; p++) {
		switch (*p) {
			case ';':
			case ',':
			case ' ':
				content_type_length = p-content_type;
				oldchar = *p;
				*p = 0;
				break;
			default:
				*p = tolower(*p);
				break;
		}
	}

	/* now try to find an appropriate POST content handler */
	if (zend_hash_find(&SG(known_post_content_types), content_type,
			content_type_length+1, (void **) &post_entry) == SUCCESS) {
		/* found one, register it for use */
		SG(request_info).post_entry = post_entry;
		post_reader_func = post_entry->post_reader;
	} else {
		/* fallback */
		SG(request_info).post_entry = NULL;
		if (!sapi_module.default_post_reader) {
			/* no default reader ? */
			SG(request_info).content_type_dup = NULL;
			sapi_module.sapi_error(E_WARNING, "Unsupported content type:  '%s'", content_type);
			return;
		}
	}
	if (oldchar) {
		*(p-1) = oldchar;
	}

	SG(request_info).content_type_dup = content_type;

	if(post_reader_func) {
		post_reader_func(TSRMLS_C);
	}

	if(sapi_module.default_post_reader) {
		sapi_module.default_post_reader(TSRMLS_C);
	}
}
```

### 预定义变量的获取
在某个局部函数中使用类似于$GLOBALS变量来预定义变量，这些变量集中存储在：EG(symbol_table)。
在通过$获取变量时，PHP内核都会通过这些变量名是否为全局变量，如果是直接从EG(symbol_table)中返回数据。
