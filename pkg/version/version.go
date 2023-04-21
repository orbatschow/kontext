package version

import (
	_ "embed"
	"strings"
)

const undefined = "undefined"

//go:generate sh -c "git name-rev --tags --name-only $(git rev-parse HEAD) > data/tag.txt"
//go:embed data/tag.txt
var Tag string

//go:generate sh -c "git rev-list --max-count=1 HEAD > data/commit.txt"
//go:embed data/commit.txt
var Commit string

func Compute() string {
	Tag = strings.ReplaceAll(Tag, "\n", "")
	Commit = strings.ReplaceAll(Commit, "\n", "")

	if Tag != undefined && len(Tag) > 0 {
		return Tag
	}
	if Commit != undefined && len(Commit) > 0 {
		return Commit
	}
	return "(devel)"
}
