package backup

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/adrg/xdg"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/kubeconfig"
	"github.com/orbatschow/kontext/pkg/logger"
)

func Create(config *config.Config) error {
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

	var backupDirectory string
	if len(config.Backup.Location) == 0 {
		backupDirectory = path.Join(xdg.DataHome, "kontext", "backup")
	}

	if _, err := os.Stat(backupDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(backupDirectory, 0755)
		if err != nil {
			return fmt.Errorf("could not create backup directory, err: '%w'", err)
		}
	}

	timestamp := strconv.Itoa(makeTimestamp())
	backupFileName := fmt.Sprintf("kubeconfig-%s.yaml", timestamp)
	backupFile, err := os.Create(path.Join(backupDirectory, backupFileName))
	if err != nil {
		return err
	}

	err = kubeconfig.Write(backupFile, apiConfig)
	if err != nil {
		return err
	}

	return nil
}

func makeTimestamp() int {
	return int(time.Now().UnixNano() / int64(time.Millisecond))
}
