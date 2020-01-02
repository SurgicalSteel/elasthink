package service

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
