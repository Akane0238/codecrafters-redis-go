package main

import (
	"sync"
)

// KVStore 为线程安全的键值存储数据库
type KVStore struct {
	mutex sync.RWMutex      // 读写锁
	data  map[string]string // 共享键值对存储
}

func NewKVStore() *KVStore {
	return &KVStore{
		data: make(map[string]string),
	}
}

// SET 命令
func (kv *KVStore) Set(key string, value string) {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	kv.data[key] = value
}

// GET 命令
func (kv *KVStore) Get(key string) (string, bool) {
	kv.mutex.RLock()
	defer kv.mutex.RUnlock()
	val, ok := kv.data[key]

	return val, ok
}
