package syncdaemon

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/deepgram/dx-asana/internal/asana"
)

type CacheMetadata struct {
	ProjectID string    `json:"project_id"`
	SyncedAt  time.Time `json:"synced_at"`
	TaskCount int       `json:"task_count"`
}

type CacheFile struct {
	Metadata CacheMetadata `json:"metadata"`
	Tasks    []asana.Task  `json:"tasks"`
}

func cacheTasks(projectID string, tasks []asana.Task) error {
	cacheDir := getCacheDir()
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}

	cacheFile := CacheFile{
		Metadata: CacheMetadata{
			ProjectID: projectID,
			SyncedAt:  time.Now(),
			TaskCount: len(tasks),
		},
		Tasks: tasks,
	}

	filePath := filepath.Join(cacheDir, fmt.Sprintf("project-%s.json", projectID))
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(cacheFile)
}

func LoadCachedTasks(projectID string) ([]asana.Task, *CacheMetadata, error) {
	filePath := filepath.Join(getCacheDir(), fmt.Sprintf("project-%s.json", projectID))
	
	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("cache not found: %w", err)
	}
	defer f.Close()

	var cacheFile CacheFile
	if err := json.NewDecoder(f).Decode(&cacheFile); err != nil {
		return nil, nil, fmt.Errorf("failed to decode cache: %w", err)
	}

	return cacheFile.Tasks, &cacheFile.Metadata, nil
}

func ClearCache(projectID string) error {
	filePath := filepath.Join(getCacheDir(), fmt.Sprintf("project-%s.json", projectID))
	return os.Remove(filePath)
}

func GetCacheSize() (int64, error) {
	var totalSize int64
	entries, err := os.ReadDir(getCacheDir())
	if err != nil {
		return 0, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			info, _ := entry.Info()
			totalSize += info.Size()
		}
	}

	return totalSize, nil
}