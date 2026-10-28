package download

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

const (
	ReleaseCompactedName = "love-%s-linux-src.tar.gz"
	archivePrefix        = "love-%s/"
)

// Extract the fetched archive file to the workdir
func (d *DownloadContext) Extract(resp *http.Response) error {
	defer resp.Body.Close()

	prefix := fmt.Sprintf(archivePrefix, d.version)

	pb := progressbar.Default(resp.ContentLength)
	defer pb.Close()

	pb.Describe(fmt.Sprintf("Downloading Love %s", d.version))

	rawReader := progressbar.NewReader(resp.Body, pb)
	r, err := gzip.NewReader(&rawReader)
	if err != nil {
		return err
	}

	tr := tar.NewReader(r)

	for {
		h, err := tr.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		name := filepath.Clean(h.Name)
		name = strings.TrimPrefix(name, prefix)

		switch h.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filepath.Join(d.workdir, name), os.FileMode(h.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Join(d.workdir, filepath.Dir(name)), 0755); err != nil {
				return err
			}

			f, err := os.OpenFile(filepath.Join(d.workdir, name), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(h.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}
	return nil
}

func (d *DownloadContext) Install(dst string) error {
	return nil
}
