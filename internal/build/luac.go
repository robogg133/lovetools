package build

import "io"

type ReadWriterSeekerCloser interface {
	io.ReadWriteSeeker
	Close() error
}

type nop struct{}

func (n *nop) Write([]byte) (int, error)      { return 0, nil }
func (n *nop) Read([]byte) (int, error)       { return 0, nil }
func (n *nop) Seek(int64, int) (int64, error) { return 0, nil }
func (n *nop) Close() error                   { return nil }

var _ ReadWriterSeekerCloser = &nop{}
