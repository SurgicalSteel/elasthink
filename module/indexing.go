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
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/SurgicalSteel/elasthink/entity"
	"github.com/SurgicalSteel/elasthink/util"
)

//CreateIndexRequestPayload is the universal request payload for create index handler
type CreateIndexRequestPayload struct {
	DocumentName string `json:"documentName"`
}

func validateCreateIndexRequestPayload(documentID int64, documentType, documentName string) error {
	err := validateDocumentType(documentType, entity.Entity.GetDocumentTypes())
	if err != nil {
		return err
	}

	if documentID <= 0 {
		return errors.New("Invalid Document ID")
	}

	if len(strings.Trim(documentName, " ")) == 0 {
		return errors.New("Document Name must not be empty")
	}

	return nil
}

//CreateIndex is the core function to create an index of a document
func CreateIndex(ctx context.Context, documentID int64, documentType string, requestPayload CreateIndexRequestPayload) Response {
	err := validateCreateIndexRequestPayload(documentID, documentType, requestPayload.DocumentName)
	if err != nil {
		return Response{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: err.Error(),
			Data:         nil,
		}
	}

	documentNameSet := util.Tokenize(requestPayload.DocumentName, moduleObj.IsUsingStopwordRemoval, moduleObj.StopwordSet)

	docType := getDocumentType(documentType, entity.Entity.GetDocumentTypes())
	errorExist := false
	errorKeys := ""

	for k, _ := range documentNameSet {
		key := fmt.Sprintf("%s:%s", docType, k)
		value := make([]interface{}, 1)
		value[0] = fmt.Sprintf("%d", documentID)
		_, err = moduleObj.Redis.SAdd(key, value)
		if err != nil {
			errorExist = true
			errorKeys = errorKeys + " " + key + ","
			log.Println("[MODULE][CREATE INDEX] failed to add index on key :", key, "and document ID:", documentID)
			continue
		}
	}

	if errorExist {
		errorKeys = strings.TrimRight(errorKeys, ",")
		errorKeys = strings.TrimLeft(errorKeys, " ")
		return Response{
			StatusCode:   http.StatusInternalServerError,
			ErrorMessage: fmt.Sprintf("Error on adding following keys :%s", errorKeys),
			Data:         nil,
		}
	} else {
		return Response{
			StatusCode:   http.StatusOK,
			ErrorMessage: "",
			Data:         nil,
		}
	}
}

//UpdateIndexRequestPayload is the universal request payload for update index handler
type UpdateIndexRequestPayload struct {
	OldDocumentName string `json:"oldDocumentName"`
	NewDocumentName string `json:"newDocumentName"`
}

func validateUpdateIndexRequestPayload(documentID int64, documentType, oldDocumentName, newDocumentName string) error {
	err := validateDocumentType(documentType, entity.Entity.GetDocumentTypes())
	if err != nil {
		return err
	}

	if documentID <= 0 {
		return errors.New("Invalid Document ID")
	}

	if len(strings.Trim(oldDocumentName, " ")) == 0 {
		return errors.New("Old Document Name must not be empty")
	}

	if len(strings.Trim(newDocumentName, " ")) == 0 {
		return errors.New("Document Name must not be empty")
	}

	return nil
}

//UpdateIndex is the core function to update the index of a document. This requires old document name and new document name
func UpdateIndex(ctx context.Context, documentID int64, documentType string, requestPayload UpdateIndexRequestPayload) Response {
	err := validateUpdateIndexRequestPayload(documentID, documentType, requestPayload.OldDocumentName, requestPayload.NewDocumentName)
	if err != nil {
		return Response{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: err.Error(),
			Data:         nil,
		}
	}

	oldDocumentNameSet := util.Tokenize(requestPayload.OldDocumentName, moduleObj.IsUsingStopwordRemoval, moduleObj.StopwordSet)
	newDocumentNameSet := util.Tokenize(requestPayload.NewDocumentName, moduleObj.IsUsingStopwordRemoval, moduleObj.StopwordSet)

	docType := getDocumentType(documentType, entity.Entity.GetDocumentTypes())

	// remove old document indexes
	isErrorRemoveExist := false
	errorRemoveKeys := ""

	for k, _ := range oldDocumentNameSet {
		key := fmt.Sprintf("%s:%s", docType, k)
		value := make([]interface{}, 1)
		value[0] = fmt.Sprintf("%d", documentID)
		_, err = moduleObj.Redis.SRem(key, value)
		if err != nil {
			isErrorRemoveExist = true
			errorRemoveKeys = errorRemoveKeys + " " + key + ","
			log.Println("[MODULE][UPDATE INDEX] failed to remove index on key :", key, "and document ID:", documentID)
			continue
		}
	}

	// add new document indexes
	isErrorAddExist := false
	errorAddKeys := ""

	for k, _ := range newDocumentNameSet {
		key := fmt.Sprintf("%s:%s", docType, k)
		value := make([]interface{}, 1)
		value[0] = fmt.Sprintf("%d", documentID)
		_, err = moduleObj.Redis.SAdd(key, value)
		if err != nil {
			isErrorAddExist = true
			errorAddKeys = errorAddKeys + " " + key + ","
			log.Println("[MODULE][UPDATE INDEX] failed to add index on key :", key, "and document ID:", documentID)
			continue
		}
	}

	if isErrorAddExist || isErrorRemoveExist {
		errorRemoveKeys = strings.TrimRight(errorRemoveKeys, ",")
		errorRemoveKeys = strings.TrimLeft(errorRemoveKeys, " ")

		errorAddKeys = strings.TrimRight(errorAddKeys, ",")
		errorAddKeys = strings.TrimLeft(errorAddKeys, " ")

		errorMessage := fmt.Sprintf("Error on removing following keys: %s and/or Error on adding following keys: %s", errorRemoveKeys, errorAddKeys)
		return Response{
			StatusCode:   http.StatusInternalServerError,
			ErrorMessage: errorMessage,
			Data:         nil,
		}
	}

	return Response{
		StatusCode:   http.StatusOK,
		ErrorMessage: "",
		Data:         nil,
	}
}
