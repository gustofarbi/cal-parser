package temp

import (
	"github.com/gustofarbi/parser/util/rand"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type dir struct {
	prefix string
}

func New() *dir {
	prefix := filepath.Join(os.TempDir(), rand.String(10))
	err := os.Mkdir(prefix, os.ModePerm) // todo maybe wrong perm
	if err != nil {
		panic(err)
	}

	return &dir{prefix: prefix}
}

func (d *dir) ListObjects(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(filepath.Join(d.prefix, dir), func(path string, d fs.DirEntry, err error) error {
		if path != "." && path != ".." {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (d *dir) PutObject(path string, r io.Reader) error {
	f, err := os.Open(filepath.Join(d.prefix, path))
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(f, r)
	return err
}

func (d *dir) GetObject(path string) (io.Reader, error) {
	return os.Open(filepath.Join(d.prefix, path))
}
