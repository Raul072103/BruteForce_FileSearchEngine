package model

type DirectoryResponse struct {
	Path          string               `utils:"path"`
	SearchRequest NetworkSearchRequest `utils:"search_request"`
}
