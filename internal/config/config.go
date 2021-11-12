package config

import (
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

// Load -
func Load(filename string, output interface{}) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(output)
	return err
}

// SelectEnvironment -
func SelectEnvironment(dir string) string {
	env := os.Getenv("AP_ENV")
	switch env {
	case "production":
		return path.Join(dir, "production")
	case "test":
		return path.Join(dir, "test")
	default:
		return dir
	}
}
