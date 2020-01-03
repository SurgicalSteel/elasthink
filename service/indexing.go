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
	"github.com/SurgicalSteel/elasthink/module"
	"github.com/SurgicalSteel/elasthink/util"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

//HandleCreateIndex handles create index (from internal endpoint)
func HandleCreateIndex(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	documentType := vars["document_type"]
	documentIDRaw := vars["document_id"]
	documentID := util.StringToInt64(documentIDRaw)

	var requestPayload module.CreateIndexRequestPayload

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &requestPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := module.CreateIndex(ctx, documentID, documentType, requestPayload)
	responsePayload := constructResponsePayload(response)

	responsePayloadJSON, err := json.Marshal(responsePayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(response.StatusCode)
	w.Write(responsePayloadJSON)
}

//HandleUpdateIndex handles update index (from internal endpoint)
func HandleUpdateIndex(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	documentType := vars["document_type"]
	documentIDRaw := vars["document_id"]
	documentID := util.StringToInt64(documentIDRaw)

	var requestPayload module.UpdateIndexRequestPayload

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &requestPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := module.UpdateIndex(ctx, documentID, documentType, requestPayload)
	responsePayload := constructResponsePayload(response)

	responsePayloadJSON, err := json.Marshal(responsePayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(response.StatusCode)
	w.Write(responsePayloadJSON)
}
