package kubeconfig

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/orbatschow/kontext/pkg/config"
	"k8s.io/client-go/tools/clientcmd/api"
)

func Test_Load(t *testing.T) {
	type args struct {
		file *os.File
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
				file: func() *os.File {
					_, caller, _, _ := runtime.Caller(0)
					kubeConfigFile := filepath.Join(caller, "..", "testdata", "01-valid-kubeconfig.yaml")
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
			name: "should throw an error, because the config points to nil",
			args: args{
				file: func() *os.File {
					file, _ := os.Open("")
					return file
				}(),
			},
			wantErr: true,
		},
		{
			name: "should built an empty kubeconfig, despite the base file being invalid",
			args: args{
				file: func() *os.File {
					_, caller, _, _ := runtime.Caller(0)
					kubeConfigFile := filepath.Join(caller, "..", "testdata", "02-invalid-kubeconfig.yaml")
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
			got, err := Read(tt.args.file)
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got: '%v'", err)
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
		apiConfig *api.Config
		config    func(file *os.File) *config.Config
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
				config: func(tmpFile *os.File) *config.Config {
					kontextConfig := &config.Config{
						Global: config.Global{
							Kubeconfig: tmpFile.Name(),
						},
					}
					return kontextConfig
				},
				apiConfig: func() *api.Config {
					_, caller, _, _ := runtime.Caller(0)
					kubeConfigFile := filepath.Join(caller, "..", "testdata", "01-valid-kubeconfig.yaml")
					file, err := os.Open(kubeConfigFile)
					if err != nil {
						t.Errorf("%v", err)
					}
					apiConfig, err := Read(file)
					if err != nil {
						t.Errorf("%v", err)
					}
					return apiConfig
				}(),
			},
			wantErr: false,
		},
		{
			name: "should throw an error due to missing config",
			args: args{
				config: func(tmpFile *os.File) *config.Config {
					kontextConfig := &config.Config{
						Global: config.Global{
							Kubeconfig: "",
						},
					}
					return kontextConfig
				},
				apiConfig: nil,
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

			err = Write(tmpFile, tt.args.apiConfig)
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got: '%v'", err)
			}

			if tt.after != nil {
				tt.after(tmpFile)
			}
		})
	}
}

func Test_Merge(t *testing.T) {
	type args struct {
		files []*os.File
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
				files: func() []*os.File {
					var buffer []*os.File

					_, caller, _, _ := runtime.Caller(0)
					filenames := []string{
						"03-kontext-merge-1.yaml",
						"04-kontext-merge-2.yaml",
						"05-kontext-merge-3.yaml",
					}

					for _, filename := range filenames {
						path := filepath.Join(caller, "..", "testdata", filename)
						file, err := os.Open(path)
						if err != nil {
							t.Errorf("%v", err)
						}
						buffer = append(buffer, file)
					}

					return buffer
				}(),
			},
			wantErr: false,
		},
		{
			name: "should merge the kubeconfig successfully, even if a kubeconfig file is invalid",
			args: args{
				files: func() []*os.File {
					var buffer []*os.File

					_, caller, _, _ := runtime.Caller(0)
					filenames := []string{
						"02-invalid-kubeconfig.yaml",
						"04-kontext-merge-2.yaml",
						"05-kontext-merge-3.yaml",
					}

					for _, filename := range filenames {
						path := filepath.Join(caller, "..", "testdata", filename)
						file, err := os.Open(path)
						if err != nil {
							t.Errorf("%v", err)
						}
						buffer = append(buffer, file)
					}

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
				t.Errorf("unexpected error, err: '%v'", err)
			}
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got: '%v'", err)
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
