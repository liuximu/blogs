<!--
author: 刘青
date: 2017-04-03
title: 类的结构和实现
type: note
source: tags: 
category: php/src
status: publish
summary: 
-->
类和函数类似，PHP内置及PHP扩展均可以实现自己的内部类，也可以由用户使用PHP代码进行定义。

### 类的数据结构
先看用户定义类：
```
class ParentClass 
{

}

interface Ifce
{
    public function iMethod();
}

funal class Tipi extends ParentClass implements Ifce
{
    public static $sa = 'aaa';
    const CA = 'bbb';

    public function __construct()
    {

    }

    public function iMethod()
    {

    }

    private function _access()
    {

    }

    public static function access()
    {
        
    }
}
```

看看类的内部存储结构：
```
//@file: Zend/zend.h
struct _zend_class_entry {
	char type;      //类型：ZEND_INTERNAL_CLASS|ZEND_USER_CLASS
	char *name;     //类名
	zend_uint name_length;
	struct _zend_class_entry *parent;   //继承的父类
	int refcount;   //引用数
	zend_bool constants_updated;
	zend_uint ce_flags; 
    /**
        ZEND_ACC_IMPLICIT_ABSTRACT_CLASS: 类存在abstract方法
        ZEND_ACC_EXPLICIT_ABSTRACT_CLASS: 在类名称前加了abstract关键字
        ZEND_ACC_FINAL_CLASS
        ZEND_ACC_INTERFACE
     */

	HashTable function_table;           //方法
	HashTable default_properties;       //默认属性
	HashTable properties_info;          //属性信息
	HashTable default_static_members;   //类本身的静态变量
	HashTable *static_members;          
	HashTable constants_table;          //常量
	const struct _zend_function_entry *builtin_functions;   //方法定义入口

	union _zend_function *constructor;
	union _zend_function *destructor;
	union _zend_function *clone;

    /* 魔法方法 */
	union _zend_function *__get;
	union _zend_function *__set;
	union _zend_function *__unset;
	union _zend_function *__isset;
	union _zend_function *__call;
	union _zend_function *__callstatic;
	union _zend_function *__tostring;
	union _zend_function *serialize_func;
	union _zend_function *unserialize_func;

	zend_class_iterator_funcs iterator_funcs;       //迭代

	/* 类句柄 */
	zend_object_value (*create_object)(zend_class_entry *class_type TSRMLS_DC);
	zend_object_iterator *(*get_iterator)(zend_class_entry *ce, zval *object, 
        int by_ref TSRMLS_DC);
	
    /* a class implements this interface */
    int (*interface_gets_implemented)(zend_class_entry *iface,
        zend_class_entry *class_type TSRMLS_DC); 
	
    union _zend_function *(*get_static_method)(zend_class_entry *ce, 
        char* method, int method_len TSRMLS_DC);

	/* serializer callbacks */
	int (*serialize)(zval *object, unsigned char **buffer, zend_uint *buf_len, 
        zend_serialize_data *data TSRMLS_DC);
	int (*unserialize)(zval **object, zend_class_entry *ce, 
        const unsigned char *buf, zend_uint buf_len, 
        zend_unserialize_data *data TSRMLS_DC);

    //类实现的接口
	zend_class_entry **interfaces;
	//类实现接口数
    zend_uint num_interfaces;

	char *filename;         //类的存放文件的绝对地址
	zend_uint line_start;   //类定义的开始行
	zend_uint line_end;     //类定义的结束行
	char *doc_comment;
	zend_uint doc_comment_len;

    //类所在模块的入口
	struct _zend_module_entry *module;
};
```

词法|语法分析跳过不讲。

### 内置类
PHP有内置类，扩展也有类。因为类名不能重复，所以不能定义同类名的类了。
PHP还有特殊的类，包括：
- self
- parent
- static：多义
    - 修饰函数内变量用于定义静态局部变量
    - 修饰类成员函数和成员变量用于声明静态成员
    - 作用域修饰符`::`前表示静态延迟绑定的特殊类


在需要获取类名时会执行`zend_do_fetch_class`函数：
```
// @file: Zend/zend_compile.c

#define ZEND_FETCH_CLASS_DEFAULT	0
#define ZEND_FETCH_CLASS_SELF		1
#define ZEND_FETCH_CLASS_PARENT		2

void zend_do_fetch_class(znode *result, znode *class_name TSRMLS_DC) /* {{{ */
{
	long fetch_class_op_number;
	zend_op *opline;

	if (class_name->op_type == IS_CONST &&
	    Z_TYPE(class_name->u.constant) == IS_STRING &&
	    Z_STRLEN(class_name->u.constant) == 0) {
		/* Usage of namespace as class name not in namespace */
		zval_dtor(&class_name->u.constant);
		zend_error(E_COMPILE_ERROR, "Cannot use 'namespace' as a class name");
		return;
	}

	fetch_class_op_number = get_next_op_number(CG(active_op_array));
	opline = get_next_op(CG(active_op_array) TSRMLS_CC);

	opline->opcode = ZEND_FETCH_CLASS;
	SET_UNUSED(opline->op1);
	opline->extended_value = ZEND_FETCH_CLASS_GLOBAL;
	CG(catch_begin) = fetch_class_op_number;
	if (class_name->op_type == IS_CONST) {
		int fetch_type;

        //重点关注
		fetch_type = zend_get_class_fetch_type(
            class_name->u.constant.value.str.val, 
            class_name->u.constant.value.str.len);
		switch (fetch_type) {
			case ZEND_FETCH_CLASS_SELF:
			case ZEND_FETCH_CLASS_PARENT:
			case ZEND_FETCH_CLASS_STATIC:
				SET_UNUSED(opline->op2);
				opline->extended_value = fetch_type;
				zval_dtor(&class_name->u.constant);
				break;
			default:
				zend_resolve_class_name(class_name, &opline->extended_value, 
                    0 TSRMLS_CC);
				opline->op2 = *class_name;
				break;
		}
	} else {
		opline->op2 = *class_name;
	}
	opline->result.u.var = get_temporary_variable(CG(active_op_array));
	opline->result.u.EA.type = opline->extended_value;
    /* FIXME: Hack so that INIT_FCALL_BY_NAME still knows this is a class */
	opline->result.op_type = IS_VAR; 
	*result = opline->result;
}
/* }}} */


int zend_get_class_fetch_type(const char *class_name, uint class_name_len) /* {{{ */
{
	if ((class_name_len == sizeof("self")-1) &&
		!memcmp(class_name, "self", sizeof("self")-1)) {
		return ZEND_FETCH_CLASS_SELF;		
	} else if ((class_name_len == sizeof("parent")-1) &&
		!memcmp(class_name, "parent", sizeof("parent")-1)) {
		return ZEND_FETCH_CLASS_PARENT;
	} else if ((class_name_len == sizeof("static")-1) &&
		!memcmp(class_name, "static", sizeof("static")-1)) {
		return ZEND_FETCH_CLASS_STATIC;
	} else {
		return ZEND_FETCH_CLASS_DEFAULT;
	}
}
/* }}} */
```
