package module

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/SurgicalSteel/elasthink/config"
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

/// Variables and structures
var redisConfig *config.RedisConfigWrap

// ElasthinkSDK is the main struct of elasthink SDK, initialized using initalize function
type ElasthinkSDK struct {
	redis                   *redis.Redis
	isUsingStopWordsRemoval bool
	stopWordRemovalData     []string
	availableDocumentType   map[string]int
}

// InitializeSpec is the payload for initialize Elasthink SDK
type InitializeSpec struct {
	redisConfig config.RedisConfig
	sdkConfig   SdkConfig
}

// SdkConfig is the configuration for initialize Elasthink SDK
type SdkConfig struct {
	isUsingStopWordsRemoval bool
	stopWordRemovalData     []string
	availableDocumentType   []string
}

// CreateIndexSpec is the spec of createIndex function
// documentName is the name of the document
// documentID is the id of the document
type CreateIndexSpec struct {
	documentName string
	documentID   string
}

/// Functions

// Initialize is the function that return ElasthinkSDK
func Initialize(initializeSpec InitializeSpec) ElasthinkSDK {

	spec := config.RedisConfigWrap{
		RedisElasthink: config.RedisConfig{
			MaxIdle:   initializeSpec.redisConfig.MaxIdle,
			MaxActive: initializeSpec.redisConfig.MaxActive,
			Address:   initializeSpec.redisConfig.Address,
			Timeout:   initializeSpec.redisConfig.Timeout,
		},
	}
	newRedis := redis.InitRedis(spec)

	availableDocumentType := make(map[string]int)
	for _, doctype := range initializeSpec.sdkConfig.availableDocumentType {
		availableDocumentType[doctype] = 1
	}

	elasthinkSDK := ElasthinkSDK{
		redis:                   newRedis,
		isUsingStopWordsRemoval: initializeSpec.sdkConfig.isUsingStopWordsRemoval,
		stopWordRemovalData:     initializeSpec.sdkConfig.stopWordRemovalData,
		availableDocumentType:   availableDocumentType,
	}
	return elasthinkSDK
}

// CreateIndex is a function to create new index based on documentType, documentID, and document name
// documentType is the type of the document, to categorize documents. For example: campaign
// documentID, is the ID of document, the key of document. For example: 1
// documentName, is the name of documennt, the value which will be indexed. For example: "we want to eat seafood on a restaurant"
func (es *ElasthinkSDK) CreateIndex(documentType string, documentID int64, documentName string) (bool, error) {
	// Validation
	err := es.validateCreateIndexRequestPayload(documentID, documentType, documentName)
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

	docType := documentType
	errorExist := false
	errorKeys := ""

	for k := range documentNameSet {
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
		return false, err
	}

	return true, nil
}

/// Private Functions

// validateCreateIndexRequestPayload validates create index request payloads
func (es *ElasthinkSDK) validateCreateIndexRequestPayload(documentID int64, documentType, documentName string) error {
	err := es.isValidFromCustomDocumentType(documentType)
	if err != nil {
		return err
	}

	if documentID <= 0 {
		return errors.New("Invalid Document ID")
	}

	if len(strings.Trim(documentName, " ")) == 0 {
		return errors.New("Document Name must not be empty")
	}

	return nil
}

// Validate is the document type is valid or not
func (es *ElasthinkSDK) isValidFromCustomDocumentType(documentType string) error {
	if _, ok := es.availableDocumentType[documentType]; ok {
		return nil
	}
	return errors.New("Invalid Document Type")
}
