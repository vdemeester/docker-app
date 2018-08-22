package types

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/docker/app/internal"
	"github.com/docker/docker/pkg/archive"
)

func (a *App) Extract(path string) error {
	if err := ioutil.WriteFile(filepath.Join(path, internal.MetadataFileName), a.Metadata(), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(path, internal.ComposeFileName), a.Composes()[0], 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(path, internal.SettingsFileName), a.Settings()[0], 0755); err != nil {
		return err
	}
	return nil
}

func (a *App) TarReader() (io.ReadCloser, error) {
	dir, err := ioutil.TempDir("", "tar-it")
	if err != nil {
		return nil, err
	}
	fmt.Println(dir)
	//defer os.RemoveAll(dir)
	if err := a.Extract(dir); err != nil {
		return nil, err
	}
	return archive.TarWithOptions(dir, &archive.TarOptions{
		Compression: archive.Uncompressed,
	})
}
