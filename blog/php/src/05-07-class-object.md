<!--
author: 刘青
date: 2017-04-04
title: 对象
type: note
source: https://github.com/reeze/tipi/blob/master/book/chapt05/05-07-class-object.markdown
tags: 
category: php/src
status: publish
summary: 
-->

对象是一个实体，它有状态，一般用变量表示，它有操作行为，用方法表示。

对象的抽象是类。

对象的通信是通过方法调用，以一种消息传递的方式进行。

### 对象的结构
我们看看对象的结构。
```
//@file: Zend/zend_types.h

typedef unsigned int zend_object_handle;
typedef struct _zend_object_handlers zend_object_handlers;

// 一个对象
typedef struct _zend_object_value {
	zend_object_handle handle; //在列表容器中的位置，得到 _zend_object
	zend_object_handlers *handlers;
} zend_object_value;

//@file: Zend/zend.h
// 一个对象
typedef struct _zend_object {
	zend_class_entry *ce;   // 对象的类
	HashTable *properties;  // 对象的属性
	HashTable *guards; /* protects from __get/__set ... recursion */
} zend_object;

//一个对象有的操作
struct _zend_object_handlers {
	/* general object functions */
	zend_object_add_ref_t					add_ref;
	zend_object_del_ref_t					del_ref;
	zend_object_clone_obj_t					clone_obj;
	/* individual object functions */
	zend_object_read_property_t				read_property;
	zend_object_write_property_t			write_property;
	zend_object_read_dimension_t			read_dimension;
	zend_object_write_dimension_t			write_dimension;
	zend_object_get_property_ptr_ptr_t		get_property_ptr_ptr;
	zend_object_get_t						get;
	zend_object_set_t						set;
	zend_object_has_property_t				has_property;
	zend_object_unset_property_t			unset_property;
	zend_object_has_dimension_t				has_dimension;
	zend_object_unset_dimension_t			unset_dimension;
	zend_object_get_properties_t			get_properties;
	zend_object_get_method_t				get_method;
	zend_object_call_method_t				call_method;
	zend_object_get_constructor_t			get_constructor;
	zend_object_get_class_entry_t			get_class_entry;
	zend_object_get_class_name_t			get_class_name;
	zend_object_compare_t					compare_objects;
	zend_object_cast_t						cast_object;
	zend_object_count_elements_t			count_elements;
	zend_object_get_debug_info_t			get_debug_info;
	zend_object_get_closure_t				get_closure;
};
```

### 对象的创建


### 对象池
> 对象池：PHP内核在运行中存储所有对象的列表（EG(objects_store)）。

```
//@file: Zend/zend_objects_API.h

//对象池中的一个槽位
typedef struct _zend_object_store_bucket {
	zend_bool destructor_called;
	zend_bool valid;
	union _store_bucket {
		struct _store_object {
			void *object;
			zend_objects_store_dtor_t dtor;
			zend_objects_free_object_storage_t free_storage;
			zend_objects_store_clone_t clone;
			const zend_object_handlers *handlers;
			zend_uint refcount;
			gc_root_buffer *buffered;
		} obj;
		struct {
			int next;
		} free_list;
	} bucket;
} zend_object_store_bucket;

//对象池
typedef struct _zend_objects_store {
	zend_object_store_bucket *object_buckets;
	zend_uint top;
	zend_uint size;
	int free_list_head;
} zend_objects_store;
```

PHP中有一套对象操作API维护对象池。

### 成员变量
对象的属性都在`properities`参数中。每个对象都有一套标准的操作函数。
获取成员变量，对象最后调用`read_property`，对应的标准函数是`zend_std_read_property`。
设置成员变量，对象最后调用`write_property`，对应的标准函数是`zend_std_write_property`。

### 成员函数
