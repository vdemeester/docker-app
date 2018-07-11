// Code generated by "esc -o bindata.go -pkg schema -ignore .*\.go -private -modtime=1518458244 data"; DO NOT EDIT.

package schema

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// _escFS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func _escFS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// _escDir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func _escDir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// _escFSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func _escFSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// _escFSMustByte is the same as _escFSByte, but panics if name is not present.
func _escFSMustByte(useLocal bool, name string) []byte {
	b, err := _escFSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// _escFSString is the string version of _escFSByte.
func _escFSString(useLocal bool, name string) (string, error) {
	b, err := _escFSByte(useLocal, name)
	return string(b), err
}

// _escFSMustString is the string version of _escFSMustByte.
func _escFSMustString(useLocal bool, name string) string {
	return string(_escFSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/data/config_schema_v3.0.json": {
		local:   "data/config_schema_v3.0.json",
		size:    11063,
		modtime: 1518458244,
		compressed: `
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
/BUAAP//72YpJjcrAAA=
`,
	},

	"/data/config_schema_v3.1.json": {
		local:   "data/config_schema_v3.1.json",
		size:    12209,
		modtime: 1518458244,
		compressed: `
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
`,
	},

	"/data/config_schema_v3.2.json": {
		local:   "data/config_schema_v3.2.json",
		size:    13755,
		modtime: 1518458244,
		compressed: `
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
`,
	},

	"/data/config_schema_v3.3.json": {
		local:   "data/config_schema_v3.3.json",
		size:    15414,
		modtime: 1518458244,
		compressed: `
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
ZccPpEaO8N0m78L+f178PwAA//9GlSx0NjwAAA==
`,
	},

	"/data/config_schema_v3.4.json": {
		local:   "data/config_schema_v3.4.json",
		size:    15797,
		modtime: 1518458244,
		compressed: `
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
`,
	},

	"/data/config_schema_v3.5.json": {
		local:   "data/config_schema_v3.5.json",
		size:    16725,
		modtime: 1518458244,
		compressed: `
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
`,
	},

	"/data/config_schema_v3.6.json": {
		local:   "data/config_schema_v3.6.json",
		size:    17007,
		modtime: 1518458244,
		compressed: `
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
DoTuq9lAU9Q4O1xV/59X/w8AAP//zRo7vm9CAAA=
`,
	},

	"/data/config_schema_v3.7.json": {
		local:   "data/config_schema_v3.7.json",
		size:    17777,
		modtime: 1518458244,
		compressed: `
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
`,
	},

	"/": {
		isDir: true,
		local: "",
	},

	"/data": {
		isDir: true,
		local: "data",
	},
}
