<!--
author: 刘青
date: 2017-03-25
title: PHP变量结构
type: note
source: https://github.com/reeze/tipi/blob/master/book/chapt03/03-01-00-variables-structure.markdown
tags: 
category: php/src
status: publish 
summary: 
-->

### 变量的结构和类型
PHP是一门弱类型语言，也有类型，分为三类：
- 标量类型：boolean, integer, float(double), string
- 复合类型：array, object
- 特殊类型：resource, NULL

PHP的弱类型的具体实现为：所有的变量使用 `zval` 数据结构保存。其包含变量的值和类型。

```
// file: Zend/zend.h
struct _zval_struct {           /* Variable information */
	zvalue_value value;		    // 变量的值
	zend_uint refcount__gc;     // 引用计数
	zend_uchar type;	        // 变量的类型:IS_NULL,IS_BOOL,IS_LONG,IS_DOUBLE,IS_STRING,IS_ARRAY,IS_OBJECT,IS_RESOURCE
	zend_uchar is_ref__gc;      // 是否为引用
};

// file: Zend/zend.h
typedef union _zvalue_value {
	long lval;					/* long value */
	double dval;				/* double value */
	struct {
		char *val;
		int len;
	} str;
	HashTable *ht;				/* hash table value */
	zend_object_value obj;
} zvalue_value;
// 是个union

// Zend/zend.types.h
typedef struct _zend_object_value {
	zend_object_handle handle;      // EG(objects_store).object_buckets的索引
	zend_object_handlers *handlers;
} zend_object_value;

typedef unsigned int zend_object_handle;
typedef struct _zend_object_handlers zend_object_handlers;
```

我们先看看HashTable的实现。

### HashTable
PHP中使用得最频繁的数据类型为字符串和数组。
哈希表通常提供Search, inSERT, Delete等操作。这些操作最坏的性能是O(n)，当良好的设计可以达到O(1)。

> 哈希表：通过哈希函数将特定的键映射到特定值的一种数据结构。

相关概念有：
- key：用于操作数据的标识
- solt/bucket：一个数据单元，值真正存放的地方
- hash function：将key映射到slot的函数。h(key) -> index
- hash collision：hash 函数将两个不同的key映射到同一个solt的情况

hash collsion的解决方案
- 链接法：通过一个链表来保持slot值，当冲突发生时使用链表来保存这些值。最坏情况下哈希表会退化为一个链表。（PHP采用）
- 开放寻址法：当冲突发生时持续寻找下一个slat直到成功。这如有让后续的key发生冲突的概率增高。

我们看看源码：
```
//@file: Zend/zend.h

/*************************
 * 数据结构
 *************************/

//哈希表，存储整个哈希表需要的基本信息
typedef struct _hashtable {
	uint nTableSize;            // hash的大小 8 * n^2
	uint nTableMask;            // nTableSize-1，索引取值的优化
	uint nNumOfElements;        // hash中元素的个数
	ulong nNextFreeElement;     // hash中下一个元素
	Bucket *pInternalPointer;	// 当前遍历的指针（foreache比for快的原因之一）
	Bucket *pListHead;          // 头元素指针
	Bucket *pListTail;          // 尾元素指针
	Bucket **arBuckets;         // 链表
	dtor_func_t pDestructor;    // 在删除时执行的回调函数，用于资源释放
	zend_bool persistent;       // bucket的内存使用操作系统还是PHP内部函数分配
	unsigned char nApplyCount;  // 被递归访问的次数
	zend_bool bApplyProtection; // 是否允许多次访问，不允许的话最多递归3次
#if ZEND_DEBUG
	int inconsistent;
#endif
} HashTable;

// 数据单元（槽位）
typedef struct bucket {
	ulong h;					// 对char *key 进行hash后的值，或者时用户指定的数字索引的值
	uint nKeyLength;            // hash关键字的长度，如果数组索引为数字，为0
	void *pData;                // 指向value，一般是用户数据的副本，若是指针数据则指向 pDataPtr
	void *pDataPtr;             // 如果是指针数据，此值会指向真正的数据
	struct bucket *pListNext;   // hash表的下一个元素
	struct bucket *pListLast;
	struct bucket *pNext;       // bucket的下一个元素
	struct bucket *pLast;
	char arKey[1];              // 当前值对应的key字符串。这个字段只能定义在最后，实现变长结构体。
} Bucket;


/*************************
 * 哈希函数
 *************************/
/*
 * DJBX33A (Daniel J. Bernstein, Times 33 with Addition)
 *
 * This is Daniel J. Bernstein's popular `times 33' hash function as
 * posted by him years ago on comp.lang.c. It basically uses a function
 * like ``hash(i) = hash(i-1) * 33 + str[i]''. This is one of the best
 * known hash functions for strings. Because it is both computed very
 * fast and distributes very well.
 *
 * The magic of number 33, i.e. why it works better than many other
 * constants, prime or not, has never been adequately explained by
 * anyone. So I try an explanation: if one experimentally tests all
 * multipliers between 1 and 256 (as RSE did now) one detects that even
 * numbers GGare not useable at all. The remaining 128 odd numbers
 * (except for the number 1) work more or less all equally well. They
 * all distribute in an acceptable way and this way fill a hash table
 * with an average percent of approx. 86%. 
 *
 * If one compares the Chi^2 values of the variants, the number 33 not
 * even has the best value. But the number 33 and a few other equally
 * good numbers like 17, 31, 63, 127 and 129 have nevertheless a great
 * advantage to the remaining numbers in the large set of possible
 * multipliers: their multiply operation can be replaced by a faster
 * operation based on just one shift plus either a single addition
 * or subtraction operation. And because a hash function has to both
 * distribute good _and_ has to be very fast to compute, those few
 * numbers should be preferred and seems to be the reason why Daniel J.
 * Bernstein also preferred it.
 *
 *
 *                  -- Ralf S. Engelschall <rse@engelschall.com>
 */

static inline ulong zend_inline_hash_func(const char *arKey, uint nKeyLength)
{
	register ulong hash = 5381;

	/* variant with the hash unrolled eight times */
	for (; nKeyLength >= 8; nKeyLength -= 8) {
		hash = ((hash << 5) + hash) + *arKey++;
		hash = ((hash << 5) + hash) + *arKey++;
		hash = ((hash << 5) + hash) + *arKey++;
		hash = ((hash << 5) + hash) + *arKey++;
		hash = ((hash << 5) + hash) + *arKey++;
		hash = ((hash << 5) + hash) + *arKey++;
		hash = ((hash << 5) + hash) + *arKey++;
		hash = ((hash << 5) + hash) + *arKey++;
	}
	switch (nKeyLength) {
		case 7: hash = ((hash << 5) + hash) + *arKey++; /* fallthrough... */
		case 6: hash = ((hash << 5) + hash) + *arKey++; /* fallthrough... */
		case 5: hash = ((hash << 5) + hash) + *arKey++; /* fallthrough... */
		case 4: hash = ((hash << 5) + hash) + *arKey++; /* fallthrough... */
		case 3: hash = ((hash << 5) + hash) + *arKey++; /* fallthrough... */
		case 2: hash = ((hash << 5) + hash) + *arKey++; /* fallthrough... */
		case 1: hash = ((hash << 5) + hash) + *arKey++; break;
		case 0: break;
EMPTY_SWITCH_DEFAULT_CASE()
	}
	return hash;
}


/*************************
 * 操作：
 *  初始化
 *  查找，插入，删除，更新
 *  迭代和循环
 *  复制，排序，倒置，销毁
 *************************/

// 哈希表初始化
ZEND_API int _zend_hash_init(
    HashTable *ht, 
    uint nSize, 
    hash_func_t pHashFunction, 
    dtor_func_t pDestructor, 
    zend_bool persistent ZEND_FILE_LINE_DC
) {
	uint i = 3;
	Bucket **tmp;

	SET_INCONSISTENT(HT_OK);

	if (nSize >= 0x80000000) {// 32位最大的值
		/* prevent overflow */
		ht->nTableSize = 0x80000000;
	} else {
		while ((1U << i) < nSize) {// 最小8，2的倍数
			i++;
		}
		ht->nTableSize = 1 << i;
	}

	ht->nTableMask = ht->nTableSize - 1
	ht->pDestructor = pDestructor;
	ht->arBuckets = NULL;
	ht->pListHead = NULL;
	ht->pListTail = NULL;
	ht->nNumOfElements = 0;
	ht->nNextFreeElement = 0;
	ht->pInternalPointer = NULL;
	ht->persistent = persistent;
	ht->nApplyCount = 0;
	ht->bApplyProtection = 1;
	
	/* Uses ecalloc() so that Bucket* == NULL */
	if (persistent) { //持久化数据能够在多个请求中访问，非持久化数据在请求结束时就会释放
		tmp = (Bucket **) calloc(ht->nTableSize, sizeof(Bucket *));
		if (!tmp) {
			return FAILURE;
		}
		ht->arBuckets = tmp;
	} else {
		tmp = (Bucket **) ecalloc_rel(ht->nTableSize, sizeof(Bucket *));
		if (tmp) {
			ht->arBuckets = tmp;
		}
	}
	
	return SUCCESS;
}


// 更新
#define zend_hash_update(ht, arKey, nKeyLength, pData, nDataSize, pDest) \
		_zend_hash_add_or_update(ht, arKey, nKeyLength, pData, nDataSize, pDest, HASH_UPDATE ZEND_FILE_LINE_CC)

// 添加
#define zend_hash_add(ht, arKey, nKeyLength, pData, nDataSize, pDest) \
		_zend_hash_add_or_update(ht, arKey, nKeyLength, pData, nDataSize, pDest, HASH_ADD ZEND_FILE_LINE_CC)

// 底层实现
ZEND_API int _zend_hash_add_or_update(
    HashTable *ht, 
    const char *arKey, 
    uint nKeyLength, 
    void *pData, 
    uint nDataSize, 
    void **pDest, 
    int flag ZEND_FILE_LINE_DC
){
	ulong h;
	uint nIndex;
	Bucket *p;

	IS_CONSISTENT(ht);

	if (nKeyLength <= 0) {
#if ZEND_DEBUG
		ZEND_PUTS("zend_hash_update: Can't put in empty key\n");
#endif
		return FAILURE;
	}

	h = zend_inline_hash_func(arKey, nKeyLength);   // 得到hash值
	nIndex = h & ht->nTableMask;                    // 得到索引

	p = ht->arBuckets[nIndex];                      // 得到槽点内容
	while (p != NULL) { // 遍历槽点列表
		if ((p->h == h) && (p->nKeyLength == nKeyLength)) {     // 找到存在的数据
			if (!memcmp(p->arKey, arKey, nKeyLength)) {
				if (flag & HASH_ADD) {
					return FAILURE;
				}
				HANDLE_BLOCK_INTERRUPTIONS();
#if ZEND_DEBUG
				if (p->pData == pData) {
					ZEND_PUTS("Fatal error in zend_hash_update: p->pData == pData\n");
					HANDLE_UNBLOCK_INTERRUPTIONS();
					return FAILURE;
				}
#endif
				if (ht->pDestructor) {
					ht->pDestructor(p->pData);
				}
				UPDATE_DATA(ht, p, pData, nDataSize);
				if (pDest) {
					*pDest = p->pData;
				}
				HANDLE_UNBLOCK_INTERRUPTIONS();
				return SUCCESS;
			}
		}
		p = p->pNext;
	}
	
    // 创建一个新的槽点
	p = (Bucket *) pemalloc(sizeof(Bucket) - 1 + nKeyLength, ht->persistent);
	if (!p) {
		return FAILURE;
	}
	memcpy(p->arKey, arKey, nKeyLength);
	p->nKeyLength = nKeyLength;
	INIT_DATA(ht, p, pData, nDataSize);
	p->h = h;
	CONNECT_TO_BUCKET_DLLIST(p, ht->arBuckets[nIndex]);     // 加入到链表后面
	if (pDest) {
		*pDest = p->pData;
	}

	HANDLE_BLOCK_INTERRUPTIONS();
	CONNECT_TO_GLOBAL_DLLIST(p, ht);
	ht->arBuckets[nIndex] = p;
	HANDLE_UNBLOCK_INTERRUPTIONS();

	ht->nNumOfElements++;
	ZEND_HASH_IF_FULL_DO_RESIZE(ht);		/* If the Hash table is full, resize it */
	return SUCCESS;
}
```

哈希表的实现大概就是这样，哈希表和双向链表的混合实现。
