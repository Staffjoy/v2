package environments

import (
	"fmt"
	"os"
	"strings"
)

const (
	// GcloudEnvVar is the environment variable for accessing the GCloud secret key
	GcloudEnvVar = "GCLOUD"
	// GcloudProjectEnvVar is the environment variable for accessign the Gcloud project name
	GcloudProjectEnvVar = "GCLOUD_PROJECT"
)

// GetPublicSentryDSN returns a Sentry Id that can be returned in JS
func GetPublicSentryDSN(secret string) (string, error) {
	// Secret DSN looks like this:
	// https://username:password@sentry.io/id
	//
	// Public DSN removes the password part to be:
	// https://username@sentry.io/id

	const prefix = "https://"
	if secret == "" {
		return "", nil
	}

	// Find the part we need to cut out. We cut off the prefix.
	startIndex := strings.Index(secret[len(prefix):], ":") + len(prefix)
	endIndex := strings.Index(secret, "@")

	if (startIndex == -1) || (endIndex == -1) || (endIndex < startIndex) {
		return "", fmt.Errorf("Unable to determine error")
	}
	return secret[0:startIndex] + secret[endIndex:], nil
}

// GetGoogleCloudProject returns the identifier of the google cloud project8
func GetGoogleCloudProject() string {
	return os.Getenv(GcloudProjectEnvVar)
}
