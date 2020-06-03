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
	"context"
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/SurgicalSteel/elasthink/entity"
)

func validateKeywordSuggestionRequest(documentType, prefix string) error {
	if len(strings.Trim(prefix, " ")) == 0 {
		return errors.New("Keyword prefix is required to get suggested keywords")
	}

	if len(strings.Trim(documentType, " ")) == 0 {
		return errors.New("Document Type is required")
	}

	err := validateDocumentType(documentType, entity.Entity.GetDocumentTypes())
	if err != nil {
		return err
	}

	return nil
}

//KeywordSuggestionResponsePayload is the universal response payload for keyword suggestion API handler
type KeywordSuggestionResponsePayload struct {
	SortedKeywords []string `json:"sortedKeywords"`
}

//SuggestKeywords is the core function for keyword suggestion (by document type and prefix)
func SuggestKeywords(ctx context.Context, documentType, prefix string) Response {
	err := validateKeywordSuggestionRequest(documentType, prefix)
	if err != nil {
		return Response{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: err.Error(),
			Data:         nil,
		}
	}
	prefix = strings.ToLower(prefix)
	docType := getDocumentType(documentType, entity.Entity.GetDocumentTypes())
	keywords, err := fetchKeywords(docType, prefix)
	if err != nil {
		return Response{
			StatusCode:   http.StatusInternalServerError,
			ErrorMessage: "There's an error when suggesting keywords.",
			Data:         nil,
		}
	}
	sort.Strings(keywords)
	return Response{
		StatusCode: http.StatusOK,
		Data:       keywords,
	}
}
