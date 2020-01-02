package entity

import ()

//DocumentType is a type that represent document type
type DocumentType string

const (
	//AdvertisementCampaignDocument is the document type that represent Airy's Advertisement Campaign document type
	AdvertisementCampaignDocument DocumentType = "advcampaign"
	//CampaignDocument is the document type that represent Airy's Campaign document type (coupon for promotion)
	CampaignDocument DocumentType = "campaign"
)
