package download

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

const (
	ReleaseCompactedName = "love-%s-win32.zip"
	archivePrefix        = "love-%s-win32/"
)

// Extract the fetched archive file to the workdir
func (d *DownloadContext) Extract(resp *http.Response) error {
	defer resp.Body.Close()

	prefix := fmt.Sprintf(archivePrefix, d.version)

	pb := progressbar.Default(resp.ContentLength)
	defer pb.Close()

	pb.Describe(fmt.Sprintf("Downloading Love %s", d.version))

	rawReader := progressbar.NewReader(resp.Body, pb)
	f, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("love-%s-zip", d.version))
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, &rawReader); err != nil {
		return err
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return err
	}

	zr, err := zip.NewReader(f, resp.ContentLength)
	if err != nil {
		return err
	}
	resp.Body.Close()
	pb.Close()

	pb = progressbar.Default(int64(len(zr.File)), "Extracting")
	defer pb.Close()
	for _, file := range zr.File {
		name := filepath.Clean(file.Name)
		name = strings.TrimPrefix(name, prefix)

		if file.Mode().IsDir() {
			os.MkdirAll(filepath.Join(d.workdir, name), file.Mode())
			pb.Add(1)
			continue
		}

		f, err := os.OpenFile(filepath.Join(d.workdir, name), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, file.Mode())
		if err != nil {
			return err
		}

		zf, err := file.Open()
		if err != nil {
			f.Close()
			return err
		}
		if _, err := io.Copy(f, zf); err != nil {
			zf.Close()
			f.Close()
			return err
		}
		zf.Close()
		f.Close()
		pb.Add(1)
	}

	return nil
}
