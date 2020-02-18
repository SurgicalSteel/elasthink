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
	"github.com/SurgicalSteel/elasthink/config"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInitRedis(t *testing.T) {
	redisConfig := config.RedisConfigWrap{
		RedisElasthink: config.RedisConfig{
			Address:   "fake.redis.address",
			MaxActive: 30,
			MaxIdle:   10,
			Timeout:   10,
		},
	}

	expectedRedis := &Redis{
		Pool: &redigo.Pool{
			MaxIdle:     redisConfig.RedisElasthink.MaxIdle,
			MaxActive:   redisConfig.RedisElasthink.MaxActive,
			IdleTimeout: time.Duration(redisConfig.RedisElasthink.Timeout) * time.Second,
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial("tcp", redisConfig.RedisElasthink.Address)
			},
		},
	}

	actualRedis := InitRedis(redisConfig)

	assert.Equal(t, expectedRedis.Pool.MaxActive, actualRedis.Pool.MaxActive)
	assert.Equal(t, expectedRedis.Pool.MaxIdle, actualRedis.Pool.MaxIdle)
	assert.Equal(t, expectedRedis.Pool.IdleTimeout, actualRedis.Pool.IdleTimeout)
	assert.NotNil(t, actualRedis.Pool.Dial)

}

func TestSAdd(t *testing.T) {
	conn := redigomock.NewConn()
	redisMock := &Redis{
		Pool: redigo.NewPool(func() (redigo.Conn, error) {
			return conn, nil
		}, 10),
	}
	cmd := conn.Command("SADD", "campaign:ganteng", 666).Expect(int64(1))
	_, err := redisMock.SAdd("campaign:ganteng", []interface{}{666})
	if err != nil {
		t.Error("Expected : ok, but found error! err:", err.Error())
		return
	}
	if conn.Stats(cmd) != 1 {
		t.Error("Command SADD is not used!")
		return
	}
	conn.Clear()
}

func TestSMembers(t *testing.T) {
	conn := redigomock.NewConn()
	redisMock := &Redis{
		Pool: redigo.NewPool(func() (redigo.Conn, error) {
			return conn, nil
		}, 10),
	}
	cmd := conn.Command("SMEMBERS", "campaign:bangun").Expect([]interface{}{"123", "234", "345", "456"})
	_, err := redisMock.SMembers("campaign:bangun")
	if err != nil {
		t.Error("Expected : ok, but found error! err:", err.Error())
		return
	}
	if conn.Stats(cmd) != 1 {
		t.Error("Command SMEMBERS is not used!")
		return
	}
	conn.Clear()
}

func TestSRem(t *testing.T) {
	conn := redigomock.NewConn()
	redisMock := &Redis{
		Pool: redigo.NewPool(func() (redigo.Conn, error) {
			return conn, nil
		}, 10),
	}
	cmd := conn.Command("SREM", "campaign:ganteng", 666).Expect(int64(1))
	_, err := redisMock.SRem("campaign:ganteng", []interface{}{666})
	if err != nil {
		t.Error("Expected : ok, but found error! err:", err.Error())
		return
	}
	if conn.Stats(cmd) != 1 {
		t.Error("Command SREM is not used!")
		return
	}
	conn.Clear()
}

func TestKeysPrefix(t *testing.T) {
	conn := redigomock.NewConn()
	redisMock := &Redis{
		Pool: redigo.NewPool(func() (redigo.Conn, error) {
			return conn, nil
		}, 10),
	}
	//test case 1 : normal
	cmd := conn.Command("KEYS", "campaign:*").Expect([]interface{}{"campaign:bangun", "campaign:tidur", "campaign:ganteng", "campaign:jalan"})
	_, err := redisMock.KeysPrefix("campaign:")
	if err != nil {
		t.Error("Expected : ok, but found error! err:", err.Error())
		return
	}
	if conn.Stats(cmd) != 1 {
		t.Error("Command KEYS is not used!")
		return
	}

	//test case 2 : expect error
	keys, err := redisMock.KeysPrefix("             ")
	assert.Equal(t, errors.New("Prefix must be defined!"), err)
	assert.Equal(t, make([]string, 0), keys)
	conn.Clear()
}
