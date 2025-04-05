package model

type DirectoryResponse struct {
	Path          string               `json:"path"`
	SearchRequest NetworkSearchRequest `json:"search_request"`
}
