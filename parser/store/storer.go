package store

import "io"

type Storer interface {
	ListObjects(dir string) ([]string, error)
	PutObject(path string, r io.Reader) error
	GetObject(path string) (io.Reader, error)
}
