package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Config — persistent application settings.
type Config struct {
	DataDir   string `json:"data_dir,omitempty"`
	BackupDir string `json:"backup_dir,omitempty"`
}

var (
	folderMu     sync.RWMutex
	backupFolder string
)

// ConfigPath returns the platform-specific path to the config file.
func ConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "LongDarkBackup", "config.json"), nil
}

// LoadConfig reads the config file. A missing file is not an error.
func LoadConfig() (Config, error) {
	var cfg Config

	path, err := ConfigPath()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

// SaveConfig writes the config file atomically, creating its directory if needed.
func SaveConfig(cfg Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	tmp, err := os.CreateTemp(dir, "config-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, 0644); err != nil {
		return err
	}

	return os.Rename(tmpName, path)
}

// BackupBrokenConfig renames an unreadable config file so it is not overwritten.
func BackupBrokenConfig() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	backupPath := path + ".bak"
	// os.Rename does not replace existing files on Windows
	os.Remove(backupPath)
	return os.Rename(path, backupPath)
}

// ResolveFolder picks the effective directory.
// Priority: config file > environment variable > default.
func ResolveFolder(configValue, envValue, defaultValue string) string {
	switch {
	case configValue != "":
		return configValue
	case envValue != "":
		return envValue
	default:
		return defaultValue
	}
}

// GetBackupFolder returns the current backup directory.
func GetBackupFolder() string {
	folderMu.RLock()
	defer folderMu.RUnlock()
	return backupFolder
}

// SetBackupFolder changes the backup directory at runtime.
func SetBackupFolder(dir string) {
	folderMu.Lock()
	backupFolder = dir
	folderMu.Unlock()
}
