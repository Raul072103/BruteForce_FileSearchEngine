package model

type FileMetadata struct {
	Path      string `json:"path"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	Extension string `json:"extension"`
	Preview   string `json:"preview"`
	Content   []byte
}
