package module

import (
	"github.com/SurgicalSteel/elasthink/config"
	"github.com/SurgicalSteel/elasthink/redis"
)

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

var redisConfig *config.RedisConfigWrap

type ElasthinkSDK struct {
	redis *redis.Redis
}

func Initialize(redisMaxIdle int,
	redisMaxActive int,
	redisTimeout int,
	redisAddress string,
	isUsingStopWordsRemoval bool,
	stopWordRemovalData []string) ElasthinkSDK {

	spec := config.RedisConfigWrap{
		RedisElasthink: config.RedisConfig{
			MaxIdle:   redisMaxIdle,
			MaxActive: redisMaxActive,
			Address:   redisAddress,
			Timeout:   redisTimeout,
		},
	}

	newRedis := redis.InitRedis(spec)

	elasthinkSDK := ElasthinkSDK{
		Redis: newRedis,
	}
	return elasthinkSDK
}

func (es *ElasthinkSDK) createIndex(index string) (bool, error) {
	// ctx := context.Background()

	// err = json.Unmarshal(body, &requestPayload)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// es.Redis.
	return true, nil
}
