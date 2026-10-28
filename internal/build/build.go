package build

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/codeglyph/go-dotignore/v2"
)

const (
	IgnoreFile = ".loveignore"
)

type Builder struct {
	tmpDir string

	useLuajit bool
	useTeal   bool

	match *dotignore.PatternMatcher

	compiledLuaFile ReadWriterSeekerCloser
}

func NewBuilder(luajit bool, teal bool) (*Builder, error) {
	dir, err := os.MkdirTemp(os.TempDir(), "lovetools-build")
	if err != nil {
		return nil, err
	}
	builder := &Builder{
		tmpDir:    dir,
		useLuajit: luajit,
		useTeal:   teal,
	}

	if luajit {
		f, err := os.OpenFile(filepath.Join(dir, "luafile.luac"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0700)
		if err != nil {
			return nil, err
		}
		builder.compiledLuaFile = f
	}

	return builder, nil
}

func (b *Builder) Build(dir string) error {
	var err error
	b.match, err = ignoreFile(dir)
	if err != nil {
		return err
	}

	if err := b.zipDir(dir); err != nil {
		return err
	}

	return nil
}
func (b *Builder) Close() error {
	return nil
}

func (b *Builder) zipDir(dir string) error {
	f, err := os.OpenFile(filepath.Join(b.tmpDir, "game.love"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
	if err != nil {
		return err
	}
	defer f.Close()
	w := zip.NewWriter(f)
	defer w.Close()

	if err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		zipPath := strings.TrimPrefix(path, dir)
		zipPath = strings.TrimPrefix(zipPath, "/")
		if matched, err := b.match.Matches(zipPath); err != nil {
			return err
		} else if matched {
			return nil
		}

		if filepath.Ext(zipPath) == ".tl" && b.useTeal {
			zipPath = strings.TrimSuffix(zipPath, ".tl")
			zipPath = zipPath + ".lua"
		}
		fWriter, err := w.Create(zipPath)
		if err != nil {
			return err
		}

		ext := filepath.Ext(path)
		if ext != ".lua" && !(ext == ".tl" && b.useTeal) {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			if _, err := f.WriteTo(fWriter); err != nil {
				f.Close()
				return err
			}
			return f.Close()
		}
		if ext == ".tl" {
			path, err = b.compileTeal(path)
			if err != nil {
				return err
			}
		}

		f, err := b.compileLuaJIT(path)
		if err != nil {
			return err
		}
		if _, err := io.Copy(fWriter, f); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	if b.useLuajit {
		os.Remove(filepath.Join(b.tmpDir, "luafile.luac"))
	}
	if b.useTeal {
		os.Remove(filepath.Join(b.tmpDir, "luafile.lua"))
	}

	return nil
}

func (b *Builder) compileTeal(filename string) (string, error) {
	if !b.useTeal {
		return filename, nil
	}

	dst := filepath.Join(b.tmpDir, "luafile.lua")
	cmd := exec.Command("tl", "gen", filename, "-o", dst)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%s: %s", err.Error(), string(output))
	}

	return dst, nil
}

func (b *Builder) compileLuaJIT(filename string) (io.ReadCloser, error) {
	if !b.useLuajit {
		return os.Open(filename)
	}
	if _, err := b.compiledLuaFile.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	cmd := exec.Command("luajit", "-b", filename, filepath.Join(b.tmpDir, "luafile.luac"))
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("%s: %s", err.Error(), string(output))
	}

	return b.compiledLuaFile, nil
}

func BuildWindows(dir string) {

}
