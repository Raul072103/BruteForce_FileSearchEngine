package model

type FileMetadata struct {
	Path      string `utils:"path"`
	Name      string `utils:"name"`
	Size      int64  `utils:"size"`
	Extension string `utils:"extension"`
	Preview   string `utils:"preview"`
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
