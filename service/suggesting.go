package service

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
	"encoding/json"
	"net/http"

	"github.com/SurgicalSteel/elasthink/module"
	"github.com/gorilla/mux"
)

func HandleKeywordSuggestion(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	documentType := vars["document_type"]
	prefix := vars["prefix"]

	response := module.SuggestKeywords(ctx, documentType, prefix)
	responsePayload := constructResponsePayload(response)

	responsePayloadJSON, err := json.Marshal(responsePayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(response.StatusCode)
	w.Write(responsePayloadJSON)
}
