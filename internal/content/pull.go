package content

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	"github.com/docker/docker/errdefs"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// Pull will pull the specified app into the store
func (s *store) Pull(ctx context.Context, reference string) error {
	resolv := docker.NewResolver(docker.ResolverOptions{})

	ref, err := referenceSum(reference)
	if err != nil {
		return err
	}

	ctx = namespaces.WithNamespace(ctx, "app")
	name, desc, err := resolv.Resolve(ctx, ref.String())
	if err != nil {
		return err
	}
	f, err := resolv.Fetcher(ctx, name)
	if err != nil {
		return err
	}

	r, err := f.Fetch(ctx, desc)
	if err != nil {
		return err
	}
	defer r.Close()

	h := images.Handlers(
		remotes.FetchHandler(s.inner, f),
		configHandler(s.inner, filepath.Join(s.path, ref.Sum())),
	)
	if err := images.Dispatch(ctx, h, desc); err != nil {
		return err
	}
	img := images.Image{
		Name:   ref.String(),
		Target: desc,
	}
	_, err = s.imageStore.Create(ctx, img)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return err
		}
		_, err = s.imageStore.Update(ctx, img)
		if err != nil {
			// if image was removed, try create again
			if errdefs.IsNotFound(err) {
				return nil
			}
			return err
		}
	}
	return nil
}

func configHandler(provider content.Provider, path string) images.HandlerFunc {
	return func(ctx context.Context, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) {
		var descs []ocispec.Descriptor
		switch desc.MediaType {
		case images.MediaTypeDockerSchema2Manifest, ocispec.MediaTypeImageManifest:
			p, err := content.ReadBlob(ctx, provider, desc.Digest)
			if err != nil {
				return nil, err
			}
			// TODO(stevvooe): We just assume oci manifest, for now. There may be
			// subtle differences from the docker version.
			var manifest ocispec.Manifest
			if err := json.Unmarshal(p, &manifest); err != nil {
				return nil, err
			}
			if len(manifest.Layers) > 1 {
				return nil, errors.New("not a valid docker app image")
			}
			descs = append(descs, manifest.Layers...)
		case images.MediaTypeDockerSchema2Layer, images.MediaTypeDockerSchema2LayerGzip,
			ocispec.MediaTypeImageLayer, ocispec.MediaTypeImageLayerGzip:
			p, err := content.ReadBlob(ctx, provider, desc.Digest)
			if err != nil {
				return nil, err
			}
			if err := ioutil.WriteFile(path, p, 0755); err != nil {
				return nil, err
			}
		default:
			// Do nothing, it's neither a manifest nor the config layer
		}
		return descs, nil
	}
}
