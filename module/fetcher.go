package module

import (
	"fmt"
	"github.com/SurgicalSteel/elasthink/entity"
	"github.com/SurgicalSteel/elasthink/util"
	"log"
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
