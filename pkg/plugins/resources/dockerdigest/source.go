package dockerdigest

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/sirupsen/logrus"
	"github.com/updatecli/updatecli/pkg/core/result"
	"github.com/updatecli/updatecli/pkg/plugins/utils/docker"
)

// Source retrieves Docker image tag digest from a registry
func (ds *DockerDigest) Source(_ context.Context, workingDir string, resultSource *result.Source) error {
	refTag := "latest"
	refName := ds.spec.Image

	if ds.spec.Tag != "" {
		if strings.HasPrefix(ds.spec.Tag, "@") {
			return fmt.Errorf("invalid tag %s: only contain a digest", ds.spec.Image)
		}

		refTagArray := strings.Split(ds.spec.Tag, "@")
		refTag = strings.TrimPrefix(refTagArray[0], ":")

		refName += ":" + refTag
	}

	ref, err := name.ParseReference(refName)
	if err != nil {
		return fmt.Errorf("invalid image %s: %w", refName, err)
	}

	opts := append(ds.options, remote.WithAuthFromKeychain(ds.keychain))
	remoteDescriptor, err := remote.Get(ref, opts...)
	if err != nil {
		if docker.IsAuthError(err) {
			logrus.Debugf("unable to retrieve image %s with authentication, falling back to anonymous: %s", refName, err)
			anonOpts := append(ds.options, remote.WithAuth(authn.Anonymous))
			remoteDescriptor, err = remote.Get(ref, anonOpts...)
			opts = anonOpts // update opts for subsequent calls like Image()
		}
		if err != nil {
			return fmt.Errorf("unable to retrieve image %s: %w", refName, err)
		}
	}

	digest := remoteDescriptor.Digest
	if ds.spec.Architecture != "" {
		image, err := remote.Image(ref, opts...)
		if err != nil {
			return fmt.Errorf("unable to retrieve image %s: %w", refName, err)
		}

		digest, err = image.Digest()
		if err != nil {
			return fmt.Errorf("unable to retrieve image digest %s: %w", refName, err)
		}
	}

	finalDigest := refTag + "@" + digest.String()
	imageDigest := ref.Context().Name() + ":" + finalDigest

	if ds.spec.HideTag {
		imageDigest = ref.Context().Name() + "@" + digest.String()
		finalDigest = "@" + digest.String()
	}

	resultSource.Result = result.SUCCESS
	resultSource.Information = finalDigest
	resultSource.Description = fmt.Sprintf("Docker Image Tag %s resolved to digest %s",
		ref.String(), imageDigest)

	return nil
}
