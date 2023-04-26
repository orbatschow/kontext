package source

import (
	"fmt"
	"os"
	"strconv"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/utils/glob"
	"github.com/pterm/pterm"
	"github.com/samber/lo"
)

// ComputeFiles computes the target files for the given SourceItem
// 1. compute all include globs and remove possible duplicates
// 2. compute all exclude globs and remove possible duplicates
// 3. build the difference for included and excluded files and return the result
func ComputeFiles(source *config.SourceItem) ([]*os.File, error) {
	log := logger.New()
	var buffer []*os.File

	log.Info("computing files for source item", log.Args("name", source.Name))

	includes, err := computeIncludes(source)
	if err != nil {
		return nil, err
	}
	excludes, err := computeExcludes(source)
	if err != nil {
		return nil, err
	}

	buffer, _ = difference(includes, excludes)

	return buffer, nil
}

func computeIncludes(source *config.SourceItem) ([]*os.File, error) {
	log := logger.New()
	var buffer []*os.File

	for _, include := range source.Include {
		log.Info("expanding include glob", log.Args("glob", include))

		matches, err := glob.Expand(glob.Pattern(include))
		if err != nil {
			return nil, fmt.Errorf("could not compute glob, err: '%w'", err)
		}
		buffer = append(buffer, matches...)

		log.Info("matched files:", lo.Map(matches, func(item *os.File, index int) pterm.LoggerArgument {
			return pterm.LoggerArgument{
				Key:   strconv.Itoa(index),
				Value: item.Name(),
			}
		}))
	}

	// remove duplicates
	buffer = lo.UniqBy(buffer, func(item *os.File) string {
		return item.Name()
	})

	return buffer, nil
}

func computeExcludes(sourceItem *config.SourceItem) ([]*os.File, error) {
	log := logger.New()
	var buffer []*os.File

	for _, exclude := range sourceItem.Exclude {
		log.Info("expanding exclude glob", log.Args("glob", exclude))

		matches, err := glob.Expand(glob.Pattern(exclude))
		if err != nil {
			return nil, fmt.Errorf("could not compute glob, err: '%w'", err)
		}

		buffer = append(buffer, matches...)

		log.Info("matched files:", lo.Map(matches, func(item *os.File, index int) pterm.LoggerArgument {
			return pterm.LoggerArgument{
				Key:   strconv.Itoa(index),
				Value: item.Name(),
			}
		}))
	}

	// remove duplicates
	buffer = lo.UniqBy(buffer, func(item *os.File) string {
		return item.Name()
	})

	return buffer, nil
}

func difference(includes []*os.File, excludes []*os.File) ([]*os.File, []*os.File) {
	var left []*os.File
	var right []*os.File

	seenLeft := map[string]struct{}{}
	seenRight := map[string]struct{}{}

	for _, elem := range includes {
		seenLeft[elem.Name()] = struct{}{}
	}

	for _, elem := range excludes {
		seenRight[elem.Name()] = struct{}{}
	}

	for _, elem := range includes {
		if _, ok := seenRight[elem.Name()]; !ok {
			left = append(left, elem)
		}
	}

	for _, elem := range excludes {
		if _, ok := seenLeft[elem.Name()]; !ok {
			right = append(right, elem)
		}
	}

	return left, right
}
