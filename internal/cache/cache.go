package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/wolv89/troster/internal/models"
)

const (
	cacheDir = "data"
	maxAge   = 24 * time.Hour
)

func cacheFile(compValue string) string {
	return filepath.Join(cacheDir, fmt.Sprintf("fixtures_%s.json", compValue))
}

func Load(compValue string) (*models.ScrapedData, error) {
	path := cacheFile(compValue)

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if time.Since(info.ModTime()) > maxAge {
		return nil, fmt.Errorf("cache expired")
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var data models.ScrapedData
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func Save(compValue string, data *models.ScrapedData) error {
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return err
	}

	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFile(compValue), raw, 0o644)
}
