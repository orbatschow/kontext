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
		config func(directory string) *config.Config
		state  *state.State
	}
	tests := []struct {
		name    string
		before  func() (string, error)
		after   func(directory string) error
		args    args
		want    func(directory string) *state.State
		wantErr bool
	}{
		{
			name: "should create a backup and add a new revision to the state",
			before: func() (string, error) {
				// create temporary working directory
				directory, err := os.MkdirTemp("", "kontext-")
				if err != nil {
					return "", err
				}
				// create backup directory
				err = os.Mkdir(filepath.Join(directory, "backup"), 0777)
				if err != nil {
					return "", err
				}
				// create state file
				stateFilePath := filepath.Join(directory, "state.json")
				_, err = os.Create(stateFilePath)
				if err != nil {
					t.Errorf("%v", err)
				}

				return directory, nil
			},
			after: func(directory string) error {
				err := os.RemoveAll(directory)
				if err != nil {
					return err
				}
				return nil
			},
			args: args{
				config: func(directory string) *config.Config {
					// current working directory
					_, caller, _, _ := runtime.Caller(0)

					kontextConfig := &config.Config{
						Global: config.Global{
							Kubeconfig: filepath.Join(caller, "..", "testdata", "01-valid-kubeconfig.yaml"),
						},
						State: config.State{
							File: filepath.Join(directory, "state.json"),
						},
						Backup: config.Backup{
							Enabled:   true,
							Directory: filepath.Join(directory, "backup"),
						},
					}
					return kontextConfig
				},
				state: &state.State{
					Group:   state.Group{},
					Context: state.Context{},
					Backup: state.Backup{
						Revisions: []state.Revision{},
					},
				},
			},
			want: func(tempDirectory string) *state.State {
				files, err := os.ReadDir(filepath.Join(tempDirectory, "backup"))
				if err != nil {
					t.Errorf("%v", err)
				}

				return &state.State{
					Group:   state.Group{},
					Context: state.Context{},
					Backup: state.Backup{
						Revisions: lo.Map(files, func(item os.DirEntry, index int) state.Revision {
							return state.Revision(filepath.Join(filepath.Join(tempDirectory, "backup"), item.Name()))
						}),
					},
				}
			},
			wantErr: false,
		},
		{
			name: "should skip the backup, as it is disabled",
			before: func() (string, error) {
				return "", nil
			},
			after: func(directory string) error {
				return nil
			},
			args: args{
				config: func(directory string) *config.Config {
					return &config.Config{
						Backup: config.Backup{
							Enabled: false,
						},
					}
				},
				state: &state.State{
					Backup: state.Backup{
						Revisions: []state.Revision{},
					},
				},
			},
			want: func(tempDirectory string) *state.State {
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
			tempDirectory, err := tt.before()
			if err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			reconciler := Reconciler{
				Config: tt.args.config(tempDirectory),
				State:  tt.args.state,
			}
			err = reconciler.Reconcile()
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got: '%v'", err)
			}

			if !tt.wantErr && !cmp.Equal(tt.want(tempDirectory), tt.args.state) {
				diff := cmp.Diff(tt.want(tempDirectory), tt.args.state)
				t.Errorf("backup.Reconcile() mismatch (-want +got):\n%s", diff)
			}

			// clean up the temporary working directory
			err = tt.after(tempDirectory)
			if err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
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
			reconciler := Reconciler{
				Config: tt.args.config,
				State:  nil,
			}
			backup, err := reconciler.create()
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
