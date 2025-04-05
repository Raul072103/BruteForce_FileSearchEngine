package model

type SearchRequest struct {
	Words     []string `json:"words"`
	Extension []string `json:"extension"`
	Name      string   `json:"name"`
}
