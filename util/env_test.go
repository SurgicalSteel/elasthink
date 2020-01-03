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

func TestGetEnv(t *testing.T) {
	type tcase struct {
		sourceEnvString string
		expected        string
	}
	testCases := make(map[string]tcase)

	testCases["env staging with uppercase and space"] = tcase{
		sourceEnvString: " STAGING",
		expected:        "staging",
	}

	testCases["env stg with space"] = tcase{
		sourceEnvString: " stg",
		expected:        "staging",
	}

	testCases["empty string"] = tcase{
		sourceEnvString: "",
		expected:        "development",
	}

	testCases["invalid env"] = tcase{
		sourceEnvString: "?B$A!J&I*N{G}A%N+",
		expected:        "development",
	}

	for ktc, vtc := range testCases {
		fmt.Println("doing test on GetEnv with test case:", ktc)
		actual := GetEnv(vtc.sourceEnvString)
		assert.Equal(t, vtc.expected, actual)
	}

}
