package resources

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"dir-agent/internal/config"
)

//go:embed assets/**
var embedded embed.FS

type InstallResult struct {
	DataPath   string `json:"data_path"`
	ConfigPath string `json:"config_path"`
}

func Install() (InstallResult, error) {
	dataPath, err := config.DataPath()
	if err != nil {
		return InstallResult{}, err
	}
	configPath, err := config.EnsureConfigFile()
	if err != nil {
		return InstallResult{}, err
	}

	if err := os.MkdirAll(dataPath, 0o755); err != nil {
		return InstallResult{}, err
	}

	err = fs.WalkDir(embedded, "assets", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		targetPath := filepath.Join(dataPath, path)
		if entry.IsDir() {
			return os.MkdirAll(targetPath, 0o755)
		}
		return copyEmbeddedFile(path, targetPath)
	})
	if err != nil {
		return InstallResult{}, err
	}

	return InstallResult{
		DataPath:   dataPath,
		ConfigPath: configPath,
	}, nil
}

func Uninstall(removeConfig bool) error {
	dataPath, err := config.DataPath()
	if err != nil {
		return err
	}
	if err := os.RemoveAll(filepath.Join(dataPath, "assets")); err != nil {
		return err
	}

	if removeConfig {
		configPath, err := config.ConfigPath()
		if err != nil {
			return err
		}
		if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func copyEmbeddedFile(sourcePath string, targetPath string) error {
	sourceFile, err := embedded.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open embedded file %s: %w", sourcePath, err)
	}
	defer sourceFile.Close()

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return err
	}

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, sourceFile)
	return err
}
