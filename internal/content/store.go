package content

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/content/local"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/images/oci"
	"github.com/containerd/containerd/metadata"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/snapshots"
	appMetadata "github.com/docker/app/types/metadata"
	"github.com/docker/docker/errdefs"
)

type Store interface {
	// Pull will pull the specified app into the store
	Pull(ctx context.Context, reference string) error
	// Get returns a tar reader of the app
	Get(ctx context.Context, reference string) (io.Reader, error)
	// Push will push the specified app from the store
	Push(ctx context.Context, reference string) (string, error)
	// Load will load the specified reader (tar) into the store
	Load(ctx context.Context, reference string, r io.Reader) error
	// List lists the docker app available in the store
	List(ctx context.Context) ([]appMetadata.AppMetadata, error)
}

func NewStore(path string) (Store, error) {
	for _, subpath := range []string{"ctrd", "apps"} {
		if err := os.MkdirAll(filepath.Join(path, subpath), 0744); err != nil {
			return nil, err
		}
	}
	db, err := bolt.Open(filepath.Join(path, "ctrdapp.db"), 0644, nil)
	if err != nil {
		return nil, err
	}
	cs, err := local.NewStore(filepath.Join(path, "content"))
	if err != nil {
		return nil, err
	}
	mdb := metadata.NewDB(db, cs, map[string]snapshots.Snapshotter{})
	if err := mdb.Init(context.TODO()); err != nil {
		return nil, err
	}
	is := metadata.NewImageStore(mdb)
	return &store{
		inner:      mdb.ContentStore(),
		imageStore: is,
		path:       filepath.Join(path, "apps"),
	}, nil
}

type store struct {
	inner      content.Store
	imageStore images.Store
	path       string
}

// Get returns a tar reader of the app
func (s *store) Get(ctx context.Context, reference string) (io.Reader, error) {
	ref, err := referenceSum(reference)
	if err != nil {
		return nil, err
	}
	return os.Open(filepath.Join(s.path, ref.Sum()))
}

// Load will load the specified reader (tar) into the store
func (s *store) Load(ctx context.Context, reference string, r io.Reader) error {
	ctx = namespaces.WithNamespace(ctx, "app")
	ref, err := referenceSum(reference)
	if err != nil {
		return err
	}

	d, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(s.path, ref.Sum()), d, 0755); err != nil {
		return err
	}
	f, err := os.Open(filepath.Join(s.path, ref.Sum()))
	if err != nil {
		return err
	}
	defer f.Close()

	ociReader, err := createOCIImage(f)
	if err != nil {
		return err
	}
	importer := &oci.V1Importer{
		ImageName: ref.String(),
	}
	imgrecs, err := importer.Import(ctx, s.inner, ociReader)
	if err != nil {
		return err
	}
	for _, img := range imgrecs {
		img.Name = ref.String()
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
	}
	return nil
}

// List lists the docker app available in the store
func (s *store) List(ctx context.Context) ([]appMetadata.AppMetadata, error) {
	ctx = namespaces.WithNamespace(ctx, "app")
	images, err := s.imageStore.List(ctx)
	if err != nil {
		return nil, err
	}
	apps := make([]appMetadata.AppMetadata, len(images))
	for i, image := range images {
		apps[i] = appMetadata.AppMetadata{
			Name: image.Name,
		}
	}
	return apps, nil
}
