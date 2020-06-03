package module

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
	"fmt"
	"log"
	"strings"

	"github.com/SurgicalSteel/elasthink/entity"
	"github.com/SurgicalSteel/elasthink/util"
)

func fetchWordIndexSets(documentType entity.DocumentType, searchTermSet map[string]int) map[string][]int64 {
	result := make(map[string][]int64)

	// set key format --> documentType:word
	for k, _ := range searchTermSet {
		key := fmt.Sprintf("%s:%s", documentType, k)
		members, err := moduleObj.Redis.SMembers(key)
		if err != nil {
			log.Println("[MODULE][FETCHER] Failed to get members of key :", key)
			continue
		}
		documentIds := util.SliceStringToInt64(members)
		result[k] = documentIds
	}

	return result
}

func fetchKeywords(documentType entity.DocumentType, prefix string) ([]string, error) {
	prefixKey := fmt.Sprintf("%s:%s", documentType, prefix)
	rawKeys, err := moduleObj.Redis.KeysPrefix(prefixKey)
	if err != nil {
		log.Printf("[MODULE][FETCHER] Failed to get keys with prefix :%s Detail :%s\n", prefixKey, err.Error())
		return []string{}, err
	}
	finalKeywords := make([]string, len(rawKeys))
	trimPrefix := fmt.Sprintf("%s:", documentType)
	for i := 0; i < len(rawKeys); i++ {
		rawKey := rawKeys[i]
		finalKeywords[i] = strings.TrimPrefix(rawKey, trimPrefix)
	}
	return finalKeywords, nil
}
