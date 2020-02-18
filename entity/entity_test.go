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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitializeEntity(t *testing.T) {
	entityDataMock := entityData{}
	stopwordDataMock := StopwordData{
		Words: []string{"tidak", "tidakkah", "tidaklah"},
	}

	assert.Nil(t, entityDataMock.documentTypes)
	assert.Equal(t, 0, len(entityDataMock.stopwordData.Words))

	entityDataMock.Initialize(stopwordDataMock)

	assert.NotNil(t, entityDataMock.documentTypes)
	assert.Equal(t, stopwordDataMock.Words, entityDataMock.stopwordData.Words)
	if len(entityDataMock.documentTypes) == 0 {
		t.Error("Expected non empty Document Types, but found empty document types")
	}
}

func TestGetDocumentTypes(t *testing.T) {
	entityDataMock := entityData{}
	stopwordDataMock := StopwordData{
		Words: []string{"tidak", "tidakkah", "tidaklah"},
	}

	assert.Nil(t, entityDataMock.GetDocumentTypes())
	assert.Equal(t, 0, len(entityDataMock.stopwordData.Words))

	entityDataMock.Initialize(stopwordDataMock)

	assert.NotNil(t, entityDataMock.GetDocumentTypes())

	if len(entityDataMock.GetDocumentTypes()) == 0 {
		t.Error("Expected non empty Document Types, but found empty document types")
	}
}

func TestGetStopwordData(t *testing.T) {
	entityDataMock := entityData{}
	stopwordDataMock := StopwordData{
		Words: []string{"tidak", "tidakkah", "tidaklah"},
	}

	assert.Nil(t, entityDataMock.GetDocumentTypes())
	assert.Equal(t, 0, len(entityDataMock.stopwordData.Words))
	entityDataMock.Initialize(stopwordDataMock)
	assert.Equal(t, 3, len(entityDataMock.stopwordData.Words))

	actualStopwordData := entityDataMock.GetStopwordData()
	assert.Equal(t, stopwordDataMock, actualStopwordData)

	if len(entityDataMock.GetDocumentTypes()) == 0 {
		t.Error("Expected non empty Document Types, but found empty document types")
	}
}
