package backup

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

	backupFileName := computeBackupFileName(file.Name())
	backup, err := os.Create(backupFileName)
	if err != nil {
		return err
	}

	err = kubeconfig.Write(backup, apiConfig)
	if err != nil {
		return err
	}

	return nil
}

func computeBackupFileName(baseFileName string) string {

	basePath := filepath.Dir(baseFileName)
	fileName := filepath.Base(fileNameWithoutExtension(baseFileName))
	timestamp := strconv.Itoa(makeTimestamp())

	backupPath := path.Join(basePath, fmt.Sprintf("%s-%s.bak", fileName, timestamp))

	return backupPath
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func makeTimestamp() int {
	return int(time.Now().UnixNano() / int64(time.Millisecond))
}
