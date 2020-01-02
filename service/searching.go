package service

import (
	"context"
	"encoding/json"
	"github.com/SurgicalSteel/elasthink/module"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

//HandleSearch handles the search for a document id (from internal & external endpoint)
func HandleSearch(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	documentType := vars["document_type"]

	var requestPayload module.SearchRequestPayload

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

	response := module.Search(ctx, documentType, requestPayload)
	responsePayload := constructResponsePayload(response)

	responsePayloadJSON, err := json.Marshal(responsePayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(response.StatusCode)
	w.Write(responsePayloadJSON)
}
