<!--
author: 刘青
date: 2017-03-28
title: 函数的种类
type: note
source: https://github.com/reeze/tipi/blob/master/book/chapt04/04-01-01-function-struct.markdown
tags: 
category: php/src
status: publish
summary: 
-->

在PHP中，函数有：
- 作用域
- 具体逻辑
- 返回值

Zend引擎将函数分为：
```
// @file Zend/zend_compile.h

#define ZEND_INTERNAL_FUNCTION				1   // 由扩展，PHP内核，Zend引擎提供的函数
#define ZEND_USER_FUNCTION					2   // 用户定义的函数
#define ZEND_OVERLOADED_FUNCTION			3
#define	ZEND_EVAL_CODE						4
#define ZEND_OVERLOADED_FUNCTION_TEMPORARY	5
```

PHP 中关于函数的数据结构有：
```
// @file Zend/zend_compile.h
typedef struct _zend_function_state {
	zend_function *function;
	void **arguments;
} zend_function_state;

typedef union _zend_function {
	zend_uchar type;	/* MUST be the first element of this struct! */

	struct {
		zend_uchar type;                    /* never used */
		char *function_name;                // 函数名
		zend_class_entry *scope;            // 函数所在的类作用域
		zend_uint fn_flags;                 // 作为方法时的范围类型
		union _zend_function *prototype;    // 函数原型
		zend_uint num_args;                 // 参数数目
		zend_uint required_num_args;        // 需要的参数数目
		zend_arg_info *arg_info;            // 参数信息
		zend_bool pass_rest_by_reference;   
		unsigned char return_reference;     // 返回值
	} common;

	zend_op_array op_array;                 // 函数中的操作
	zend_internal_function internal_function;
} zend_function;


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

	zend_op *opcodes;
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

### 用户函数（ZEND_USER_FUNCTION）

定义一个用户函数：
```
function tipi($name) {
    $r = "Hi, " . $name;

    return $r;
}
```

PHP中的实现是：
```
// @file Zend/zend_compile.h
// 函数执行的数据保存在下面的结构体中
struct _zend_execute_data {
	struct _zend_op *opline;

	zend_function_state function_state;
	zend_function *fbc; /* Function Being Called */
	
    zend_class_entry *called_scope;
	zend_op_array *op_array;
	zval *object;
	union _temp_variable *Ts;
	zval ***CVs;
	HashTable *symbol_table;
	struct _zend_execute_data *prev_execute_data;
	zval *old_error_reporting;
	zend_bool nested;
	zval **original_return_value;
	zend_class_entry *current_scope;
	zend_class_entry *current_called_scope;
	zval *current_this;
	zval *current_object;
	struct _zend_op *call_opline;
};

```

### 内部函数（ZEND_INTERNAL_FUNCTION）
ZEND_INTERNAL_FUNCTION函数是由扩展、PHP内核、Zend引擎提供的内部函数，一般用“C/C++”编写，可以直接在用户脚本中调用的函数。

```
// @file Zend/zend_compile.h

typedef struct _zend_internal_function {
	/* Common elements */
	zend_uchar type;
	char * function_name;
	zend_class_entry *scope;
	zend_uint fn_flags;
	union _zend_function *prototype;
	zend_uint num_args;
	zend_uint required_num_args;
	zend_arg_info *arg_info;
	zend_bool pass_rest_by_reference;
	unsigned char return_reference;
	/* END of common elements */

	void (*handler)(INTERNAL_FUNCTION_PARAMETERS);  // 调用方法
	struct _zend_module_entry *module;              // 所属模块
} zend_internal_function;
```

最常见的操作是在模块初始化时Zend遍历每个载入的扩展模块，为模块中`function_entry`指明的每一个函数（module->functions）创建一个`zend_internal_function` 结构，并将其type制为`ZEND_INTERNAL_FUNCTION`，写入全局的函数表。

### 变量函数
PHP支持变量函数的概念：一个变量名后面加括号，PHP将寻找与变量的值同名的函数，并尝试执行它。
```
$func = 'print_r';
$func('i am print_r function');
```

### 函数间的转换
对比 `zend_function`和`zend_internal_function`和`zend_op_array`，可以发现它们的共同点特别多。转换关系为：
- zend_function可以和zend_op_array互转
- zend_function可以和zend_internal_function互转
- zend_op_array不可以通过zend_function和zend_internal_function互转
