package helpers

import "encoding/base64"

func BasicAuth64(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
