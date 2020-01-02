package module

import (
	"errors"
	"github.com/SurgicalSteel/elasthink/entity"
	"strings"
)

func validateDocumentType(docType string) error {
	docType = strings.ToLower(docType)
	documentType := entity.DocumentType(docType)
	switch documentType {
	case
		entity.AdvertisementCampaignDocument,
		entity.CampaignDocument:
		return nil
	default:
		return errors.New("Invalid Document Type")
	}
}

func getDocumentType(docType string) entity.DocumentType {
	docType = strings.ToLower(docType)
	documentType := entity.DocumentType(docType)
	switch documentType {
	case
		entity.AdvertisementCampaignDocument,
		entity.CampaignDocument:
		return documentType
	default:
		return ""
	}
}
