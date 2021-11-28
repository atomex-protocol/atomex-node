package config

import (
	"context"
	"os"
	"path"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

var (
	valid = validator.New()
)

// Load -
func Load(ctx context.Context, filename string, output interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(output); err != nil {
		return err
	}

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	val := reflect.ValueOf(output)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		return valid.StructCtx(cancelCtx, output)
	case reflect.Array:
		return valid.VarCtx(cancelCtx, output, "dive")
	}
	return nil
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
