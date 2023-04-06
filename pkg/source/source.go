package source

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/samber/lo"
)

func Expand(source *config.Source) ([]string, error) {
	log := logger.New()
	var files []string

	// add all file globs from source.Include
	for _, include := range source.Include {
		log.Info("expanding include glob", log.Args("glob", include))

		buffer, err := computeGlob(include)
		if err != nil {
			return nil, fmt.Errorf("could not compute glob, err: '%w'", err)
		}
		files = append(files, buffer...)

		log.Info("found matching files", log.ArgsFromMap(computeFileMap(buffer)))
	}

	// remove duplicates
	files = lo.Uniq(files)
	log.Info("computed kubeconfig files without duplicates", log.ArgsFromMap(computeFileMap(files)))

	// remove all file globs from source.Exclude
	for _, exclude := range source.Exclude {
		log.Warn("expanding exclude glob", log.Args("glob", exclude))

		buffer, err := computeGlob(exclude)
		if err != nil {
			return nil, fmt.Errorf("could not compute glob, err: '%w'", err)
		}
		files, _ = lo.Difference(files, buffer)

		log.Warn("found matching files", log.ArgsFromMap(computeFileMap(buffer)))
	}

	log.Info("computed final kubeconfig files", log.ArgsFromMap(computeFileMap(files)))
	return files, nil
}

// compute glob takes a glob and returns a slice of absolute filepaths
func computeGlob(glob string) ([]string, error) {
	var basePath string
	basePath, glob = doublestar.SplitPattern(glob)
	fileSystem := os.DirFS(basePath)

	buffer, err := doublestar.Glob(fileSystem, glob)
	if err != nil {
		return nil, err
	}

	for index, file := range buffer {
		buffer[index] = path.Join(basePath, file)
	}

	return buffer, err
}

func computeFileMap(files []string) map[string]any {
	buffer := map[string]any{}

	for index, file := range files {
		buffer[strconv.Itoa(index)] = file
	}

	return buffer
}
