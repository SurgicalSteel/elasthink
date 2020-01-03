package util

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
import ()

//CreateWordSet create a set of words from a slice of string
func CreateWordSet(words []string) map[string]int {
	result := make(map[string]int)

	for _, vw := range words {
		if _, okr := result[vw]; !okr {
			result[vw] = 1
		}
	}

	return result
}

//WordsSetUnion do a union from two given words sets
func WordsSetUnion(ma, mb map[string]int) map[string]int {
	result := make(map[string]int)
	for ka, _ := range ma {
		if _, oka := result[ka]; !oka {
			result[ka] = 1
		}
	}

	for kb, _ := range mb {
		if _, okb := result[kb]; !okb {
			result[kb] = 1
		}
	}
	return result
}

//WordsSetSubtraction do a substraction between two sets ma - mb
func WordsSetSubtraction(ma, mb map[string]int) map[string]int {
	result := ma
	for kb, _ := range mb {
		if _, okb := result[kb]; okb {
			delete(result, kb)
		}
	}
	return result
}

//WordsSetIntersection find the intersection between two sets (ma & mb)
func WordsSetIntersection(ma, mb map[string]int) map[string]int {
	result := make(map[string]int)
	for ka, _ := range ma {
		if _, okb := mb[ka]; okb {
			result[ka] = 1
		}
	}
	return result
}
