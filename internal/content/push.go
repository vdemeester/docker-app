package content

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// Push will push the specified app from the store
func (s *store) Push(ctx context.Context, reference string) (string, error) {
	resolv := docker.NewResolver(docker.ResolverOptions{
		Credentials: func(host string) (string, string, error) {
			// Only one host
			return "vdemeester", "docEC08eYpdker", nil
		},
	})

	ref, err := referenceSum(reference)
	if err != nil {
		return "", err
	}

	r, err := os.Open(filepath.Join(s.path, ref.Sum()))
	if err != nil {
		return "", err
	}
	defer r.Close()
	ctx = namespaces.WithNamespace(ctx, "app")

	img, err := s.imageStore.Get(ctx, ref.String())
	if err != nil {
		return "", err
	}
	fmt.Println(img.Target)
	fmt.Println(img.Target.MediaType)
	fmt.Println(img.Target.Digest)

	p, err := resolv.Pusher(ctx, img.Name)
	if err != nil {
		return "", err
	}
	fmt.Println("p", p)

	/*
		w, err := p.Push(ctx, img.Target)
		if err != nil {
			return "", err
		}
		defer w.Close()
	*/

	var m sync.Mutex
	manifestStack := []ocispec.Descriptor{}
	pushHandler := remotes.PushHandler(p, s.inner)

	filterHandler := images.HandlerFunc(func(ctx context.Context, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) {
		fmt.Println("filter", desc)
		switch desc.MediaType {
		case images.MediaTypeDockerSchema2Manifest, ocispec.MediaTypeImageManifest,
			images.MediaTypeDockerSchema2ManifestList, ocispec.MediaTypeImageIndex:
			m.Lock()
			manifestStack = append(manifestStack, desc)
			m.Unlock()
			return nil, images.ErrStopHandler
		default:
			fmt.Println("not a manifest")
			return nil, nil
		}
	})

	h := images.Handlers(
		images.ChildrenHandler(s.inner),
		filterHandler,
		pushHandler,
	)

	err = images.Dispatch(ctx, h, img.Target)
	if err != nil {
		return "", err
	}
	fmt.Println(img)

	// Iterate in reverse order as seen, parent always uploaded after child
	for i := len(manifestStack) - 1; i >= 0; i-- {
		_, err := pushHandler(ctx, manifestStack[i])
		if err != nil {
			// TODO(estesp): until we have a more complete method for index push, we need to report
			// missing dependencies in an index/manifest list by sensing the "400 Bad Request"
			// as a marker for this problem
			if (manifestStack[i].MediaType == ocispec.MediaTypeImageIndex ||
				manifestStack[i].MediaType == images.MediaTypeDockerSchema2ManifestList) &&
				errors.Cause(err) != nil && strings.Contains(errors.Cause(err).Error(), "400 Bad Request") {
				return "", errors.Wrap(err, "manifest list/index references to blobs and/or manifests are missing in your target registry")
			}
			return "", err
		}
	}
	fmt.Println("foo")

	return "", err
}
