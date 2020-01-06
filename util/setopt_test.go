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

func TestCreateWordSet(t *testing.T) {
	type tcase struct {
		sourceWordSlice []string
		expected        map[string]int
	}

	testCases := make(map[string]tcase)

	testCases["empty word slice"] = tcase{
		sourceWordSlice: make([]string, 0),
		expected:        make(map[string]int),
	}

	testCases["normal word slice"] = tcase{
		sourceWordSlice: []string{"invisible", "means", "can", "not", "be", "seen"},
		expected: map[string]int{
			"invisible": 1,
			"means":     1,
			"can":       1,
			"not":       1,
			"be":        1,
			"seen":      1,
		},
	}

	testCases["repeated normal word slice"] = tcase{
		sourceWordSlice: []string{"this", "one", "and", "that", "one", "are", "yours", "and", "shut", "up"},
		expected: map[string]int{
			"this":  1,
			"one":   1,
			"and":   1,
			"that":  1,
			"are":   1,
			"yours": 1,
			"shut":  1,
			"up":    1,
		},
	}

	for ktc, vtc := range testCases {
		fmt.Println("doing test on CreateWordSet with test case:", ktc)
		actual := CreateWordSet(vtc.sourceWordSlice)
		isEqual := reflect.DeepEqual(vtc.expected, actual)
		if !isEqual {
			t.Fatal("Result WordSet is not same with what we expected.")
			continue
		}
	}
}
func TestWordsSetUnion(t *testing.T) {
	type tcase struct {
		sourceSetA map[string]int
		sourceSetB map[string]int
		expected   map[string]int
	}

	testCases := make(map[string]tcase)

	testCases["empty setA and empty setB"] = tcase{
		sourceSetA: make(map[string]int),
		sourceSetB: make(map[string]int),
		expected:   make(map[string]int),
	}

	testCases["empty setA and filled setB"] = tcase{
		sourceSetA: make(map[string]int),
		sourceSetB: map[string]int{
			"the":   1,
			"quick": 1,
			"brown": 1,
			"fox":   1,
			"jumps": 1,
			"over":  1,
			"lazy":  1,
			"dog":   1,
		},
		expected: map[string]int{
			"the":   1,
			"quick": 1,
			"brown": 1,
			"fox":   1,
			"jumps": 1,
			"over":  1,
			"lazy":  1,
			"dog":   1,
		},
	}

	testCases["filled setA and empty setB"] = tcase{
		sourceSetA: map[string]int{
			"the":   1,
			"quick": 1,
			"brown": 1,
			"fox":   1,
			"jumps": 1,
			"over":  1,
			"lazy":  1,
			"dog":   1,
		},
		sourceSetB: make(map[string]int),
		expected: map[string]int{
			"the":   1,
			"quick": 1,
			"brown": 1,
			"fox":   1,
			"jumps": 1,
			"over":  1,
			"lazy":  1,
			"dog":   1,
		},
	}

	testCases["filled setA and filled setB with no overlapping items"] = tcase{
		sourceSetA: map[string]int{
			"invincible": 1,
			"means":      1,
			"can":        1,
			"not":        1,
			"be":         1,
			"seen":       1,
		},
		sourceSetB: map[string]int{
			"the":   1,
			"quick": 1,
			"brown": 1,
			"fox":   1,
			"jumps": 1,
			"over":  1,
			"lazy":  1,
			"dog":   1,
		},
		expected: map[string]int{
			"invincible": 1,
			"means":      1,
			"can":        1,
			"not":        1,
			"be":         1,
			"seen":       1,
			"the":        1,
			"quick":      1,
			"brown":      1,
			"fox":        1,
			"jumps":      1,
			"over":       1,
			"lazy":       1,
			"dog":        1,
		},
	}

	testCases["filled setA and filled setB with overlapping items"] = tcase{
		sourceSetA: map[string]int{
			"it":         1,
			"is":         1,
			"fun":        1,
			"to":         1,
			"do":         1,
			"the":        1,
			"impossible": 1,
		},
		sourceSetB: map[string]int{
			"impossible": 1,
			"is":         1,
			"nothing":    1,
		},
		expected: map[string]int{
			"it":         1,
			"is":         1,
			"fun":        1,
			"to":         1,
			"do":         1,
			"the":        1,
			"impossible": 1,
			"nothing":    1,
		},
	}

	for ktc, vtc := range testCases {
		fmt.Println("doing test on WordSetUnion with test case:", ktc)
		actual := WordsSetUnion(vtc.sourceSetA, vtc.sourceSetB)
		isEqual := reflect.DeepEqual(vtc.expected, actual)
		if !isEqual {
			t.Fatal("Result WordSet is not same with what we expected.")
			continue
		}
	}

}
func TestWordsSetSubtraction(t *testing.T) {
	type tcase struct {
		sourceSetA map[string]int
		sourceSetB map[string]int
		expected   map[string]int
	}

	testCases := make(map[string]tcase)

	testCases["empty setA and empty setB"] = tcase{
		sourceSetA: make(map[string]int),
		sourceSetB: make(map[string]int),
		expected:   make(map[string]int),
	}

	testCases["empty setA and filled setB"] = tcase{
		sourceSetA: make(map[string]int),
		sourceSetB: map[string]int{
			"the":   1,
			"quick": 1,
			"brown": 1,
			"fox":   1,
			"jumps": 1,
			"over":  1,
			"lazy":  1,
			"dog":   1,
		},
		expected: make(map[string]int),
	}

	testCases["filled setA and empty setB"] = tcase{
		sourceSetA: map[string]int{
			"the":   1,
			"quick": 1,
			"brown": 1,
			"fox":   1,
			"jumps": 1,
			"over":  1,
			"lazy":  1,
			"dog":   1,
		},
		sourceSetB: make(map[string]int),
		expected: map[string]int{
			"the":   1,
			"quick": 1,
			"brown": 1,
			"fox":   1,
			"jumps": 1,
			"over":  1,
			"lazy":  1,
			"dog":   1,
		},
	}

	testCases["filled setA and filled setB with no overlapping items"] = tcase{
		sourceSetA: map[string]int{
			"invincible": 1,
			"means":      1,
			"can":        1,
			"not":        1,
			"be":         1,
			"seen":       1,
		},
		sourceSetB: map[string]int{
			"the":   1,
			"quick": 1,
			"brown": 1,
			"fox":   1,
			"jumps": 1,
			"over":  1,
			"lazy":  1,
			"dog":   1,
		},
		expected: map[string]int{
			"invincible": 1,
			"means":      1,
			"can":        1,
			"not":        1,
			"be":         1,
			"seen":       1,
		},
	}

	testCases["filled setA and filled setB with overlapping items"] = tcase{
		sourceSetA: map[string]int{
			"it":         1,
			"is":         1,
			"fun":        1,
			"to":         1,
			"do":         1,
			"the":        1,
			"impossible": 1,
		},
		sourceSetB: map[string]int{
			"impossible": 1,
			"is":         1,
			"nothing":    1,
		},
		expected: map[string]int{
			"it":  1,
			"fun": 1,
			"to":  1,
			"do":  1,
			"the": 1,
		},
	}

	for ktc, vtc := range testCases {
		fmt.Println("doing test on WordSetSubtraction with test case:", ktc)
		actual := WordsSetSubtraction(vtc.sourceSetA, vtc.sourceSetB)
		isEqual := reflect.DeepEqual(vtc.expected, actual)
		if !isEqual {
			t.Fatal("Result WordSet is not same with what we expected.")
			continue
		}
	}
}
func TestWordsSetIntersection(t *testing.T) {
	type tcase struct {
		sourceSetA map[string]int
		sourceSetB map[string]int
		expected   map[string]int
	}

	testCases := make(map[string]tcase)

	testCases["empty setA and empty setB"] = tcase{
		sourceSetA: make(map[string]int),
		sourceSetB: make(map[string]int),
		expected:   make(map[string]int),
	}

	testCases["empty setA and filled setB"] = tcase{
		sourceSetA: make(map[string]int),
		sourceSetB: map[string]int{
			"the":   1,
			"quick": 1,
			"brown": 1,
			"fox":   1,
			"jumps": 1,
			"over":  1,
			"lazy":  1,
			"dog":   1,
		},
		expected: make(map[string]int),
	}

	testCases["filled setA and empty setB"] = tcase{
		sourceSetA: map[string]int{
			"the":   1,
			"quick": 1,
			"brown": 1,
			"fox":   1,
			"jumps": 1,
			"over":  1,
			"lazy":  1,
			"dog":   1,
		},
		sourceSetB: make(map[string]int),
		expected:   make(map[string]int),
	}

	testCases["filled setA and filled setB with no overlapping items"] = tcase{
		sourceSetA: map[string]int{
			"invincible": 1,
			"means":      1,
			"can":        1,
			"not":        1,
			"be":         1,
			"seen":       1,
		},
		sourceSetB: map[string]int{
			"the":   1,
			"quick": 1,
			"brown": 1,
			"fox":   1,
			"jumps": 1,
			"over":  1,
			"lazy":  1,
			"dog":   1,
		},
		expected: make(map[string]int),
	}

	testCases["filled setA and filled setB with overlapping items"] = tcase{
		sourceSetA: map[string]int{
			"it":         1,
			"is":         1,
			"fun":        1,
			"to":         1,
			"do":         1,
			"the":        1,
			"impossible": 1,
		},
		sourceSetB: map[string]int{
			"impossible": 1,
			"is":         1,
			"nothing":    1,
		},
		expected: map[string]int{
			"impossible": 1,
			"is":         1,
		},
	}

	for ktc, vtc := range testCases {
		fmt.Println("doing test on WordSetIntersection with test case:", ktc)
		actual := WordsSetIntersection(vtc.sourceSetA, vtc.sourceSetB)
		isEqual := reflect.DeepEqual(vtc.expected, actual)
		if !isEqual {
			t.Fatal("Result WordSet is not same with what we expected.")
			continue
		}
	}
}
