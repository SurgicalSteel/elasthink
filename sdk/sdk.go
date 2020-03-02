package module

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/SurgicalSteel/elasthink/config"
	"github.com/SurgicalSteel/elasthink/entity"
	"github.com/SurgicalSteel/elasthink/redis"
	"github.com/SurgicalSteel/elasthink/util"
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

// ElasthinkSDK is the main struct of elasthink SDK, initialized using initalize function
type ElasthinkSDK struct {
	redis                   *redis.Redis
	isUsingStopWordsRemoval bool
	stopWordRemovalData     []string
	entity                  EntityData
}

// CreateIndexSpec is the spec of createIndex function
// documentName is the name of the document
// documentID is the id of the document
type CreateIndexSpec struct {
	documentName string
	documentID   string
}

// Initialize is the function that return ElasthinkSDK
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
		redis:                   newRedis,
		isUsingStopWordsRemoval: isUsingStopWordsRemoval,
		stopWordRemovalData:     stopWordRemovalData,
	}
	return elasthinkSDK
}

func (es *ElasthinkSDK) CreateIndex(documentType string, documentID string, documentName string) (bool, error) {
	ctx := context.Background()
	// Validation
	err := validateCreateIndexRequestPayload(documentID, documentType, documentName)
	if err != nil {
		return false, err
	}

	// Tokenize document name set
	stopword := make(map[string]int)
	for _, k := range es.stopWordRemovalData {
		stopword[k] = 1
	}
	redis := es.redis
	documentNameSet := util.Tokenize(documentName, es.isUsingStopWordsRemoval, stopword)

	docType := getDocumentType(documentType, entity.Entity.GetDocumentTypes())
	errorExist := false
	errorKeys := ""

	for k, _ := range documentNameSet {
		key := fmt.Sprintf("%s:%s", docType, k)
		value := make([]interface{}, 1)
		value[0] = fmt.Sprintf("%d", documentID)
		_, err := redis.SAdd(key, value)
		if err != nil {
			errorExist = true
			errorKeys = errorKeys + " " + key + ","
			log.Println("[MODULE][CREATE INDEX] failed to add index on key :", key, "and document ID:", documentID)
			continue
		}
	}

	if errorExist {
		errorKeys = strings.TrimRight(errorKeys, ",")
		errorKeys = strings.TrimLeft(errorKeys, " ")
		return false, Response{
			StatusCode:   http.StatusInternalServerError,
			ErrorMessage: fmt.Sprintf("Error on adding following keys :%s", errorKeys),
			Data:         nil,
		}
	} else {
		return false, Response{
			StatusCode:   http.StatusOK,
			ErrorMessage: "",
			Data:         nil,
		}
	}

	return true, nil
}
