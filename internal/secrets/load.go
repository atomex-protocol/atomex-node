package secrets

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
)

// Load - loads secret from environment variable. if it's not found checks docker secrets.
func Load(name string) string {
	if secret := os.Getenv(name); secret != "" {
		return secret
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("/run/secrets/%s", name))
	if err != nil {
		log.Err(err).Msg("ioutil.ReadFile")
		return ""
	}

	return string(data)
}
