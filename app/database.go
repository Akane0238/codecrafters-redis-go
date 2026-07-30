package main

import (
	"fmt"
	"sync"
	"time"
)

type ValueType int

const (
	StringType ValueType = iota
	ListType
)

// KVStore 为线程安全的键值存储数据库
type KVStore struct {
	mutex sync.RWMutex    // 读写锁
	data  map[string]item // 共享键值对存储
}

// item 结构体存储值与Expire参数
type item struct {
	value      any
	valueType  ValueType
	expireTime time.Time // 过期的时间戳，零值表示不过期
}

// --- Methods ---
func NewKVStore() *KVStore {
	return &KVStore{
		data: make(map[string]item),
	}
}

// SET 命令
func (kv *KVStore) Set(key string, value string, ttl time.Duration) {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()

	var expire time.Time
	if ttl > 0 {
		expire = time.Now().Add(ttl)
	}
	kv.data[key] = item{
		value:      value,
		valueType:  StringType,
		expireTime: expire,
	}
}

// GET 命令
func (kv *KVStore) Get(key string) (string, bool) {
	kv.mutex.RLock()

	item, exist := kv.data[key]

	if !exist || item.valueType != StringType {
		kv.mutex.RUnlock()
		return "", false
	}

	// 数据已过期，需要删除
	if !item.expireTime.IsZero() && time.Now().After(item.expireTime) {
		// 先释放读锁，再尝试获得写锁删除
		kv.mutex.RUnlock()
		kv.mutex.Lock()
		defer kv.mutex.Unlock()

		// 释放读锁的间隙 value 可能已被其他协程删除或更新
		item, exist := kv.data[key]
		if !exist {
			return "", false
		}

		// 二次确认后执行安全删除
		if !item.expireTime.IsZero() && time.Now().After(item.expireTime) {
			delete(kv.data, key)
		}
		return "", false
	}

	kv.mutex.RUnlock()
	return item.value.(string), true
}

// RPUSH 命令
func (kv *KVStore) RPush(list_key string, elements []string) int {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	val, ok := kv.data[list_key]
	if !ok {
		// 创建List
		list := make([]string, 0)
		list = append(list, elements...)
		kv.data[list_key] = item{
			value:     list,
			valueType: ListType,
		}
		return len(list)
	}

	if val.valueType != ListType {
		fmt.Println("Unalble to push data to not 'ListType' container.")
		return -1
	}

	// fmt.Println("Before: ", val.value.([]string)) // debug
	val.value = append(val.value.([]string), elements...)
	kv.data[list_key] = val
	// fmt.Println("After: ", val.value.([]string)) // debug

	return len(val.value.([]string))
}

// LRANGE 命令
func (kv *KVStore) LRange(list_key string, start int, stop int) []string {
	kv.mutex.RLock()
	defer kv.mutex.RUnlock()

	item, ok := kv.data[list_key]

	if !ok || item.valueType != ListType {
		return nil
	}

	size := len(kv.data[list_key].value.([]string))
	if start >= size || start > stop {
		return nil
	}
	if stop >= size {
		stop = size - 1
	}

	fmt.Printf("start:%d  stop:%d\n", start, stop)     // debug
	fmt.Println(item.value.([]string)[start : stop+1]) // debug

	return item.value.([]string)[start : stop+1]
}
