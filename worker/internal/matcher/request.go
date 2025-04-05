package matcher

import (
	"BruteForce_SearchEnginer/common/model"
	"strings"
)

type RequestMatcher interface {
	MatchFile(fileMetadata model.FileMetadata, request model.SearchRequest) bool
}

type requestMatcher struct {
	typeMap model.FileTypesConfig
}

func New(typeMap model.FileTypesConfig) RequestMatcher {
	return &requestMatcher{typeMap: typeMap}
}

func (r *requestMatcher) MatchFile(fileMetadata model.FileMetadata, request model.SearchRequest) bool {
	var wordSearchRequest, fileTxtType bool

	if len(request.Words) != 0 {
		wordSearchRequest = true
	} else {
		wordSearchRequest = false
	}

	if r.typeMap.GetTypeByExtension(fileMetadata.Extension) == ".txt" {
		fileTxtType = true
	} else {
		fileTxtType = false
	}

	// if we have a request for finding words, and it is not a .txt file than this request automatically fails.
	if wordSearchRequest && fileTxtType == false {
		return false
	}

	// extension search
	if len(request.Extension) > 0 {
		if _, exists := request.Extension[fileMetadata.Extension]; exists == false {
			return false
		}
	}

	// name search
	if request.Name == "" {
		if strings.HasPrefix(fileMetadata.Name, request.Name) == false {
			return false
		}
	}

	// words search
	if wordSearchRequest {
		words := strings.Split(string(fileMetadata.Content), "")

		var matchedWords = 0

		for _, word := range words {
			if _, exists := request.Words[word]; exists {
				matchedWords += 1
			}
		}

		if matchedWords < len(request.Words) {
			return false
		}
	}

	return true
}
