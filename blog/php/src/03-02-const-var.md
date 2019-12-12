<!--
author: 刘青
date: 2017-03-26
title: PHP常量
type: note
source: https://github.com/reeze/tipi/blob/master/book/chapt03/03-02-const-var.markdown
tags: 
category: php/src
status: publish 
summary: 
-->

常量是特殊的变量。


### define定义常量的过程
```
// @file: Zend/zend_constants.h
[c]
typedef struct _zend_constant {
	zval value;         // 变量本身
	int flags;          // 一些标志
	char *name;         // 常量名
	uint name_len;
	int module_number;  // 模块号
} zend_constant;


// flags
#define CONST_CS				(1<<0)				/* Case Sensitive */
#define CONST_PERSISTENT		(1<<1)				/* Persistent */
#define CONST_CT_SUBST			(1<<2)				/* Allow compile-time substitution */

[php]
define('TIPI', 'Thinking In PHP Internal');


/* {{{ proto bool define(string constant_name, mixed value, boolean case_insensitive=false)
   Define a new constant */
ZEND_FUNCTION(define)
{
	char *name;
	int name_len;
	zval *val;
	zval *val_free = NULL;
	zend_bool non_cs = 0;
	int case_sensitive = CONST_CS;
	zend_constant c;

    // 为变量赋值
	if (zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "sz|b", &name, 
        &name_len, &val, &non_cs) == FAILURE) {
		return;
	}

	if(non_cs) {
		case_sensitive = 0;
	}

	/* class constant, check if there is name and make sure class is valid & exists */
	if (zend_memnstr(name, "::", sizeof("::") - 1, name + name_len)) {
		zend_error(E_WARNING, "Class constants cannot be defined or redefined");
		RETURN_FALSE;
	}

repeat:
	switch (Z_TYPE_P(val)) {
		case IS_LONG:
		case IS_DOUBLE:
		case IS_STRING:
		case IS_BOOL:
		case IS_RESOURCE:
		case IS_NULL:
			break;
		case IS_OBJECT:     // 对象处理
			if (!val_free) {
				if (Z_OBJ_HT_P(val)->get) {
					val_free = val = Z_OBJ_HT_P(val)->get(val TSRMLS_CC);
					goto repeat;
				} else if (Z_OBJ_HT_P(val)->cast_object) {
					ALLOC_INIT_ZVAL(val_free);
					if (Z_OBJ_HT_P(val)->cast_object(val, val_free, IS_STRING TSRMLS_CC) == SUCCESS) {
						val = val_free;
						break;
					}
				}
			}
			/* no break */
		default:
			zend_error(E_WARNING,"Constants may only evaluate to scalar values");
			if (val_free) {
				zval_ptr_dtor(&val_free);
			}
			RETURN_FALSE;
	}
	
    // 创建结构体对象
	c.value = *val;
	zval_copy_ctor(&c.value);
	if (val_free) {
		zval_ptr_dtor(&val_free);
	}
	c.flags = case_sensitive; /* non persistent */
	c.name = zend_strndup(name, name_len);
	c.name_len = name_len+1;
	c.module_number = PHP_USER_CONSTANT;
	if (zend_register_constant(&c TSRMLS_CC) == SUCCESS) {
		RETURN_TRUE;
	} else {
		RETURN_FALSE;
	}
}
/* }}} */
```

### 标准常量的初始化
通过define()函数定义的常量的模块编号都是 `PHP_USER_CONSTANT`，表示用户定义。PHP有内置常量，如错误报警级别 `E_ALL`，为标准常量。
标准常量会在Zend引擎启动后进行注册。

### 魔术常量
魔术常量的值是变化的，不区分大小写，会在语法解析时将常量的内容赋值进行替换。更像是一个占位符。

- __LINE__
- __DIR__
- __FUNCTION__
- __CLASS__
- __METHOD__
-__NAMESPACE__
