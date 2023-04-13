package backup

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/adrg/xdg"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/kubeconfig"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/utils/glob"
)

const (
	backupFileGlob = "kubeconfig-*.yaml"
)

type Filename string

func Reconcile(config *config.Config) error {
	err := create(config)
	if err != nil {
		return err
	}

	backups, err := list(config)
	if err != nil {
		return err
	}
}

func create(config *config.Config) error {
	log := logger.New()

	if !config.Backup.Enabled {
		log.Warn("skipping backup, it is disabled")
		return nil
	}

	file, err := os.Open(config.Global.Kubeconfig)
	if err != nil {
		return err
	}
	apiConfig, err := kubeconfig.Read(file)
	if err != nil {
		return err
	}

	if len(config.Backup.Location) == 0 {
		config.Backup.Location = path.Join(xdg.DataHome, "kontext", "backup")
	}

	if _, err := os.Stat(config.Backup.Location); os.IsNotExist(err) {
		err = os.MkdirAll(config.Backup.Location, 0755)
		if err != nil {
			return fmt.Errorf("could not create backup directory, err: '%w'", err)
		}
	}

	backupFileName, err := computeBackupFileName(config)
	backupFile, err := os.Create(string(backupFileName))
	if err != nil {
		return err
	}

	err = kubeconfig.Write(backupFile, apiConfig)
	if err != nil {
		return err
	}

	return nil
}

func remove(config *config.Config) error {
	log := logger.New()

	if !config.Backup.Enabled {
		log.Warn("skipping backup, it is disabled")
		return nil
	}

	file, err := os.Open(config.Global.Kubeconfig)
	if err != nil {
		return err
	}
	apiConfig, err := kubeconfig.Read(file)
	if err != nil {
		return err
	}

	if len(config.Backup.Location) == 0 {
		config.Backup.Location = path.Join(xdg.DataHome, "kontext", "backup")
	}

	if _, err := os.Stat(config.Backup.Location); os.IsNotExist(err) {
		err = os.MkdirAll(config.Backup.Location, 0755)
		if err != nil {
			return fmt.Errorf("could not create backup directory, err: '%w'", err)
		}
	}

	backupFileName, err := computeBackupFileName(config)
	backupFile, err := os.Create(string(backupFileName))
	if err != nil {
		return err
	}

	err = kubeconfig.Write(backupFile, apiConfig)
	if err != nil {
		return err
	}

	return nil
}

func list(config *config.Config) ([]*os.File, error) {
	pattern := glob.Pattern(path.Join(config.Backup.Location, backupFileGlob))
	buffer, err := glob.Expand(pattern)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func sortByDate([]os.File) []Filename {

}

func computeBackupFileName(config *config.Config) (Filename, error) {
	// compute the current timestamp
	timestamp := int(time.Now().UnixNano() / int64(time.Millisecond))
	// compute the backup file name
	backupFileName := fmt.Sprintf("kubeconfig-%s.yaml", timestamp)
	// compute the backup file path
	backupFilePath := path.Join(config.Backup.Location, backupFileName)
	return Filename(backupFilePath), nil
}
