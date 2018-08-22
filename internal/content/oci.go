package content

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/pkg/archive"
	digest "github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

func createOCIImage(r io.Reader) (io.Reader, error) {
	// Create an OCI tarball from our app tarball
	// index.json to describe it ?
	// blobs/… for the layer (i.e. where we are going to put our reader)
	dir, err := ioutil.TempDir("", "create-oci-image")
	if err != nil {
		return nil, err
	}
	fmt.Println(dir)
	now := time.Now()
	// 0. Get diffid (content hashed)
	var diffID digest.Digest
	// 1. write the layer and digest
	layerDesc, err := createLayerBlob(dir, r)
	if err != nil {
		return nil, err
	}
	diffID = layerDesc.Digest
	//	fmt.Println("layer", layerDesc)
	imageConfig := ocispec.Image{
		Created:      &now,
		Architecture: "config",
		OS:           "config",
		Config: ocispec.ImageConfig{
			Labels: map[string]string{},
		},
		RootFS: ocispec.RootFS{
			Type:    "layers",
			DiffIDs: []digest.Digest{diffID}, //nope { payloadDesc.Digest},
		},
		History: []ocispec.History{
			{CreatedBy: "COPY configfile /"},
		},
	}
	// 2. write the oci config and digest
	configDesc, err := createConfigFile(dir, imageConfig)
	if err != nil {
		return nil, err
	}
	// 3. write the manifest (?) and digest
	var manifest ocispec.Manifest
	manifest.SchemaVersion = 2
	manifest.Config = configDesc
	manifest.Layers = []ocispec.Descriptor{layerDesc}
	manifestDesc, err := createManifestFile(dir, manifest)
	if err != nil {
		return nil, err
	}
	// fmt.Println("manifest", manifestDesc)
	// 4. write index.json & co
	var index ocispec.Index
	index.SchemaVersion = 2
	index.Manifests = []ocispec.Descriptor{manifestDesc}
	if err := createIndexFile(dir, index); err != nil {
		return nil, err
	}
	// 5. write layout
	if err := createLayoutFile(dir); err != nil {
		return nil, err
	}
	return archive.TarWithOptions(dir, &archive.TarOptions{
		Compression: archive.Uncompressed,
	})
}

func createLayoutFile(root string) error {
	var layout ocispec.ImageLayout
	layout.Version = ocispec.ImageLayoutVersion
	contents, err := json.Marshal(layout)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(root, ocispec.ImageLayoutFile), contents, 0644)
}

func createLayerBlob(root string, inTar io.Reader) (ocispec.Descriptor, error) {
	return createBlob(root, inTar, "application/vnd.docker.image.rootfs.diff.tar")
}

func createIndexFile(root string, index ocispec.Index) error {
	content, err := json.Marshal(index)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(root, "index.json"), content, 0644)
}

func createManifestFile(root string, manifest ocispec.Manifest) (ocispec.Descriptor, error) {
	content, err := json.Marshal(manifest)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	return createBlob(root, bytes.NewBuffer(content), "application/vnd.docker.distribution.manifest.v2+json")
}

func createConfigFile(root string, config ocispec.Image) (ocispec.Descriptor, error) {
	content, err := json.Marshal(config)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	return createBlob(root, bytes.NewBuffer(content), "application/vnd.docker.container.image.v1+json")
}

func createBlob(root string, stream io.Reader, mediaType string) (ocispec.Descriptor, error) {
	name := filepath.Join(root, "blobs", "sha256", ".tmp-blob")
	err := os.MkdirAll(filepath.Dir(name), 0700)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	f, err := os.Create(name)
	if err != nil {
		return ocispec.Descriptor{}, err
	}
	defer f.Close()

	digester := digest.SHA256.Digester()
	tee := io.TeeReader(stream, digester.Hash())
	size, err := io.Copy(f, tee)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	if err := f.Sync(); err != nil {
		return ocispec.Descriptor{}, err
	}

	if err := f.Chmod(0644); err != nil {
		return ocispec.Descriptor{}, err
	}

	if err := digester.Digest().Validate(); err != nil {
		return ocispec.Descriptor{}, err
	}

	err = os.Rename(name, filepath.Join(filepath.Dir(name), digester.Digest().Hex()))
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	return ocispec.Descriptor{
		Digest:    digester.Digest(),
		Size:      size,
		MediaType: mediaType,
	}, nil
}
