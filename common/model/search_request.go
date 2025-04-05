package model

type SearchRequest struct {
	Words     map[string]any
	Extension map[string]any
	Name      string
}

type NetworkSearchRequest struct {
	Words     []string `json:"words"`
	Extension []string `json:"extension"`
	Name      string   `json:"name"`
}

// ConvertNetworkSearchRequest converts a NetworkSearchRequest to a SearchRequest
func ConvertNetworkSearchRequest(nsr NetworkSearchRequest) SearchRequest {
	sr := SearchRequest{
		Words:     make(map[string]any),
		Extension: make(map[string]any),
		Name:      nsr.Name,
	}

	for _, word := range nsr.Words {
		sr.Words[word] = struct{}{}
	}

	for _, ext := range nsr.Extension {
		sr.Extension[ext] = struct{}{}
	}

	return sr
}
