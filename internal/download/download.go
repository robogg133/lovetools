package download

import (
	"context"
	"fmt"
	"net/http"
)

const (
	ReleasePrefix = "https://github.com/love2d/love/releases/download/%s/"
)

type DownloadContext struct {
	ctx     context.Context
	workdir string
	version string
	url     string
}

func NewDownloadContext(ctx context.Context, version, workdir string) *DownloadContext {
	return &DownloadContext{
		ctx:     ctx,
		workdir: workdir,
		version: version,
		url:     fmt.Sprintf(ReleasePrefix, version) + fmt.Sprintf(ReleaseCompactedName, version),
	}
}

func (d *DownloadContext) Fetch() (*http.Response, error) {
	fmt.Println(d.url)
	req, err := http.NewRequestWithContext(d.ctx, "GET", d.url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch: excpecting status code 200 but got %d", resp.StatusCode)
	}

	return resp, nil
}
