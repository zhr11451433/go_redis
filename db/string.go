package db

import (
	"errors"
	"time"
)

func (r *Redis) Set(key, value string) {
	r.GetSet(key, value)
}

// GetSet 有这个 Key：更新旧值。
// 没有这个 Key：创建新值
func (r *Redis) GetSet(key, value string) (old string, existed bool) {
	for i := range r.Key {
		entry := &r.Key[i]
		if entry.KeyName != key {
			continue
		}
		if entry.KeyTime != 0 && entry.KeyTime <= time.Now().Unix() {
			existed = false
		} else {
			old = entry.KeyValue
			existed = true
		}
		entry.KeyValue = value
		entry.KeyTime = 0
		entry.KeyType = "string"

		return old, existed

	}

	r.Key = append(r.Key, RedisKey{
		KeyName:  key,
		KeyValue: value,
		KeyType:  "string",
	})
	return "", false
}

func (r *Redis) Get(key string) (string, error) {
	entry, ok := r.Find(key)
	if !ok {
		return "", errors.New("key not found")
	}
	if entry.KeyType != "string" {
		return "", errors.New("key is not string")
	}
	return entry.KeyValue, nil
}

func (r *Redis) MGet(Keys ...string) []any {
	keys := make([]any, 0, len(Keys))
	for _, k := range Keys {
		entry, ok := r.Find(k)
		//key 不存在	返回 nil
		//key 存在但非 String	返回 nil（不报错）
		if !ok || entry.KeyType != "string" {
			keys = append(keys, nil)
			continue
		}
		keys = append(keys, entry.KeyValue)
	}
	return keys
}
