package env

import (
	"log"
	"os"
)

func MustEnv(key string) string {

	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("[ERR] enviromental variable \"%s\" don't specified ", key)
	}

	return value
}
