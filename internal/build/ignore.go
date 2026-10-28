package build

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/codeglyph/go-dotignore/v2"
)

func ignoreFile(dir string) (*dotignore.PatternMatcher, error) {
	ignore, err := os.Open(filepath.Join(dir, IgnoreFile))
	if err != nil {
		return dotignore.NewPatternMatcher([]string{"/.git/"})
	} else {
		defer ignore.Close()
		var buff bytes.Buffer
		if _, err := buff.Write([]byte("/.git/\n/.loveignore\n")); err != nil {
			return nil, err
		}
		if _, err := ignore.WriteTo(&buff); err != nil {
			return nil, err
		}

		return dotignore.NewPatternMatcherFromReader(&buff)

	}
}
