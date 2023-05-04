package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/orbatschow/kontext/pkg/backup/revision"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/kubeconfig"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
)

type Filename string
type Directory string

type Reconciler struct {
	Config *config.Config
	State  *state.State
}

// Reconcile creates a new backup revision (if backups are enabled), updates the state and cleans up old revisions
func (r *Reconciler) Reconcile() error {
	log := logger.New()

	if !r.Config.Backup.Enabled {
		log.Warn("skipping backup, it is disabled")
		return nil
	}

	// create a new backup
	backupFile, err := r.create()
	if err != nil {
		return err
	}

	// add the new backup revision and remove revisions, that exceed the limit
	revisionReconciler := revision.Reconciler{
		Config: r.Config,
		State:  r.State,
		Backup: backupFile,
	}
	revisions, err := revisionReconciler.Reconcile()
	if err != nil {
		return err
	}
	r.State.Backup.Revisions = revisions

	err = state.Write(r.Config, r.State)
	if err != nil {
		return err
	}

	return nil
}

// create creates a new backup revision
func (r *Reconciler) create() (*os.File, error) {
	log := logger.New()

	file, err := os.Open(r.Config.Global.Kubeconfig)
	if err != nil {
		return nil, err
	}
	apiConfig, err := kubeconfig.Read(file)
	if err != nil {
		return nil, err
	}

	backupFilename := computeBackupFileName(r.Config)

	if _, err := os.Stat(r.Config.Backup.Directory); os.IsNotExist(err) {
		err = os.MkdirAll(r.Config.Backup.Directory, 0755)
		if err != nil {
			return nil, fmt.Errorf("could not create backup directory, err: '%w'", err)
		}
	}

	backupFile, err := os.OpenFile(string(backupFilename), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	err = kubeconfig.Write(backupFile, apiConfig)
	if err != nil {
		return nil, err
	}
	log.Trace("created new backup", log.Args("file", backupFile.Name()))

	return backupFile, nil
}

// computeBackupFileName builds the file name for the new backup
func computeBackupFileName(config *config.Config) Filename {
	// compute the current timestamp
	timestamp := int(time.Now().UnixNano() / int64(time.Millisecond))
	// compute the backup file name
	backupFileName := fmt.Sprintf("kubeconfig-%d.yaml", timestamp)
	// compute the backup file path
	backupFilePath := filepath.Join(config.Backup.Directory, backupFileName)
	return Filename(backupFilePath)
}
