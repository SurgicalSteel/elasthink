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
	"github.com/SurgicalSteel/elasthink/entity"
	"sort"
)

//RankByShowCount is the additional struct for document ranking purpose based on its ShowCount
type RankByShowCount []entity.SearchResultRankData

func (r RankByShowCount) Len() int           { return len(r) }
func (r RankByShowCount) Less(i, j int) bool { return r[i].ShowCount > r[j].ShowCount }
func (r RankByShowCount) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

//rankSearchResult ranks search result (document id by its appeareance count). word indexes is a map with word as a key and slice of ids as value. Returns ordered search result rank slice.
func rankSearchResult(wordIndexes map[string][]int64) []entity.SearchResultRankData {
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

	result := make([]entity.SearchResultRankData, len(counterMap))

	iterator := 0
	for kcm, vcm := range counterMap {
		result[iterator] = entity.SearchResultRankData{
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
