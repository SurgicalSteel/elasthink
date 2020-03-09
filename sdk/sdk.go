package module

import (
	"errors"
	"fmt"
	"log"
	"sort"
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

// ElasthinkSDK is the main struct of elasthink SDK, initialized using initalize function
type ElasthinkSDK struct {
	redis                   *redis.Redis
	isUsingStopWordsRemoval bool
	stopWordRemovalData     []string
	availableDocumentType   map[string]int
}

// InitializeSpec is the payload for initialize Elasthink SDK
type InitializeSpec struct {
	redisConfig RedisConfig
	sdkConfig   SdkConfig
}

//RedisConfig is the basic configuration for a redis connection
type RedisConfig struct {
	Address   string
	MaxActive int
	MaxIdle   int
	Timeout   int
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

// SearchResultRankData is the search result datum
type SearchResultRankData struct {
	ID        int64
	ShowCount int
	Rank      int
}

// SearchResult is the result of Search, it have array of search result datum
type SearchResult struct {
	RankedResultList RankByShowCount
}

//RankByShowCount is the additional struct for document ranking purpose based on its ShowCount
type RankByShowCount []SearchResultRankData

// Len overrides Len function of RankByShowCount
func (r RankByShowCount) Len() int { return len(r) }

// Less overrides Less function of RankByShowCount
func (r RankByShowCount) Less(i, j int) bool {
	return r[i].ShowCount > r[j].ShowCount
}

// Swap overrides Swap function of RankByShowCount
func (r RankByShowCount) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

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
	err := es.validateCreateIndexSpec(documentID, documentType, documentName)
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
			log.Println("[SDK_REDIS][CREATE INDEX] failed to add index on key :", key, "and document ID:", documentID)
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
	redis := es.redis

	// Validate
	err := es.validateUpdateIndexSpec(documentID, documentType, spec.oldDocumentName, spec.newDocumentName)
	if err != nil {
		return false, err
	}

	// Tokenize
	stopword := make(map[string]int)
	for _, k := range es.stopWordRemovalData {
		stopword[k] = 1
	}
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
			log.Println("[SDK_REDIS][UPDATE INDEX] failed to remove index on key :", key, "and document ID:", documentID)
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
			log.Println("[SDK_REDIS][UPDATE INDEX] failed to add index on key :", key, "and document ID:", documentID)
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
	searchTerm := spec.searchTerm
	documentType := spec.documentType

	ret := SearchResult{RankedResultList: make([]SearchResultRankData, 0)}

	err := es.validateSearchSpec(documentType, searchTerm)
	if err != nil {
		return ret, err
	}

	stopword := make(map[string]int)
	for _, k := range es.stopWordRemovalData {
		stopword[k] = 1
	}
	searchTermSet := util.Tokenize(searchTerm, es.isUsingStopWordsRemoval, stopword)
	if len(searchTermSet) == 0 {
		return ret, nil
	}

	docType := documentType

	wordIndexSets := es.fetchWordIndexSets(docType, searchTermSet)

	if len(wordIndexSets) == 0 {
		return ret, err
	}

	rankedSearchResult := rankSearchResult(wordIndexSets)
	ret.RankedResultList = rankedSearchResult

	return ret, nil
}

/// Private Functions

// validateCreateIndexSpec validates create index spec
func (es *ElasthinkSDK) validateCreateIndexSpec(documentID int64, documentType, documentName string) error {
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

// validateUpdateIndexSpec validate update index spec
func (es *ElasthinkSDK) validateUpdateIndexSpec(documentID int64, documentType, oldDocumentName, newDocumentName string) error {
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

// validateSearchSpec validate search spec
func (es *ElasthinkSDK) validateSearchSpec(documentType, searchTerm string) error {
	if len(strings.Trim(searchTerm, " ")) == 0 {
		return errors.New("Search Term is required")
	}

	if len(strings.Trim(documentType, " ")) == 0 {
		return errors.New("Document Type is required")
	}

	err := es.isValidFromCustomDocumentType(documentType)
	if err != nil {
		return err
	}

	return nil
}

// fetchWordIndexSets
func (es *ElasthinkSDK) fetchWordIndexSets(documentType string, searchTermSet map[string]int) map[string][]int64 {
	result := make(map[string][]int64)

	// set key format --> documentType:word
	for k := range searchTermSet {
		key := fmt.Sprintf("%s:%s", documentType, k)
		members, err := es.redis.SMembers(key)
		if err != nil {
			log.Println("[SDK_REDIS][FETCHER] Failed to get members of key :", key)
			continue
		}
		documentIds := util.SliceStringToInt64(members)
		result[k] = documentIds
	}

	return result
}

//rankSearchResult ranks search result (document id by its appeareance count). word indexes is a map with word as a key and slice of ids as value. Returns ordered search result rank slice.
func rankSearchResult(wordIndexes map[string][]int64) []SearchResultRankData {
	counterMap := make(map[int64]int)
	for _, ids := range wordIndexes {
		for i := 0; i < len(ids); i++ {
			if vcm, ok := counterMap[ids[i]]; ok {
				counterMap[ids[i]] = vcm + 1
			} else {
				counterMap[ids[i]] = 1
			}
		}
	}

	result := make([]SearchResultRankData, len(counterMap))

	iterator := 0
	for kcm, vcm := range counterMap {
		result[iterator] = SearchResultRankData{
			ID:        kcm,
			ShowCount: vcm,
		}
		iterator++
	}

	//sort by appeareance count (descending)
	sort.Sort(RankByShowCount(result))

	//assign rank to each search result data
	for i := 0; i < len(result); i++ {
		temp := result[i]
		temp.Rank = i + 1
		result[i] = temp
	}

	return result
}
