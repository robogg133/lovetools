package config

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

var (
	FolderRoot          = filepath.Join(xdg.DataHome, "lovetools")
	FolderInstallations = filepath.Join(FolderRoot, "installations")
)
