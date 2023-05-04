package kubeconfig

import (
	"fmt"
	"io"
	"os"

	"github.com/orbatschow/kontext/pkg/logger"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func Read(file *os.File) (*api.Config, error) {
	log := logger.New()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read kubeconfig, err: '%w'", err)
	}

	buffer, err := clientcmd.Load(data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal kubeconfig, err: '%w'", err)
	}
	log.Debug("read kubeconfig", log.Args("file", file.Name()))

	return buffer, nil
}

func Write(file *os.File, apiConfig *api.Config) error {
	log := logger.New()

	if apiConfig == nil {
		return fmt.Errorf("invalid api config")
	}

	var buffer []byte

	buffer, err := clientcmd.Write(*apiConfig)
	if err != nil {
		return fmt.Errorf("persist new kubeconfig, err: '%w'", err)
	}

	_, err = file.Write(buffer)
	if err != nil {
		return err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	log.Debug("wrote kubeconfig", log.Args("file", file.Name()))
	return nil
}

func Merge(files ...*os.File) (*api.Config, error) {
	var buffer []string

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
	return matches, nil
}
