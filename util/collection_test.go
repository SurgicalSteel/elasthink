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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceStringToInt64(t *testing.T) {
	type tcase struct {
		sourceSlice []string
		expected    []int64
	}
	testCases := make(map[string]tcase)

	testCases["empty slice"] = tcase{
		sourceSlice: make([]string, 0),
		expected:    make([]int64, 0),
	}

	testCases["slice of numbers in string"] = tcase{
		sourceSlice: []string{"1", "1", "2", "3", "5", "8", "13"},
		expected:    []int64{1, 1, 2, 3, 5, 8, 13},
	}

	testCases["slice of numbers in string and some random characters"] = tcase{
		sourceSlice: []string{"111", "222", "333", "???", "555", "???"},
		expected:    []int64{111, 222, 333, 0, 555, 0},
	}

	for ktc, vtc := range testCases {
		fmt.Println("doing test on SliceStringToInt64 with test case:", ktc)
		actual := SliceStringToInt64(vtc.sourceSlice)
		assert.ElementsMatch(t, actual, vtc.expected)
	}

}
