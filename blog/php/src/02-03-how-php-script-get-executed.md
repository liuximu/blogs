<!--
author: 刘青
date: 2017-03-20
title: PHP脚本的执行
type: note
source: https://github.com/reeze/tipi/blob/master/book/chapt02/02-03-00-how-php-script-get-executed.markdown
tags: 
category: php/src
status: publish 
summary: 
-->

大概的看了PHP的生命周期，PHP的SAPI。这里讨论PHP最后的执行。

目前编程语言分为两大类：
- 编译型：运行之前必须对源代码进行编译，运行的是编译后的目标文件
- 解释型：直接运行。实现为有解释器来（运行时）编译并执行这些代码

PHP时解释型语言。以PHP命令行程序为例解释PHP脚本的执行方式。
```
// 编写 hello.php
<?php
$str = "Hello\n";
echo $str;

// 命令行执行
php hello.php
Hello
```
以上代码的执行步骤包括：
- php程序完成准备工作后启动PHP及Zend引擎，加载注册的扩展模块
- 编译：Zend读取脚本，词法分析，语法分析，编译为opcode（安装了opcode缓存则可以跳过）
- 执行opcode


#### 词法分析
PHP使用lex生成词法解析器。

#### 语法分析
PHP使用bison生成语法解析器。

#### 编译代码为opcode

> [opcode](https://en.wikipedia.org/wiki/Opcode)：operation code。计算机指令的一部分。

opcode又叫字节码。java运行在JVM上，有其专有的字节码；C#运行在.Net平台上，有其专有的CIL；PHP运行在Zend虚拟机中，也有其专有的字节码。

```
// Zend/zend_compile.h

struct _zend_op {
	opcode_handler_t handler;   // 执行该opcode时调用的函数
	znode result;               // 脚本执行结果
	znode op1;
	znode op2;
	ulong extended_value;       // 脚本执行时需要更多的信息
	uint lineno;
	zend_uchar opcode;          // opcode代码
};

//print语句的编译函数：
void zend_do_print(znode *result, const znode *arg TSRMLS_DC) /* {{{ */
{
	zend_op *opline = get_next_op(CG(active_op_array) TSRMLS_CC);

	opline->result.op_type = IS_TMP_VAR;    // 设置返回值的类型为临时变量
	opline->result.u.var = get_temporary_variable(CG(active_op_array));
	opline->opcode = ZEND_PRINT;            // 指定opcode
	opline->op1 = *arg;                     // 设置参数
	SET_UNUSED(opline->op2);
	*result = opline->result;               // 返回执行结果
}
/* }}} */

// echo语句的编译函数
void zend_do_echo(const znode *arg TSRMLS_DC) /* {{{ */
{
	zend_op *opline = get_next_op(CG(active_op_array) TSRMLS_CC);

	opline->opcode = ZEND_ECHO;
	opline->op1 = *arg;
	SET_UNUSED(opline->op2);
}

// 脚本编译为opcode保存在数组中
struct _zend_op_array {
	/* Common elements */
	zend_uchar type;
	char *function_name;		
	zend_class_entry *scope;
	zend_uint fn_flags;
	union _zend_function *prototype;
	zend_uint num_args;
	zend_uint required_num_args;
	zend_arg_info *arg_info;
	zend_bool pass_rest_by_reference;
	unsigned char return_reference;
	/* END of common elements */

	zend_bool done_pass_two;

	zend_uint *refcount;

	zend_op *opcodes;           // opcode数组
	zend_uint last, size;

	zend_compiled_variable *vars;
	int last_var, size_var;

	zend_uint T;

	zend_brk_cont_element *brk_cont_array;
	int last_brk_cont;
	int current_brk_cont;

	zend_try_catch_element *try_catch_array;
	int last_try_catch;

	/* static variables support */
	HashTable *static_variables;

	zend_op *start_op;
	int backpatch_count;

	zend_uint this_var;

	char *filename;
	zend_uint line_start;
	zend_uint line_end;
	char *doc_comment;
	zend_uint doc_comment_len;
	zend_uint early_binding; /* the linked list of delayed declarations */

	void *reserved[ZEND_MAX_RESERVED_RESOURCES];
};
```

#### 执行opcode
这个过程可能会重复进行编译-执行过程。
```
// Zend/zend_vm_execute.h
```
