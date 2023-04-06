package kubeconfig

import (
	"io"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/orbatschow/kontext/pkg/config"
	"k8s.io/client-go/tools/clientcmd/api"
)

func Test_Load(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *api.Config
		wantErr bool
	}{
		{
			name: "should parse the kubeconfig successfully",
			args: args{
				reader: func() io.Reader {
					_, caller, _, _ := runtime.Caller(0)
					kubeConfigFile := path.Join(caller, "..", "testdata", "01-valid-kubeconfig.yaml")
					file, err := os.Open(kubeConfigFile)
					if err != nil {
						t.Errorf("%v", err)
					}
					return file
				}(),
			},
			want: &api.Config{
				CurrentContext: "kind-kontext",
			},
			wantErr: false,
		},
		{
			name: "should throw an error, because the KontextConfig points to nil",
			args: args{
				reader: func() io.Reader {
					file, _ := os.Open("")
					return file
				}(),
			},
			wantErr: true,
		},
		{
			name: "should built an empty kubeconfig, despite the base file being invalid",
			args: args{
				reader: func() io.Reader {
					_, caller, _, _ := runtime.Caller(0)
					kubeConfigFile := path.Join(caller, "..", "testdata", "02-invalid-kubeconfig.yaml")
					file, err := os.Open(kubeConfigFile)
					if err != nil {
						t.Errorf("%v", err)
					}
					return file
				}(),
			},
			want: &api.Config{
				CurrentContext: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.reader)
			if !tt.wantErr && err != nil {
				t.Errorf("expected error")
			}

			// check if the current context is equal
			if !tt.wantErr && tt.want.CurrentContext != got.CurrentContext {
				t.Errorf("want: '%s', got: '%s'", tt.want.CurrentContext, got.CurrentContext)
			}
		})
	}
}

func Test_Write(t *testing.T) {
	type args struct {
		APIConfig     *api.Config
		KontextConfig func(file *os.File) *config.Config
	}
	tests := []struct {
		name    string
		after   func(file *os.File)
		args    args
		wantErr bool
	}{
		{
			name: "should save the kubeconfig successfully",
			after: func(file *os.File) {
				err := os.Remove(file.Name())
				if err != nil {
					t.Errorf("%v", err)
				}
			},
			args: args{
				KontextConfig: func(tmpFile *os.File) *config.Config {
					kontextConfig := &config.Config{
						Global: config.Global{
							Kubeconfig:                tmpFile.Name(),
							ConfirmKubeconfigOverride: false,
						},
					}
					return kontextConfig
				},
				APIConfig: func() *api.Config {
					_, caller, _, _ := runtime.Caller(0)
					kubeConfigFile := path.Join(caller, "..", "testdata", "01-valid-kubeconfig.yaml")
					file, err := os.Open(kubeConfigFile)
					if err != nil {
						t.Errorf("%v", err)
					}
					apiConfig, err := Load(file)
					if err != nil {
						t.Errorf("%v", err)
					}
					return apiConfig
				}(),
			},
			wantErr: false,
		},
		{
			name: "should throw an error due to missing file",
			args: args{
				KontextConfig: func(tmpFile *os.File) *config.Config {
					kontextConfig := &config.Config{
						Global: config.Global{
							Kubeconfig:                "",
							ConfirmKubeconfigOverride: false,
						},
					}
					return kontextConfig
				},
				APIConfig: func() *api.Config {
					_, caller, _, _ := runtime.Caller(0)
					kubeConfigFile := path.Join(caller, "..", "testdata", "01-valid-kubeconfig.yaml")
					file, err := os.Open(kubeConfigFile)
					if err != nil {
						t.Errorf("%v", err)
					}
					apiConfig, err := Load(file)
					if err != nil {
						t.Errorf("%v", err)
					}
					return apiConfig
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "kontext-kubeconfig-*.yaml")
			if err != nil {
				t.Errorf("%v", err)
			}
			kontextConfig := tt.args.KontextConfig(tmpFile)

			err = Write(kontextConfig, tt.args.APIConfig)
			if !tt.wantErr && err != nil {
				t.Errorf("expected error")
			}

			if tt.after != nil {
				tt.after(tmpFile)
			}
		})
	}
}

func Test_Merge(t *testing.T) {
	type args struct {
		files []string
	}
	tests := []struct {
		name    string
		args    args
		want    *api.Config
		wantErr bool
	}{
		{
			name: "should merge the kubeconfig successfully",
			args: args{
				files: func() []string {
					var buffer []string

					_, caller, _, _ := runtime.Caller(0)
					mergeFileOne := path.Join(caller, "..", "testdata", "03-kontext-merge-1.yaml")
					mergeFileTwo := path.Join(caller, "..", "testdata", "04-kontext-merge-2.yaml")
					mergeFileThree := path.Join(caller, "..", "testdata", "05-kontext-merge-3.yaml")

					buffer = append(buffer, mergeFileOne)
					buffer = append(buffer, mergeFileTwo)
					buffer = append(buffer, mergeFileThree)

					return buffer
				}(),
			},
			wantErr: false,
		},
		{
			name: "should merge the kubeconfig successfully, even if a kubeconfig file is invalid",
			args: args{
				files: func() []string {
					var buffer []string

					_, caller, _, _ := runtime.Caller(0)
					mergeFileOne := path.Join(caller, "..", "testdata", "02-invalid-kubeconfig.yaml")
					mergeFileTwo := path.Join(caller, "..", "testdata", "04-kontext-merge-2.yaml")
					mergeFileThree := path.Join(caller, "..", "testdata", "05-kontext-merge-3.yaml")

					buffer = append(buffer, mergeFileOne)
					buffer = append(buffer, mergeFileTwo)
					buffer = append(buffer, mergeFileThree)

					return buffer
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Merge(tt.args.files...)
			if !tt.wantErr && err != nil {
				t.Errorf("expected error")
			}

			if !tt.wantErr && tt.want != nil && got == nil {
				t.Errorf("want: '%v', got: '%v'", tt.want, got)
			}

			if !tt.wantErr && got.CurrentContext == "" {
				t.Errorf("want: '%v', got: '%v'", tt.want, got)
			}
		})
	}
}
