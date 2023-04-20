package backup

import (
	"os"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/samber/lo"
)

// reconcileRevisions adds a new backup revision to the current state
func reconcileRevisions(config *config.Config, currentState *state.State, backupFile *os.File) error {
	revisions := currentState.Backup.Revisions

	// add the new revision
	revisions = append(revisions, state.Revision(backupFile.Name()))

	// if the length of the current revisions do not exceed the maximum
	// revision size, return the current revisions
	if len(revisions) <= config.Backup.Revisions || len(revisions) == 0 {
		return nil
	}

	// get all revisions, that exceed the maximum revision size
	overflow := lo.Slice(revisions, 0, len(revisions)-config.Backup.Revisions)

	for _, revision := range overflow {
		err := removeRevision(revision)
		if err != nil {
			return err
		}
	}

	// compute the new revisions
	currentState.Backup.Revisions = lo.Slice(revisions, len(revisions)-config.Backup.Revisions, len(revisions))

	return nil
}

func removeRevision(revision state.Revision) error {
	log := logger.New()

	file, err := os.Open(string(revision))
	if err != nil {
		return err
	}

	log.Info("removing backup revision", log.Args("file", file.Name()))
	err = os.Remove(file.Name())
	if err != nil {
		return err
	}

	return nil
}
