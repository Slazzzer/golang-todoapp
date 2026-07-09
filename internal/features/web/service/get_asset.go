package web_service

import (
	"fmt"
	"os"
	"path"

	core_errors "github.com/Slazzzer/golang-todoapp/internal/core/errors"
)

// allowedAssets ограничивает раздачу статикой фронтенда и задаёт Content-Type.
var allowedAssets = map[string]string{
	"styles.css": "text/css; charset=utf-8",
	"app.js":     "application/javascript; charset=utf-8",
}

func (s *WebService) GetAsset(name string) ([]byte, string, error) {
	contentType, ok := allowedAssets[name]
	if !ok {
		return nil, "", fmt.Errorf("asset %q not allowed: %w", name, core_errors.ErrNotFound)
	}

	assetPath := path.Join(
		os.Getenv("PROJECT_ROOT"),
		"public",
		name,
	)

	content, err := s.webRepository.GetFile(assetPath)
	if err != nil {
		return nil, "", fmt.Errorf("get_asset from repository: %w", err)
	}

	return content, contentType, nil
}
