package chain

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

// LoadSecret - loads secret from environment variable. if it's not found checks docker secrets.
func LoadSecret(name string) (string, error) {
	if secret := os.Getenv(name); secret != "" {
		return secret, nil
	}

	f, err := os.Open("/run/secrets/watch-tower")
	if err != nil {
		return "", err
	}
	defer f.Close()

	var secrets map[string]string
	if err := json.NewDecoder(f).Decode(&secrets); err != nil {
		return "", err
	}

	if secret, ok := secrets[name]; ok {
		return secret, nil
	}

	return "", errors.Errorf("unknown secret: %s", name)
}
