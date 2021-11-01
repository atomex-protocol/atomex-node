package config

import (
	"os"

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
