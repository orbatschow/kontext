package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/adrg/xdg"
	"github.com/google/go-cmp/cmp"
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
					File: func() string {
						_, caller, _, _ := runtime.Caller(0)
						path := filepath.Join(caller, "..", "testdata", "01-valid-config.yaml")
						return path
					}(),
				},
			},
			want: &Config{
				Global: Global{
					Kubeconfig: filepath.Join(xdg.ConfigHome, "kontext", "kubeconfig.yaml"),
				},
				Backup: Backup{
					Enabled: true,
				},
				State: State{
					File:    filepath.Join(xdg.StateHome, "kontext", "state.json"),
					History: History{},
				},
				Group: Group{
					Items: []GroupItem{
						{
							Name: "default",
							Context: Context{
								Default: "kind-local",
								Selection: Selection{
									Default: "kind-local",
									Sort:    "desc",
								},
							},
							Sources: []string{
								"default",
							},
						},
						{
							Name:    "dev",
							Context: Context{},
							Sources: []string{
								"dev",
							},
						},
					},
					Selection: Selection{
						Default: "dev",
						Sort:    "asc",
					},
				},
				Source: Source{
					Items: []SourceItem{{
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
					File: func() string {
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
				},
				Backup: Backup{
					Enabled:   true,
					Revisions: DefaultBackupRevisionLimit,
					Directory: filepath.Join(xdg.DataHome, "kontext", "backup"),
				},
				State: State{
					File: filepath.Join(xdg.StateHome, "kontext", "state.json"),
					History: History{
						Size: DefaultStateHistoryLimit,
					},
				},
				Group: Group{
					Items:     []GroupItem{},
					Selection: Selection{},
				},
				Source: Source{
					Items: []SourceItem{},
				},
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
