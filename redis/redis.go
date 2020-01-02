//package redis is where we place all funcs related to redis operations
package redis

import (
	"errors"
	"fmt"
	"github.com/SurgicalSteel/elasthink/config"
	redigo "github.com/gomodule/redigo/redis"
	"strings"
	"sync"
	"time"
)

// Redis main struct
type Redis struct {
	Pool  *redigo.Pool
	mutex sync.Mutex
}

//NetworkTCP is the default network TCP
const NetworkTCP string = "tcp"

//InitRedis is a func to initialize redis that we are going to use
func InitRedis(redisConfig config.RedisConfigWrap) (*Redis, error) {
	newRedis := &Redis{
		Pool: &redigo.Pool{
			MaxIdle:     redisConfig.RedisElasthink.MaxIdle,
			MaxActive:   redisConfig.RedisElasthink.MaxActive,
			IdleTimeout: time.Duration(redisConfig.RedisElasthink.Timeout) * time.Second,
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial(NetworkTCP, redisConfig.RedisElasthink.Address)
			},
		},
	}
	return newRedis, nil
}

// SAdd add an item into a set
func (r *Redis) SAdd(key string, args []interface{}) (int64, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	return redigo.Int64(conn.Do("SADD", redigo.Args{key}.AddFlat(args)...))
}

// SMembers get members of a set
func (r *Redis) SMembers(key string) ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	return redigo.Strings(conn.Do("SMEMBERS", key))
}

// SRem remove an item from a set
func (r *Redis) SRem(keyRedis string, members []interface{}) (int64, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	return redigo.Int64(conn.Do("SREM", redigo.Args{keyRedis}.AddFlat(members)...))
}

// KeysPrefix get keys by a defined prefix
func (r *Redis) KeysPrefix(prefix string) ([]string, error) {
	prefix = strings.Trim(prefix, " ")
	if len(prefix) == 0 {
		return make([]string, 0), errors.New("Prefix must be defined!")
	}

	conn := r.Pool.Get()
	defer conn.Close()

	finalKeyPrefix := fmt.Sprintf("%s*", prefix)

	return redigo.Strings(conn.Do("KEYS", finalKeyPrefix))
}
