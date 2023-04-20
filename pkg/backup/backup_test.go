package backup

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/kubeconfig"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/samber/lo"
)

func Test_Reconcile(t *testing.T) {
	type args struct {
		config *config.Config
		state  *state.State
	}
	tests := []struct {
		name    string
		after   func(*config.Config)
		args    args
		want    func(*config.Config) *state.State
		wantErr bool
	}{
		{
			name: "should create a backup and add a new revision to the state",
			after: func(config *config.Config) {
				directory := filepath.Dir(config.State.File)
				err := os.RemoveAll(directory)
				if err != nil {
					t.Errorf("%v", err)
				}
			},
			args: args{
				config: func() *config.Config {
					// temporary directory
					tempDirectory, err := os.MkdirTemp("", "kontext-")
					if err != nil {
						t.Errorf("%v", err)
					}

					// current working directory
					_, caller, _, _ := runtime.Caller(0)

					// backup directory, based on temporary directory
					tempBackupDir := filepath.Join(tempDirectory, "backup")
					err = os.MkdirAll(filepath.Join(tempDirectory, "backup"), 0755)
					if err != nil {
						t.Errorf("%v", err)
					}

					// config, that shall be backed up
					kubeconfigFilepath := filepath.Join(caller, "..", "testdata", "01-valid-kubeconfig.yaml")

					// temporary state file
					stateFilePath := filepath.Join(tempDirectory, "state.json")
					_, err = os.Create(stateFilePath)
					if err != nil {
						t.Errorf("%v", err)
					}

					kontextConfig := &config.Config{
						Global: config.Global{
							Kubeconfig: kubeconfigFilepath,
						},
						State: config.State{
							File: stateFilePath,
						},
						Backup: func() config.Backup {
							if err != nil {
								t.Errorf("%v", err)
							}
							return config.Backup{
								Enabled:   true,
								Directory: tempBackupDir,
							}
						}(),
					}
					return kontextConfig
				}(),
				state: &state.State{
					Group:   state.Group{},
					Context: state.Context{},
					Backup: state.Backup{
						Revisions: []state.Revision{},
					},
				},
			},
			want: func(config *config.Config) *state.State {
				files, err := os.ReadDir(config.Backup.Directory)
				if err != nil {
					t.Errorf("%v", err)
				}

				return &state.State{
					Group:   state.Group{},
					Context: state.Context{},
					Backup: state.Backup{
						Revisions: lo.Map(files, func(item os.DirEntry, index int) state.Revision {
							return state.Revision(filepath.Join(config.Backup.Directory, item.Name()))
						}),
					},
				}
			},
			wantErr: false,
		},
		{
			name: "should skip the backup, as it is disabled",
			args: args{
				config: &config.Config{
					Backup: config.Backup{
						Enabled: false,
					},
				},
				state: &state.State{
					Backup: state.Backup{
						Revisions: []state.Revision{},
					},
				},
			},
			want: func(config *config.Config) *state.State {
				return &state.State{
					Group:   state.Group{},
					Context: state.Context{},
					Backup: state.Backup{
						Revisions: []state.Revision{},
					},
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Reconcile(tt.args.config, tt.args.state)
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got: '%v'", err)
			}

			if !tt.wantErr && !cmp.Equal(tt.want(tt.args.config), tt.args.state) {
				diff := cmp.Diff(tt.want(tt.args.config), tt.args.state)
				t.Errorf("backup.Reconcile() mismatch (-want +got):\n%s", diff)
			}

			// clean up the new backup
			if tt.after != nil {
				tt.after(tt.args.config)
			}
		})
	}
}

func Test_create(t *testing.T) {
	type args struct {
		config *config.Config
	}
	tests := []struct {
		name    string
		after   func(file *os.File)
		args    args
		wantErr bool
	}{
		{
			name: "should create a backup successfully",
			after: func(file *os.File) {
				err := os.Remove(file.Name())
				if err != nil {
					t.Errorf("%v", err)
				}
			},
			args: args{
				config: func() *config.Config {
					_, caller, _, _ := runtime.Caller(0)
					kubeconfigFilepath := filepath.Join(caller, "..", "testdata", "01-valid-kubeconfig.yaml")
					tempBackupDirectory, err := os.MkdirTemp("", "kontext-")

					kontextConfig := &config.Config{
						Global: config.Global{
							Kubeconfig: kubeconfigFilepath,
						},
						Backup: func() config.Backup {
							if err != nil {
								t.Errorf("%v", err)
							}
							return config.Backup{
								Enabled:   true,
								Directory: tempBackupDirectory,
							}
						}(),
					}
					return kontextConfig
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// original is the kubeconfig file, that shall be included within the backup
			original, err := os.Open(tt.args.config.Global.Kubeconfig)
			if err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			// read the original kubeconfig
			want, err := kubeconfig.Read(original)
			if err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			// create the backup
			backup, err := create(tt.args.config)
			if err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			got, err := kubeconfig.Read(backup)
			if err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			// compare both kubeconfig files
			if !tt.wantErr && want != nil && !cmp.Equal(want, got) {
				diff := cmp.Diff(want, got)
				t.Errorf("backup.create() mismatch (-want +got):\n%s", diff)
			}

			// clean up the new backup
			if tt.after != nil {
				tt.after(backup)
			}
		})
	}
}

func Test_computeBackupFilename(t *testing.T) {
	type args struct {
		config *config.Config
	}
	tests := []struct {
		name string
		args args
		want Filename
	}{
		{
			name: `should compute the backup filename that matches '$BACKUP_DIRECTORY/kubeconfig-\d+.yaml`,
			args: args{
				config: &config.Config{
					Backup: config.Backup{
						Directory: "/tmp/kontext",
					},
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeBackupFileName(tt.args.config)

			match, _ := regexp.MatchString(`(.*)kubeconfig-(\d+).yaml`, string(got))

			if !match {
				t.Errorf("backup.computeBackupFilename(), computed file does not match regular expression")
			}
		})
	}
}
