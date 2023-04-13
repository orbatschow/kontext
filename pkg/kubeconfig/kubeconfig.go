package kubeconfig

import (
	"fmt"
	"io"
	"os"

	"github.com/orbatschow/kontext/pkg/logger"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func Read(reader io.Reader) (*api.Config, error) {
	log := logger.New()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read kubeconfig, err: '%w'", err)
	}

	buffer, err := clientcmd.Load(data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal kubeconfig, err: '%w'", err)
	}
	log.Debug("read kubeconfig")

	return buffer, nil
}

func Write(writer io.Writer, apiConfig *api.Config) error {
	log := logger.New()

	if apiConfig == nil {
		return fmt.Errorf("invalid api config")
	}

	var buffer []byte

	buffer, err := clientcmd.Write(*apiConfig)
	if err != nil {
		return fmt.Errorf("persist new kubeconfig, err: '%w'", err)
	}

	_, err = writer.Write(buffer)
	if err != nil {
		return err
	}

	log.Debug("wrote kubeconfig")
	return nil
}

func Merge(files ...*os.File) (*api.Config, error) {
	var buffer []string
	// log := logger.New()

	for _, file := range files {
		buffer = append(buffer, file.Name())
	}

	loadingRules := clientcmd.ClientConfigLoadingRules{
		Precedence: buffer,
	}

	matches, err := loadingRules.Load()
	if err != nil {
		return nil, fmt.Errorf("could not merge kubeconfig files, err: '%w'", err)
	}
	// TODO: replace with print function, that prints a table
	// log.Info("merged kubeconfig", log.ArgsFromMap(buildFileMap(files)))
	return matches, nil
}
