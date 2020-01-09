//Package redis is where we place all funcs related to redis operations
package redis

// Elasthink, An alternative to elasticsearch engine written in Go for small set of documents that uses inverted index to build the index and utilizes redis to store the indexes.
// Copyright (C) 2020 Yuwono Bangun Nagoro (a.k.a SurgicalSteel)
//
// This file is part of Elasthink
//
// Elasthink is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Elasthink is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
