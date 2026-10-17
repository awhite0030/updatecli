package lock

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/updatecli/updatecli/pkg/core/httpclient"
)

func TestGetProviderHashesOpenTofu(t *testing.T) {
	mockResponse := `{
		"packages": {
			"linux_amd64": {
				"hashes": [
					"h1:mock-h1-hash1",
					"zh:mock-zh-hash1"
				]
			},
			"darwin_arm64": {
				"hashes": [
					"h1:mock-h1-hash2",
					"zh:mock-zh-hash2"
				]
			}
		}
	}`

	mockClient := &httpclient.MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "https://registry.opentofu.org/v1/providers/hashicorp/aws/5.10.0/download/linux/amd64", req.URL.String())
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
			}, nil
		},
	}

	spec := Spec{
		File:      "testdata/terraform.lock.hcl",
		Provider:  "registry.opentofu.org/hashicorp/aws",
		Platforms: []string{"linux_amd64"},
	}

	resource, err := New(spec)
	require.NoError(t, err)

	// Inject the mock client
	resource.httpClient = mockClient

	hashes, err := resource.getProviderHashes("5.10.0")
	require.NoError(t, err)

	expectedHashes := []string{
		"h1:mock-h1-hash1",
		"h1:mock-h1-hash2",
		"zh:mock-zh-hash1",
		"zh:mock-zh-hash2",
	}

	assert.ElementsMatch(t, expectedHashes, hashes)
}
