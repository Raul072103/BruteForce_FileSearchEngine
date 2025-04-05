package repo

import (
	"BruteForce_SearchEnginer/common/model"
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrFileNotFound = errors.New("the given path doesn't exist")
)

type FileRepo interface {
	Read(path string) (*model.FileMetadata, error)
	Stats(path string) (*model.FileMetadata, error)
}

type fileRepo struct {
	typeMap model.FileTypesConfig
}

func New(typeMap model.FileTypesConfig) FileRepo {
	return &fileRepo{typeMap: typeMap}
}

// ReadFile reads the contents of the file, if it exists, at the given path and returns the content of that file.
func (fr *fileRepo) Read(path string) (*model.FileMetadata, error) {
	file, err := fr.Stats(path)
	if err != nil {
		return file, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if fr.typeMap.GetTypeByExtension(file.Extension) == "txt" {
		file.Content = data
	}

	return file, nil
}

// Stats if the file exists it returns an instance of models.File, else a nil.
func (fr *fileRepo) Stats(path string) (*model.FileMetadata, error) {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, ErrFileNotFound
	}

	if err != nil {
		return nil, err
	}

	file := model.FileMetadata{
		Path:      path,
		Name:      fileInfo.Name(),
		Size:      fileInfo.Size(),
		Extension: getFileExtension(path, fileInfo.IsDir()),
	}

	return &file, err
}

func getFileExtension(path string, isDir bool) string {
	if isDir {
		return ""
	}
	return filepath.Ext(path)
}
