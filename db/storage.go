package db

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Load(path string) (*Redis, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	var redis Redis
	err = json.NewDecoder(file).Decode(&redis)
	if err != nil {
		return nil, err
	}
	return &redis, nil
}

// 整体覆盖

func Save(path string, redis *Redis) error {
	data, err := json.MarshalIndent(redis, "", "  ")
	if err != nil {
		return err
	}
	return atomicWrite(path, data)
}

func atomicWrite(path string, data []byte) (err error) {
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		if err != nil {
			_ = tmp.Close()
			_ = os.Remove(tmpName)
		}
	}()

	if _, err = tmp.Write(data); err != nil {
		return err
	}
	if err = tmp.Sync(); err != nil {
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}
	if err = os.Chmod(tmpName, 0644); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

//清空

func FlushDB(path string, redis *Redis) error {
	redis.Key = []RedisKey{}
	if err := Save(path, redis); err != nil {
		return err
	}
	return nil
}
