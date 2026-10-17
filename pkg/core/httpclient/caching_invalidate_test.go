package httpclient

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCachingTransport_InvalidatesCacheOnNonGet(t *testing.T) {
	// Arrange
	srv, hits := newCountingServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	ct := newCachingTransport(http.DefaultTransport)
	client := &http.Client{Transport: ct}

	// Act 1: GET populates the cache
	resp1, err := client.Get(srv.URL + "/resource")
	require.NoError(t, err)
	io.ReadAll(resp1.Body) //nolint:errcheck
	resp1.Body.Close()

	// Act 2: POST to same URL or any URL invalidates the cache? Wait, if we want to invalidate on non-GET requests to prevent stale cache hits...
	resp2, err := client.Post(srv.URL+"/submit", "application/json", nil)
	require.NoError(t, err)
	io.ReadAll(resp2.Body) //nolint:errcheck
	resp2.Body.Close()

	// Act 3: GET again, should miss cache if invalidated
	resp3, err := client.Get(srv.URL + "/resource")
	require.NoError(t, err)
	io.ReadAll(resp3.Body) //nolint:errcheck
	resp3.Body.Close()

	// Assert
	assert.Equal(t, int64(3), hits.Load(), "cache should be invalidated, causing 3 hits")
}
