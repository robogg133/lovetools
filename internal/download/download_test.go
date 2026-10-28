package download_test

import (
	"testing"

	"github.com/robogg133/lovetools/internal/download"
)

func TestDownload(t *testing.T) {
	d := download.NewDownloadContext(t.Context(), "11.5", t.ArtifactDir())
	resp, err := d.Fetch()
	if err != nil {
		t.Fatal(err)
	}

	if err := d.Extract(resp); err != nil {
		t.Fatal(err)
	}

}
