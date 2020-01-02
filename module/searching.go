package module

import (
	"context"
	"errors"
	"github.com/SurgicalSteel/elasthink/entity"
	"net/http"
	"strings"
)

//SearchRequestPayload is the universal request payload for search handlers
type SearchRequestPayload struct {
	SearchTerm string `json:"searchTerm"`
}

func validateSearchRequestPayload(documentType, searchTerm string) error {
	if len(strings.Trim(searchTerm, " ")) == 0 {
		return errors.New("Search Term is required")
	}

	if len(strings.Trim(documentType, " ")) == 0 {
		return errors.New("Document Type is required")
	}

	err := validateDocumentType(documentType)
	if err != nil {
		return err
	}

	return nil
}

//SearchResponsePayload is the universal response payload for search handlers
type SearchResponsePayload struct {
	RankedResultList []entity.SearchResultRankData `json:"rankedResultList"`
}

//Search is the core function of searching a document
func Search(ctx context.Context, documentType string, requestPayload SearchRequestPayload) Response {
	err := validateSearchRequestPayload(documentType, requestPayload.SearchTerm)
	if err != nil {
		return Response{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: err.Error(),
			Data:         nil,
		}
	}

	searchTermSet := tokenizeIndonesianSearchTerm(requestPayload.SearchTerm)
	if len(searchTermSet) == 0 {
		return Response{
			StatusCode:   http.StatusOK,
			ErrorMessage: "",
			Data:         nil,
		}
	}

	docType := getDocumentType(documentType)

	wordIndexSets := fetchWordIndexSets(docType, searchTermSet)

	if len(wordIndexSets) == 0 {
		return Response{
			StatusCode:   http.StatusOK,
			ErrorMessage: "",
			Data:         nil,
		}
	}

	rankedSearchResult := rankSearchResult(wordIndexSets)
	searchResponsePayload := SearchResponsePayload{RankedResultList: rankedSearchResult}

	return Response{
		StatusCode:   http.StatusOK,
		ErrorMessage: "",
		Data:         searchResponsePayload,
	}
}
