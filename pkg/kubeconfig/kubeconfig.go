package kubeconfig

import (
	"fmt"
	"io"
	"strconv"

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
	log.Info("read kubeconfig")

	return buffer, nil
}

func Write(writer io.Writer, apiConfig *api.Config) error {
	log := logger.New()
	var buffer []byte

	buffer, err := clientcmd.Write(*apiConfig)
	if err != nil {
		return fmt.Errorf("persist new kubeconfig, err: '%w'", err)
	}

	_, err = writer.Write(buffer)
	if err != nil {
		return err
	}

	log.Info("wrote kubeconfig")
	return nil
}

func Merge(files ...string) (*api.Config, error) {
	log := logger.New()
	loadingRules := clientcmd.ClientConfigLoadingRules{
		Precedence: files,
	}

	buffer, err := loadingRules.Load()
	if err != nil {
		return nil, fmt.Errorf("could not merge kubeconfig files, err: '%w'", err)
	}
	log.Info("merged kubeconfig", log.ArgsFromMap(buildFileMap(files)))
	return buffer, nil
}

func buildFileMap(files []string) map[string]any {
	buffer := map[string]any{}

	for index, file := range files {
		buffer[strconv.Itoa(index)] = file
	}

	return buffer
}
