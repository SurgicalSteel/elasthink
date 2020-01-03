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
import ()

//DocumentType is a type that represent document type
type DocumentType string

const (
	//AdvertisementCampaignDocument is the document type that represent Airy's Advertisement Campaign document type
	AdvertisementCampaignDocument DocumentType = "advcampaign"
	//CampaignDocument is the document type that represent Airy's Campaign document type (coupon for promotion)
	CampaignDocument DocumentType = "campaign"
)
