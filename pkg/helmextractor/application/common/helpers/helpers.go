package helpers

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// GetEnvVar retrieves the environment variable with the supplied name and fails
// if it is not able to do so
func GetEnvVar(envVarName string) (string, error) {
	envVar := os.Getenv(envVarName)
	if envVar == "" {
		return "", fmt.Errorf("the environment variable '%s' is not set", envVarName)
	}

	return envVar, nil
}

// MustGetEnvVar returns the value of the environment variable with the indicated name or panics.
// It is intended to be used in the INTERNALS of the server when we can guarantee (through orderly
// coding) that the environment variable was set at server startup.
// Since the env is required, kill the app if the env is not set. In the event a variable is not super
// required, set a sensible default or don't call this method
func MustGetEnvVar(envVarName string) string {
	val, err := GetEnvVar(envVarName)
	if err != nil {
		log.Panicf("mandatory environment variable %s not found", envVarName)
		os.Exit(1)
	}

	return val
}
