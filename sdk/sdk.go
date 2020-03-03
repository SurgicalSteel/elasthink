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

// CreateIndexSpec is the spec of CreateIndex function
// documentType is the type of the document
// documentName is the name of the document
// documentID is the id of the document
type CreateIndexSpec struct {
	documentType string
	documentName string
	documentID   int64
}

// UpdateIndexSpec is the spec of UpdateIndex function
// oldDocumentName is the name of the old document which will be replaced by new document
// newDocumentName is the name of the new document
// documentID is the id of the document
type UpdateIndexSpec struct {
	documentType    string
	oldDocumentName string
	newDocumentName string
	documentID      int64
}

// SearchSpec is the spec of Search function
type SearchSpec struct {
	documentType string
	searchTerm   string
}

// SearchResultRankData is the search result
type SearchResultRankData struct {
	ID        int64
	ShowCount int
	Rank      int
}

// SearchResult is the result of Search
type SearchResult struct {
	RankedResultList []SearchResultRankData
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
func (es *ElasthinkSDK) CreateIndex(spec CreateIndexSpec) (bool, error) {
	documentID := spec.documentID
	documentType := spec.documentType
	documentName := spec.documentName

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

//UpdateIndex is function to update previously created index
func (es *ElasthinkSDK) UpdateIndex(spec UpdateIndexSpec) (bool, error) {
	documentID := spec.documentID
	documentType := spec.documentType
	oldDocumentName := spec.oldDocumentName
	newDocumentName := spec.newDocumentName

	// Validate
	err := es.validateUpdateIndexRequestPayload(documentID, documentType, spec.oldDocumentName, spec.newDocumentName)
	if err != nil {
		return false, err
	}

	// Tokenize
	stopword := make(map[string]int)
	for _, k := range es.stopWordRemovalData {
		stopword[k] = 1
	}
	redis := es.redis
	oldDocumentNameSet := util.Tokenize(oldDocumentName, es.isUsingStopWordsRemoval, stopword)
	newDocumentNameSet := util.Tokenize(newDocumentName, es.isUsingStopWordsRemoval, stopword)

	docType := spec.documentType

	// remove old document indexes
	isErrorRemoveExist := false
	errorRemoveKeys := ""

	for k := range oldDocumentNameSet {
		key := fmt.Sprintf("%s:%s", docType, k)
		value := make([]interface{}, 1)
		value[0] = fmt.Sprintf("%d", documentID)
		_, err = redis.SRem(key, value)
		if err != nil {
			isErrorRemoveExist = true
			errorRemoveKeys = errorRemoveKeys + " " + key + ","
			log.Println("[MODULE][UPDATE INDEX] failed to remove index on key :", key, "and document ID:", documentID)
			continue
		}
	}

	// add new document indexes
	isErrorAddExist := false
	errorAddKeys := ""

	for k := range newDocumentNameSet {
		key := fmt.Sprintf("%s:%s", docType, k)
		value := make([]interface{}, 1)
		value[0] = fmt.Sprintf("%d", documentID)
		_, err = redis.SAdd(key, value)
		if err != nil {
			isErrorAddExist = true
			errorAddKeys = errorAddKeys + " " + key + ","
			log.Println("[MODULE][UPDATE INDEX] failed to add index on key :", key, "and document ID:", documentID)
			continue
		}
	}

	if isErrorAddExist || isErrorRemoveExist {
		errorRemoveKeys = strings.TrimRight(errorRemoveKeys, ",")
		errorRemoveKeys = strings.TrimLeft(errorRemoveKeys, " ")

		errorAddKeys = strings.TrimRight(errorAddKeys, ",")
		errorAddKeys = strings.TrimLeft(errorAddKeys, " ")

		errorMessage := fmt.Sprintf("Error on removing following keys: %s and/or Error on adding following keys: %s", errorRemoveKeys, errorAddKeys)
		return false, errors.New(errorMessage)
	}

	return true, nil
}

//Search is the core function of searching a document
func (es *ElasthinkSDK) Search(spec SearchSpec) (SearchResult, error) {
	// err := validateSearchRequestPayload(documentType, requestPayload.SearchTerm)
	// if err != nil {
	// 	return Response{
	// 		StatusCode:   http.StatusBadRequest,
	// 		ErrorMessage: err.Error(),
	// 		Data:         nil,
	// 	}
	// }

	// searchTermSet := util.Tokenize(requestPayload.SearchTerm, moduleObj.IsUsingStopwordRemoval, moduleObj.StopwordSet)
	// if len(searchTermSet) == 0 {
	// 	return Response{
	// 		StatusCode:   http.StatusOK,
	// 		ErrorMessage: "",
	// 		Data:         nil,
	// 	}
	// }

	// docType := getDocumentType(documentType, entity.Entity.GetDocumentTypes())

	// wordIndexSets := fetchWordIndexSets(docType, searchTermSet)

	// if len(wordIndexSets) == 0 {
	// 	return Response{
	// 		StatusCode:   http.StatusOK,
	// 		ErrorMessage: "",
	// 		Data:         nil,
	// 	}
	// }

	// rankedSearchResult := rankSearchResult(wordIndexSets)
	// searchResponsePayload := SearchResponsePayload{RankedResultList: rankedSearchResult}

	// return Response{
	// 	StatusCode:   http.StatusOK,
	// 	ErrorMessage: "",
	// 	Data:         searchResponsePayload,
	// }
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

func (es *ElasthinkSDK) validateUpdateIndexRequestPayload(documentID int64, documentType, oldDocumentName, newDocumentName string) error {
	err := es.isValidFromCustomDocumentType(documentType)
	if err != nil {
		return err
	}

	if documentID <= 0 {
		return errors.New("Invalid Document ID")
	}

	if len(strings.Trim(oldDocumentName, " ")) == 0 {
		return errors.New("Old Document Name must not be empty")
	}

	if len(strings.Trim(newDocumentName, " ")) == 0 {
		return errors.New("Document Name must not be empty")
	}

	return nil
}
