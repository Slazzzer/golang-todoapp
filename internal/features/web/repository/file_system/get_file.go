package web_fs_repository

import (
	"errors"
	"fmt"
	"os"

	core_errors "github.com/Slazzzer/golang-todoapp/internal/core/errors"
)

func (r *WebFSRepository) GetFile(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf(
				"file: %s: %w",
				filePath,
				core_errors.ErrNotFound,
			)
		}
		return nil, fmt.Errorf(
			"get_file: %s: %w",
			filePath,
			err,
		)
	}
	return file, nil
}
