package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0o755); err != nil {
		return nil, fmt.Errorf("create storage dir: %w", err)
	}
	return &LocalStorage{basePath: basePath}, nil
}

func (s *LocalStorage) Save(userID uuid.UUID, originalName string, r io.Reader) (storedName, urlPath string, err error) {
	ext := strings.ToLower(filepath.Ext(originalName))
	if ext != ".pdf" {
		return "", "", fmt.Errorf("only pdf files are allowed")
	}

	storedName = uuid.New().String() + ext
	userDir := filepath.Join(s.basePath, userID.String())
	if err := os.MkdirAll(userDir, 0o755); err != nil {
		return "", "", fmt.Errorf("create user dir: %w", err)
	}

	destPath := filepath.Join(userDir, storedName)
	f, err := os.Create(destPath)
	if err != nil {
		return "", "", fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		_ = os.Remove(destPath)
		return "", "", fmt.Errorf("write file: %w", err)
	}

	urlPath = filepath.ToSlash(filepath.Join("/storage", userID.String(), storedName))
	return storedName, urlPath, nil
}

func (s *LocalStorage) Open(userID uuid.UUID, storedName string) (string, error) {
	path := filepath.Join(s.basePath, userID.String(), storedName)
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("file not found")
	}
	return path, nil
}

func (s *LocalStorage) BasePath() string { return s.basePath }
