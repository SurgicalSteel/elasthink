package entity

//SearchResultRankData is the core struct that represent search result rank item
type SearchResultRankData struct {
	ID        int64 `json:"id"`
	ShowCount int   `json:"showCount"`
	Rank      int   `json:"rank"`
}
