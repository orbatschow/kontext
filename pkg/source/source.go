package source

import (
	"fmt"
	"os"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/utils/glob"
	"github.com/samber/lo"
)

func ComputeFiles(source *config.SourceItem) ([]*os.File, error) {
	log := logger.New()
	var buffer []*os.File

	// add all file globs from source.Include
	for _, include := range source.Include {
		log.Info("expanding include glob", log.Args("glob", include))

		matches, err := glob.Expand(glob.Pattern(include))
		if err != nil {
			return nil, fmt.Errorf("could not compute glob, err: '%w'", err)
		}
		buffer = append(buffer, matches...)

		// TODO: replace with print function, that prints a table
		// log.Info("found matching buffer", log.ArgsFromMap(computeFileMap(matches)))
	}

	// remove duplicates
	buffer = lo.UniqBy(buffer, func(item *os.File) string {
		return item.Name()
	})
	// TODO: replace with print function, that prints a table
	// log.Info("computed kubeconfig buffer without duplicates", log.ArgsFromMap(computeFileMap(buffer)))

	// remove all file globs from source.Exclude
	for _, exclude := range source.Exclude {
		log.Warn("expanding exclude glob", log.Args("glob", exclude))

		matches, err := glob.Expand(glob.Pattern(exclude))
		if err != nil {
			return nil, fmt.Errorf("could not compute glob, err: '%w'", err)
		}

		for i, file := range buffer {
			_, ok := lo.Find(matches, func(item *os.File) bool {
				return item.Name() == file.Name()
			})

			if ok {
				buffer[i] = nil
			}
		}

		buffer = lo.Without(buffer, nil)

		// TODO: replace with print function, that prints a table
		// log.Warn("found matching buffer", log.ArgsFromMap(computeFileMap(matches)))
	}

	// TODO: replace with print function, that prints a table
	// log.Info("computed final kubeconfig buffer", log.ArgsFromMap(computeFileMap(buffer)))
	return buffer, nil
}
