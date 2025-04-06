package model

type FileMetadata struct {
	Path      string `json:"path"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	Extension string `json:"extension"`
	Preview   string `json:"preview"`
	Content   []byte
}

// ConvertToResponse converts FileMetadata into FileSearchResponse with a given preview
func ConvertToResponse(fm *FileMetadata, preview string) FileSearchResponse {
	return FileSearchResponse{
		Path:      fm.Path,
		Name:      fm.Name,
		Size:      fm.Size,
		Extension: fm.Extension,
		Preview:   preview,
	}
}
