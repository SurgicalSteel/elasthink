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
import (
	"fmt"
	"reflect"
	"testing"
)

func TestTokenizer(t *testing.T) {
	type tcase struct {
		sourceString           string
		isUsingStopwordRemoval bool
		stopwordSet            map[string]int
		expected               map[string]int
	}

	testCases := make(map[string]tcase)

	testCases["normal testcase"] = tcase{
		sourceString:           "quick brown fox",
		isUsingStopwordRemoval: false,
		stopwordSet:            make(map[string]int),
		expected: map[string]int{
			"quick": 1,
			"brown": 1,
			"fox":   1,
		},
	}

	testCases["using stopwords"] = tcase{
		sourceString:           "quick brown fox",
		isUsingStopwordRemoval: true,
		stopwordSet: map[string]int{
			"quick": 1,
			"fox":   1,
		},
		expected: map[string]int{
			"brown": 1,
		},
	}

	testCases["invalid word"] = tcase{
		sourceString:           " ==== ==== ",
		isUsingStopwordRemoval: false,
		stopwordSet:            make(map[string]int),
		expected:               make(map[string]int),
	}

	for ktc, vtc := range testCases {
		fmt.Println("doing test on Tokenizer with test case:", ktc)
		actual := Tokenize(vtc.sourceString, vtc.isUsingStopwordRemoval, vtc.stopwordSet)
		isEqual := reflect.DeepEqual(vtc.expected, actual)
		if !isEqual {
			t.Fatal("Result WordSet is not same with what we expected.")
			continue
		}
	}
}
