package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/kubeconfig"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
)

const (
	defaultRevisionLimit = 10
)

type Filename string
type Directory string

// Reconcile creates a new backup revision (if backups are enabled), updates the state and cleans up old revisions
func Reconcile(config *config.Config, currentState *state.State) error {
	log := logger.New()

	if !config.Backup.Enabled {
		log.Warn("skipping backup, it is disabled")
		return nil
	}

	// create a new backup
	backupFile, err := create(config)
	if err != nil {
		return err
	}

	// add the new backup revision to the current currentState
	addRevision(currentState, backupFile)

	err = enforceRevisionLimit(config, currentState)
	if err != nil {
		return err
	}

	err = state.Write(config, currentState)
	if err != nil {
		return err
	}

	return nil
}

// create creates a new backup revision
func create(config *config.Config) (*os.File, error) {
	file, err := os.Open(config.Global.Kubeconfig)
	if err != nil {
		return nil, err
	}
	apiConfig, err := kubeconfig.Read(file)
	if err != nil {
		return nil, err
	}

	backupFilename := computeBackupFileName(config)

	if _, err := os.Stat(config.Backup.Directory); os.IsNotExist(err) {
		err = os.MkdirAll(config.Backup.Directory, 0755)
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

	return backupFile, nil
}

// addRevision adds a new backup revision to the current state
func addRevision(currentState *state.State, backupFile *os.File) {
	currentState.Backup.Revisions = append(currentState.Backup.Revisions, state.Revision(backupFile.Name()))
}

// enforceRevisionLimit will check if the given revision limit is reached and
// remove old backup revisions if necessary
func enforceRevisionLimit(config *config.Config, state *state.State) error {
	log := logger.New()

	var revisionLimit int
	if config.Backup.Revisions == nil {
		revisionLimit = defaultRevisionLimit
	} else {
		revisionLimit = *config.Backup.Revisions
	}
	if revisionLimit >= len(state.Backup.Revisions) {
		return nil
	}

	outdatedRevisions := len(state.Backup.Revisions) - revisionLimit
	if outdatedRevisions <= 0 {
		return nil
	}

	for i := 0; i < outdatedRevisions; i++ {
		file, err := os.Open(string(state.Backup.Revisions[i]))
		if err != nil {
			return err
		}

		log.Info("removing backup revision", log.Args("file", file.Name()))
		err = os.Remove(file.Name())
		if err != nil {
			return err
		}
	}

	return nil
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
