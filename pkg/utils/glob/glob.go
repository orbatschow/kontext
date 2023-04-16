package glob

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
)

type Pattern string

// Expand takes a glob and returns a slice of absolute filepaths
func Expand(glob Pattern) ([]*os.File, error) {
	var buffer []*os.File

	basePath, pattern := doublestar.SplitPattern(string(glob))
	fileSystem := os.DirFS(basePath)

	matches, err := doublestar.Glob(fileSystem, pattern)
	if err != nil {
		return nil, err
	}

	for _, filename := range matches {
		file, err := os.Open(filepath.Join(basePath, filename))
		if err != nil {
			return nil, fmt.Errorf("could not open file: '%s'", filename)
		}
		buffer = append(buffer, file)
	}

	return buffer, err
}
