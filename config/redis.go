package config

import (
	"fmt"
	"gopkg.in/gcfg.v1"
)

var redisConfig *RedisConfigWrap

//RedisConfigWrap is A wrapper for reading all redis connections (can be multiple connections).
type RedisConfigWrap struct {
	RedisElasthink RedisConfig
}

//RedisConfig is the basic configuration for a redis connection
type RedisConfig struct {
	Address   string
	MaxActive int
	MaxIdle   int
	Timeout   int
}

func readRedisConfig(path, env string) error {
	redisConfig = &RedisConfigWrap{}
	fileName := fmt.Sprintf("%s/redis/%s.ini", path, env)
	err := gcfg.ReadFileInto(redisConfig, fileName)
	return err
}

//GetRedisConfig gets the redis config that has been initializad
func GetRedisConfig() *RedisConfigWrap {
	return redisConfig
}
