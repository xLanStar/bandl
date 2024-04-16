package bandl

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type IConfig interface{}

type Config struct{}

func InitConfig(path string, out IConfig) {
	if _, err := os.Stat(path); err == nil {
		LoadConfig(path, out)
	} else {
		SaveConfig(path, out)
	}
}

func LoadConfig(path string, out IConfig) error {
	log.Printf("Load from %s\n", path)

	if _, err := os.Stat(path); err != nil {
		return err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, out)
	if err != nil {
		return err
	}

	return nil
}

func SaveConfig(path string, config IConfig) error {
	log.Printf("Save to %s\n", path)

	file, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.WriteString(f, string(file))
	if err != nil {
		return err
	}

	return nil
}
