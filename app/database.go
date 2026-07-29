package main

import (
	"sync"
	"time"
)

// KVStore 为线程安全的键值存储数据库
type KVStore struct {
	mutex sync.RWMutex    // 读写锁
	data  map[string]item // 共享键值对存储
}

// item 结构体存储值与Expire参数
type item struct {
	value      string
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
		expireTime: expire,
	}
}

// GET 命令
func (kv *KVStore) Get(key string) (string, bool) {
	kv.mutex.RLock()

	item, exist := kv.data[key]

	if !exist {
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
	return item.value, true
}
