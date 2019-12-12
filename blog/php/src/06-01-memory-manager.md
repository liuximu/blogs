<!--
author: 刘青
date: 2017-04-07
title: 内存管理
type: note
source: https://github.com/reeze/tipi/blob/master/book/chapt06/06-01-memory-management-overview.markdown
tags: 
category: php/src
status: drash
summary: 
-->

计算器开机后操作系统就进入内存了，它管理计算器的资源，提供内存管理单元用于处理CPU对内存的访问。
其他的应用程序需要使用内存时需要向操作系统申请，它无法直接对内存进行访问。

由于想操作系统事情内存空间会引发系统调用。系统调用会将CUP从用户态切换到内核，切换代价巨大。对于需要频繁使用和释放内存的应用（Web服务器等）会自己在用户态进行内存管理，比如一次性申请较大块内存，使用时也不及时归还操作系统。

PHP不需要显式的对内存进行管理，这些工作由Zend引擎管理。PHP内部有自己的内存管理体系，会自动将不再使用的内存垃圾进行释放。 
### 内存配置
在php.ini中有如下配置项：
```
memory_limit = 32M
```

如果PHP环境没有禁用ini_set()，也可以动态配置：
```
ini_set("memory_limit", "128M");
```

可以通过 `memory_get_usage()` 和 `memory_get_peak_usage()` 查看内存使用情况。

前面了解到的引用计数，函数表，符号表，常量表等，我们都刻意避免不必要的内存浪费。

### 内存管理
内存管理包括：
- 申请内存
- 释放内存

在PHP的内存管理也包含这样的内容，Zend内核中以宏的形式作为借口提供给外部使用。

PHP的内存管理可以看成是分层的：
- 存储层（storage）：面向底层。通过`malloc()`, `mmap()` 等函数向系统真正的申请内存，通过`free()`释放内存。通常申请大块内存。不同平台有不同的内存管理方案。
- 堆层（heap）：控制整个PHP内存管理的过程。
- 接口层（emalloc/efree）：面向应用层。

![](https://raw.githubusercontent.com/reeze/tipi/master/book/images/chapt06/06-02-01-zend-memeory-manager.jpg)

**接口层**
```
//@file: Zend/zend_alloc.h

/* Standard wrapper macros */
#define emalloc(size)						_emalloc((size) ZEND_FILE_LINE_CC ZEND_FILE_LINE_EMPTY_CC)
#define safe_emalloc(nmemb, size, offset)	_safe_emalloc((nmemb), (size), (offset) ZEND_FILE_LINE_CC ZEND_FILE_LINE_EMPTY_CC)
#define efree(ptr)							_efree((ptr) ZEND_FILE_LINE_CC ZEND_FILE_LINE_EMPTY_CC)
#define ecalloc(nmemb, size)				_ecalloc((nmemb), (size) ZEND_FILE_LINE_CC ZEND_FILE_LINE_EMPTY_CC)
#define erealloc(ptr, size)					_erealloc((ptr), (size), 0 ZEND_FILE_LINE_CC ZEND_FILE_LINE_EMPTY_CC)
#define safe_erealloc(ptr, nmemb, size, offset)	_safe_erealloc((ptr), (nmemb), (size), (offset) ZEND_FILE_LINE_CC ZEND_FILE_LINE_EMPTY_CC)
#define erealloc_recoverable(ptr, size)		_erealloc((ptr), (size), 1 ZEND_FILE_LINE_CC ZEND_FILE_LINE_EMPTY_CC)
#define estrdup(s)							_estrdup((s) ZEND_FILE_LINE_CC ZEND_FILE_LINE_EMPTY_CC)
#define estrndup(s, length)					_estrndup((s), (length) ZEND_FILE_LINE_CC ZEND_FILE_LINE_EMPTY_CC)
#define zend_mem_block_size(ptr)			_zend_mem_block_size((ptr) TSRMLS_CC ZEND_FILE_LINE_CC ZEND_FILE_LINE_EMPTY_CC)
```
这些都是一些宏，官方文档建议使用宏来实现解耦。

**堆层**
```
//@file: Zend/zend_alloc.c

/* mm block type */
typedef struct _zend_mm_block_info { 
	size_t _size;   //block的大小
	size_t _prev;   //前一个block的大小
} zend_mm_block_info; 

typedef struct _zend_mm_block {
	zend_mm_block_info info;
} zend_mm_block;

typedef struct _zend_mm_small_free_block { //双向链表
	zend_mm_block_info info;
	struct _zend_mm_free_block *prev_free_block;    //前一个块
	struct _zend_mm_free_block *next_free_block;    //后一个块
} zend_mm_small_free_block;

typedef struct _zend_mm_free_block { //双向链表 + 树结构
	zend_mm_block_info info;
	struct _zend_mm_free_block *prev_free_block;
	struct _zend_mm_free_block *next_free_block;

	struct _zend_mm_free_block **parent;    //父节点
	struct _zend_mm_free_block *child[2];   //两个子节点
} zend_mm_free_block;

struct _zend_mm_heap {          // 堆
	int                 use_zend_alloc;                 //是否使用zend内存管理器
	void               *(*_malloc)(size_t);             //内存分配函数
	void                (*_free)(void*);                //内存释放函数
	void               *(*_realloc)(void*, size_t);     
	size_t              free_bitmap;                    //小块空闲内存标识
	size_t              large_free_bitmap;              //大块空闲内存标识
	size_t              block_size;                     //一次内存分配的段大小
	size_t              compact_size;                   //压缩操作边界值
	zend_mm_segment    *segments_list;                  //段指针列表
	zend_mm_storage    *storage;                        //所调用的存储层
	size_t              real_size;                      //堆真实大小
	size_t              real_peak;                      //堆真实大小的峰值
	size_t              limit;                          //堆的内存边界
	size_t              size;                           //堆大小
	size_t              peak;                           //对大小的峰值
	size_t              reserve_size;                   //备用堆大小
	void               *reserve;                        //备用堆
	int                 overflow;                       //内存溢出数
	int                 internal;               
#if ZEND_MM_CACHE
	unsigned int        cached;                         //已缓存大小
	zend_mm_free_block *cache[ZEND_MM_NUM_BUCKETS];     //缓存数组
#endif

    //PHP中的内存管理的主要工作就是维护下面三个列表
	zend_mm_free_block *free_buckets[ZEND_MM_NUM_BUCKETS*2];        //小块内存数组
	zend_mm_free_block *large_free_buckets[ZEND_MM_NUM_BUCKETS];    //大块内存数组
	zend_mm_free_block *rest_buckets[2];                            //剩余内存数组
};
```

看看堆层的初始化：
```
//@file: Zend/zend_alloc.c

ZEND_API zend_mm_heap *zend_mm_startup(void)
{
	int i;
	size_t seg_size;
    //从配置中找到内存分配方案
	char *mem_type = getenv("ZEND_MM_MEM_TYPE");
	char *tmp;
	const zend_mm_mem_handlers *handlers;
	zend_mm_heap *heap;

	if (mem_type == NULL) {
		i = 0;
	} else {
		for (i = 0; mem_handlers[i].name; i++) {
			if (strcmp(mem_handlers[i].name, mem_type) == 0) {
				break;
			}
		}
		if (!mem_handlers[i].name) {
			fprintf(stderr, "Wrong or unsupported zend_mm storage type '%s'\n", mem_type);
			fprintf(stderr, "  supported types:\n");
/* See http://support.microsoft.com/kb/190351 */
#ifdef PHP_WIN32
			fflush(stderr);
#endif
			for (i = 0; mem_handlers[i].name; i++) {
				fprintf(stderr, "    '%s'\n", mem_handlers[i].name);
			}
/* See http://support.microsoft.com/kb/190351 */
#ifdef PHP_WIN32
			fflush(stderr);
#endif
			exit(255);
		}
	}
	handlers = &mem_handlers[i];

    //从配置中找到段的大小
	tmp = getenv("ZEND_MM_SEG_SIZE");
	if (tmp) {
		seg_size = zend_atoi(tmp, 0);
		if (zend_mm_low_bit(seg_size) != zend_mm_high_bit(seg_size)) {
			fprintf(stderr, "ZEND_MM_SEG_SIZE must be a power of two\n");
/* See http://support.microsoft.com/kb/190351 */
#ifdef PHP_WIN32
			fflush(stderr);
#endif
			exit(255);
		} else if (seg_size < ZEND_MM_ALIGNED_SEGMENT_SIZE + ZEND_MM_ALIGNED_HEADER_SIZE) {
			fprintf(stderr0, "ZEND_MM_SEG_SIZE is too small\n");
/* See http://support.microsoft.com/kb/190351 */
#ifdef PHP_WIN32
			fflush(stderr);
#endif
			exit(255);
		}
	} else {
		seg_size = ZEND_MM_SEG_SIZE;
	}

    //得到堆
	heap = zend_mm_startup_ex(handlers, seg_size, ZEND_MM_RESERVE_SIZE, 0, NULL);
	if (heap) {
		tmp = getenv("ZEND_MM_COMPACT");
		if (tmp) {
			heap->compact_size = zend_atoi(tmp, 0);
		} else {
			heap->compact_size = 2 * 1024 * 1024;
		}
	}
	return heap;
}

```

我们重点看`zend_mm_startup_ex`:
```
/* Notes:
 * - This function may alter the block_sizes values to match platform alignment
 * - This function does *not* perform sanity checks on the arguments
 */
ZEND_API zend_mm_heap *zend_mm_startup_ex(const zend_mm_mem_handlers *handlers, size_t block_size, size_t reserve_size, int internal, void *params)
{
	zend_mm_storage *storage;
	zend_mm_heap    *heap;

#if ZEND_MM_HEAP_PROTECTION
	if (_mem_block_start_magic == 0) {
		zend_mm_random((unsigned char*)&_mem_block_start_magic, sizeof(_mem_block_start_magic));
	}
	if (_mem_block_end_magic == 0) {
		zend_mm_random((unsigned char*)&_mem_block_end_magic, sizeof(_mem_block_end_magic));
	}
#endif
#if ZEND_MM_COOKIES
	if (_zend_mm_cookie == 0) {
		zend_mm_random((unsigned char*)&_zend_mm_cookie, sizeof(_zend_mm_cookie));
	}
#endif

    //2^n意味着1打头，其他都是0，高位的0和低位的0在同一个位置
	if (zend_mm_low_bit(block_size) != zend_mm_high_bit(block_size)) {
		fprintf(stderr, "'block_size' must be a power of two\n");
/* See http://support.microsoft.com/kb/190351 */
#ifdef PHP_WIN32
		fflush(stderr);
#endif
		exit(255);
	}
	storage = handlers->init(params);
	if (!storage) {
		fprintf(stderr, "Cannot initialize zend_mm storage [%s]\n", handlers->name);
/* See http://support.microsoft.com/kb/190351 */
#ifdef PHP_WIN32
		fflush(stderr);
#endif
		exit(255);
	}
	storage->handlers = handlers;

	heap = malloc(sizeof(struct _zend_mm_heap));
	if (heap == NULL) {
		fprintf(stderr, "Cannot allocate heap for zend_mm storage [%s]\n", handlers->name);
#ifdef PHP_WIN32
		fflush(stderr);
#endif
		exit(255);
	}
	heap->storage = storage;
	heap->block_size = block_size;
	heap->compact_size = 0;
	heap->segments_list = NULL;
	zend_mm_init(heap);
# if ZEND_MM_CACHE_STAT
	memset(heap->cache_stat, 0, sizeof(heap->cache_stat));
# endif

	heap->use_zend_alloc = 1;
	heap->real_size = 0;
	heap->overflow = 0;
	heap->real_peak = 0;
	heap->limit = ZEND_MM_LONG_CONST(1)<<(ZEND_MM_NUM_BUCKETS-2);
	heap->size = 0;
	heap->peak = 0;
	heap->internal = internal;
	heap->reserve = NULL;
	heap->reserve_size = reserve_size;
	if (reserve_size > 0) {
		heap->reserve = _zend_mm_alloc_int(heap, reserve_size ZEND_FILE_LINE_CC ZEND_FILE_LINE_EMPTY_CC);
	}
	if (internal) {
		int i;
		zend_mm_free_block *p, *q, *orig;
		zend_mm_heap *mm_heap = _zend_mm_alloc_int(heap, sizeof(zend_mm_heap)  ZEND_FILE_LINE_CC ZEND_FILE_LINE_EMPTY_CC);

		*mm_heap = *heap;

		p = ZEND_MM_SMALL_FREE_BUCKET(mm_heap, 0);
		orig = ZEND_MM_SMALL_FREE_BUCKET(heap, 0);
		for (i = 0; i < ZEND_MM_NUM_BUCKETS; i++) {
			q = p;
			while (q->prev_free_block != orig) {
				q = q->prev_free_block;
			}
			q->prev_free_block = p;
			q = p;
			while (q->next_free_block != orig) {
				q = q->next_free_block;
			}
			q->next_free_block = p;
			p = (zend_mm_free_block*)((char*)p + sizeof(zend_mm_free_block*) * 2);
			orig = (zend_mm_free_block*)((char*)orig + sizeof(zend_mm_free_block*) * 2);
			if (mm_heap->large_free_buckets[i]) {
				mm_heap->large_free_buckets[i]->parent = &mm_heap->large_free_buckets[i];
			}
		}
		mm_heap->rest_buckets[0] 
            = mm_heap->rest_buckets[1] 
            = ZEND_MM_REST_BUCKET(mm_heap);

		free(heap);
		heap = mm_heap;
	}
	return heap;
}

#ifdef _WIN64
# define ZEND_MM_LONG_CONST(x)	(x##i64)
#else
# define ZEND_MM_LONG_CONST(x)	(x##L)
#endif

#define ZEND_MM_NUM_BUCKETS (sizeof(size_t) << 3)

#define ZEND_MM_NUM_BUCKETS (sizeof(size_t) << 3)

#define ZEND_MM_SMALL_FREE_BUCKET(heap, index) \
	(zend_mm_free_block*) ((char*)&heap->free_buckets[index * 2] + \
		sizeof(zend_mm_free_block*) * 2 - \
		sizeof(zend_mm_small_free_block))

#define ZEND_MM_REST_BUCKET(heap) \
	(zend_mm_free_block*)((char*)&heap->rest_buckets[0] + \
		sizeof(zend_mm_free_block*) * 2 - \
		sizeof(zend_mm_small_free_block))
```
