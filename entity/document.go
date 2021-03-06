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
)

//DocumentType is a type that represent document type
type DocumentType string

const (
	//These document types below are just for example.
	//You can create your own document type and don't forget to modify the module/document.go for document type validation.

	//AdvertisementCampaignDocument is the document type that represent Advertisement Campaign document type
	AdvertisementCampaignDocument DocumentType = "advcampaign"
	//CampaignDocument is the document type that represent Campaign document type (coupon for promotion)
	CampaignDocument DocumentType = "campaign"
)

/*
//IsValid checks if the document type is a valid (registered) document type in entity const
func (dt DocumentType) IsValid() error {
	switch dt {
	case AdvertisementCampaignDocument, CampaignDocument:
		return nil
	}
	return errors.New("Invalid Document Type")
}
*/
//IsValidFromCustomType checks if the document type is a valid (registered) in a documentTypeMap (Custom Document Type)
func (dt DocumentType) IsValidFromCustomDocumentType(documentTypeMap map[DocumentType]int) error {
	if _, ok := documentTypeMap[dt]; ok {
		return nil
	}
	return errors.New("Invalid Document Type")
}
