package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pterm/pterm"
)

func Test_Read(t *testing.T) {
	type args struct {
		Environment map[string]string
		Reader      *Client
	}
	tests := []struct {
		name    string
		before  func(t *testing.T, environment map[string]string)
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "should read the config file successfully",
			args: args{
				Reader: &Client{
					Path: func() string {
						_, caller, _, _ := runtime.Caller(0)
						path := filepath.Join(caller, "..", "testdata", "01-valid-config.yaml")
						return path
					}(),
				},
			},
			want: &Config{
				Global: Global{
					Kubeconfig: "/home/nilsorbat/.config/kontext/kubeconfig.yaml",
					Verbosity:  pterm.LogLevelInfo,
				},
				Backup: Backup{
					Enabled: true,
				},
				State: State{
					Path:    os.ExpandEnv("$HOME/.local/state/kontext/state.json"),
					History: History{},
				},
				Groups: []Group{
					{
						Name:    "default",
						Context: "kind-local",
						Sources: []string{
							"default",
						},
					},
					{
						Name:    "dev",
						Context: "",
						Sources: []string{
							"dev",
						},
					},
				},
				Sources: []Source{
					{
						Name: "default",
						Include: []string{
							os.ExpandEnv("$HOME/.config/kontext/**/*.yaml"),
						},
						Exclude: []string{
							os.ExpandEnv("$HOME/.config/kontext/**/*prod*.yaml"),
						},
					},
					{
						Name: "dev",
						Include: []string{
							os.ExpandEnv("$HOME/.config/kontext/dev/**/*.yaml"),
						},
						Exclude: []string{
							os.ExpandEnv("$HOME/.config/kontext/dev/**/*prod*.yaml"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should read the config file successfully and set proper default values",
			before: func(t *testing.T, environment map[string]string) {
				for key, value := range environment {
					t.Setenv(key, value)
				}
			},
			args: args{
				Environment: func() map[string]string {
					_, caller, _, _ := runtime.Caller(0)
					path := filepath.Join(caller, "..", "testdata", "03-valid-kubeconfig.yaml")

					return map[string]string{
						"KUBECONFIG": path,
					}
				}(),
				Reader: &Client{
					Path: func() string {
						_, caller, _, _ := runtime.Caller(0)
						path := filepath.Join(caller, "..", "testdata", "02-valid-config-default-values.yaml")
						return path
					}(),
				},
			},
			want: &Config{
				Global: Global{
					Kubeconfig: func() string {
						_, caller, _, _ := runtime.Caller(0)
						path := filepath.Join(caller, "..", "testdata", "03-valid-kubeconfig.yaml")

						return path
					}(),
					Verbosity: pterm.LogLevelInfo,
				},
				Backup: Backup{
					Enabled: true,
				},
				State: State{
					Path:    os.ExpandEnv("$HOME/.local/state/kontext/state.json"),
					History: History{},
				},
				Groups:  []Group{},
				Sources: []Source{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t, tt.args.Environment)
			}

			got, err := tt.args.Reader.Read()

			if !tt.wantErr && err != nil {
				t.Errorf("config.Read() = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got: '%v'", err)
			}

			if !tt.wantErr && tt.want != nil && !cmp.Equal(tt.want, got) {
				diff := cmp.Diff(tt.want, got)
				t.Errorf("reader.Read() mismatch (-want +got):\n%s", diff)
			}

		})
	}
}
