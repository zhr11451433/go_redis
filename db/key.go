package db

import (
	"errors"
	"math/rand"
	"time"
)

func (r *Redis) Del(key string) int {
	for i := range r.Key {
		if r.Key[i].KeyName != key {
			continue
		}
		r.Key = append(r.Key[:i], r.Key[i+1:]...)
		return 1
	}
	return 0
}

func (r *Redis) Exists(key string) int {
	_, ok := r.Find(key)
	if !ok {
		return 0
	}
	return 1
}

func (r *Redis) Expire(key string, duration time.Duration) int {
	redisKey, ok := r.Find(key)
	if !ok {
		return 0
	}
	if duration <= 0 {
		//不过期
		redisKey.KeyTime = 0
		return 1
	}
	redisKey.KeyTime = time.Now().Add(duration).Unix()
	return 1
}

func (r *Redis) findAt(key string, now int64) (*RedisKey, bool) {
	for i := range r.Key {
		if r.Key[i].KeyName == key &&
			(r.Key[i].KeyTime == 0 || r.Key[i].KeyTime > now) {
			return &r.Key[i], true
		}
	}
	return nil, false
}

func (r *Redis) Find(key string) (*RedisKey, bool) {
	return r.findAt(key, time.Now().Unix())
}

func (r *Redis) Persist(key string) int {
	redisKey, ok := r.Find(key)
	if !ok {
		return 0
	}
	if redisKey.KeyTime == 0 {
		return 0
	}
	redisKey.KeyTime = 0
	return 1
}

func (r *Redis) Ttl(key string) int64 {
	now := time.Now().Unix()

	entry, ok := r.findAt(key, now)
	if !ok {
		return -2
	}

	if entry.KeyTime == 0 {
		return -1
	}

	return entry.KeyTime - now
}

// RandOmKey 需要防止抽到过期key  蓄水池抽样
func (r *Redis) RandOmKey() (string, bool) {
	now := time.Now().Unix()
	//暂存最终随机选中的 key 名称
	var name string
	//目前遍历到的有效 key 数量
	count := 0

	for _, entry := range r.Key {
		if entry.KeyTime != 0 && entry.KeyTime <= now {
			continue
		}

		count++
		if rand.Intn(count) == 0 {
			name = entry.KeyName
		}
	}

	return name, count > 0
}

func (r *Redis) ReName(oldKey, newKey string) error {
	now := time.Now().Unix()

	// 1. 源键必须存在且有效
	old, oldOk := r.findAt(oldKey, now)
	if !oldOk {
		return errors.New("key not found")
	}

	// 2. 目标键有效时，重命名失败
	if _, newOk := r.findAt(newKey, now); newOk {
		return errors.New("key already exists")
	}

	// 3. 目标键不存在或已过期。
	// 如果切片里还留着目标键的过期记录，就将它删除。
	for i := len(r.Key) - 1; i >= 0; i-- {
		entry := &r.Key[i]

		if entry.KeyName == newKey &&
			entry.KeyTime != 0 &&
			entry.KeyTime <= now {
			r.Key = append(r.Key[:i], r.Key[i+1:]...)
		}
	}

	// 4. 删除记录可能改变切片，因此重新查找源键
	old, oldOk = r.findAt(oldKey, now)
	if !oldOk {
		return errors.New("key not found")
	}

	// 5. 修改源键名称
	old.KeyName = newKey
	return nil
}

func (r *Redis) Type(key string) (string, bool) {
	entry, ok := r.Find(key)
	if !ok {
		return "", false
	}
	return entry.KeyType, true
}
