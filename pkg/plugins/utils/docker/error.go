package docker

import (
	"strings"

	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
)

// IsAuthError returns true if the error is due to an authentication issue
func IsAuthError(err error) bool {
	if err == nil {
		return false
	}
	if terr, ok := err.(*transport.Error); ok {
		if terr.StatusCode == 401 || terr.StatusCode == 403 {
			return true
		}
	}
	if strings.Contains(err.Error(), "DENIED: Not Authorized") || strings.Contains(err.Error(), "error getting credentials") {
		return true
	}
	return false
}
