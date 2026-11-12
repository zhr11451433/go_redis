package db

import (
	"errors"
	"strconv"
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

type Pair struct {
	Key   string
	Value string
}

func (r *Redis) MSet(P ...Pair) {
	for _, p := range P {
		key := p.Key
		value := p.Value
		r.Set(key, value)
	}
}

// SetNx 仅在 key 不存在时设置 key-value。若 key 存在则不执行任何操作
func (r *Redis) SetNx(key, value string) int {
	_, ok := r.Find(key)
	if ok {
		return 0
	}
	r.Set(key, value)
	return 1
}

func (r *Redis) StrLen(key string) int {
	entry, ok := r.Find(key)
	if !ok {
		return 0
	}
	//类型不对简化返回0
	if entry.KeyType != "string" {
		return 0
	}
	return len(entry.KeyValue)
}

func (r *Redis) IncrBy(key string, n int) (int, error) {
	entry, ok := r.Find(key)
	if !ok {
		r.Set(key, strconv.Itoa(n))
		return n, nil
	}
	m, err := strconv.Atoi(entry.KeyValue)
	if err != nil {
		return 0, err
	}
	r.Set(key, strconv.Itoa(n+m))
	return m + n, nil
}

func (r *Redis) Incr(key string) (int, error) {
	return r.IncrBy(key, 1)
}

func (r *Redis) DecrBY(key string, n int) (int, error) {
	entry, ok := r.Find(key)
	if !ok {
		r.Set(key, strconv.Itoa(-n))
		return -n, nil
	}
	m, err := strconv.Atoi(entry.KeyValue)
	if err != nil {
		return 0, err
	}
	r.Set(key, strconv.Itoa(m-n))
	return m - n, nil
}

func (r *Redis) Decr(key string) (int, error) {
	return r.DecrBY(key, 1)
}

func (r *Redis) Append(key, value string) (int, error) {
	entry, ok := r.Find(key)
	if !ok {
		r.Set(key, value)
		return len(value), nil
	}
	if entry.KeyType != "string" {
		return 0, errors.New("key is not string")
	}
	entry.KeyValue = entry.KeyValue + value
	return r.StrLen(key), nil
}
