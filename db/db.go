package db

type RedisKey struct {
	KeyName  string `json:"key_name"`
	KeyValue string `json:"key_value"`
	KeyTime  int64  `json:"key_time"`
	KeyType  string `json:"key_type"`
}

type Redis struct {
	Key []RedisKey `json:"keys"`
}
