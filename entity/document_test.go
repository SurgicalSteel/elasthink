package entity

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
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidFromCustomDocumentType(t *testing.T) {
	type tcase struct {
		documentType    DocumentType
		documentTypeMap map[DocumentType]int
		expectedError   error
	}
	var testDocumentTypeA DocumentType = "type_a"
	var testDocumentTypeB DocumentType = "type_b"
	var testDocumentTypeC DocumentType = "type_c"
	var testDocumentTypeD DocumentType = "type_d"
	testCases := make(map[string]tcase)
	testCases["Valid Document Type"] = tcase{
		documentType: testDocumentTypeD,
		documentTypeMap: map[DocumentType]int{
			testDocumentTypeA: 1,
			testDocumentTypeB: 1,
			testDocumentTypeC: 1,
			testDocumentTypeD: 1,
		},
		expectedError: nil,
	}
	testCases["Invalid Document Type"] = tcase{
		documentType: testDocumentTypeD,
		documentTypeMap: map[DocumentType]int{
			testDocumentTypeA: 1,
			testDocumentTypeB: 1,
			testDocumentTypeC: 1,
		},
		expectedError: errors.New("Invalid Document Type"),
	}

	for ktc, vtc := range testCases {
		fmt.Println("doing test on IsValid of Document Type with test case:", ktc)
		actual := vtc.documentType.IsValidFromCustomDocumentType(vtc.documentTypeMap)
		assert.Equal(t, vtc.expectedError, actual)
	}

}
