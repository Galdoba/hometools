package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	JSON = "json"
)

type ConfigFile struct {
	app                 string
	serializationMethod string
	path                string
	RawData             []byte
}

func New(appName, serializationMethod string) (*ConfigFile, error) {
	cf := ConfigFile{}
	cf.app = appName
	cf.serializationMethod = serializationMethod
	path, err := configPath(cf.app)
	if err != nil {
		return nil, err
	}
	cf.path = path
	switch cf.serializationMethod {
	default:
		return nil, fmt.Errorf("method '%v' is not supported", serializationMethod)
	case JSON:
	}
	dir := filepath.Dir(cf.path)
	os.MkdirAll(dir, 0666)
	return &cf, nil
}

type Convertor interface {
	ToBytes() ([]byte, error)
	FromBytes([]byte) error
}

func (cf *ConfigFile) Write(cnv Convertor) error {
	bt, err := cnv.ToBytes()
	if err != nil {
		return fmt.Errorf("config to bytes conversion failed: %v", err)
	}
	path, err := configPath(cf.app)
	if err != nil {
		return fmt.Errorf("config path not detected: %v", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}
	if err = f.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate config file: %v", err)
	}
	if _, err = f.Write(bt); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}
	return f.Close()
}

func (cf *ConfigFile) Read() ([]byte, error) {
	bt, err := os.ReadFile(cf.path)
	return bt, err
}

func (cf *ConfigFile) Path() string {
	return cf.path
}

func (cf *ConfigFile) Dir() string {
	return filepath.Dir(cf.path) + string(filepath.Separator)
}

func configPath(appName string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	home = strings.ReplaceAll(home, `\`, "/")
	path := fmt.Sprintf("%v/.config/galdoba/%v/%v.config", home, appName, appName)

	return path, nil
}
