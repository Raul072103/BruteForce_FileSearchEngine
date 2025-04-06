package model

type SearchRequest struct {
	Words     map[string]any
	Extension map[string]any
	Name      string
}

type NetworkSearchRequest struct {
	Words     []string `utils:"words"`
	Extension []string `utils:"extension"`
	Name      string   `utils:"name"`
}

type FileSearchResponse struct {
	Path      string `utils:"path"`
	Name      string `utils:"name"`
	Size      int64  `utils:"size"`
	Extension string `utils:"extension"`
	Preview   string `utils:"preview"`
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

// ConvertSearchRequest converts a SearchRequest to a NetworkSearchRequest
func ConvertSearchRequest(sr SearchRequest) NetworkSearchRequest {
	nsr := NetworkSearchRequest{
		Words:     make([]string, 0, len(sr.Words)),
		Extension: make([]string, 0, len(sr.Extension)),
		Name:      sr.Name,
	}

	// Extract keys from Words map
	for word := range sr.Words {
		nsr.Words = append(nsr.Words, word)
	}

	// Extract keys from Extension map
	for ext := range sr.Extension {
		nsr.Extension = append(nsr.Extension, ext)
	}

	return nsr
}
