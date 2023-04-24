package revision

import (
	"os"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/samber/lo"
)

type Reconciler struct {
	Config *config.Config
	State  *state.State
	Backup *os.File
}

func (r *Reconciler) Reconcile() ([]state.Revision, error) {
	revisions := r.State.Backup.Revisions

	// add the new revision
	revisions = append(revisions, state.Revision(r.Backup.Name()))

	// if the length of the current revisions does not exceed the maximum
	// revision size skip the cleanup
	if len(revisions) <= r.Config.Backup.Revisions || len(revisions) == 0 {
		return revisions, nil
	}

	// get all revisions, that exceed the maximum revision size
	overflow := lo.Slice(revisions, 0, len(revisions)-r.Config.Backup.Revisions)

	// remove the previously matched conditions
	for _, revision := range overflow {
		err := deleteRevisions(revision)
		if err != nil {
			return nil, err
		}
	}

	// compute the new revisions
	revisions = lo.Slice(revisions, len(revisions)-r.Config.Backup.Revisions, len(revisions))

	return revisions, nil
}

func deleteRevisions(revision state.Revision) error {
	log := logger.New()

	log.Info("removing backup revision", log.Args("file", revision))
	err := os.Remove(string(revision))
	if err != nil {
		return err
	}

	return nil
}
