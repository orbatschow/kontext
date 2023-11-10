package revision

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/samber/lo"
)

func before(revisionCount int) (string, error) {
	// create temporary working directory
	directory, err := os.MkdirTemp("", "kontext-")
	if err != nil {
		return "", err
	}
	// create backup directory
	backupDirectory := filepath.Join(directory, "backup")
	err = os.Mkdir(backupDirectory, 0777)
	if err != nil {
		return "", err
	}

	_, err = generateTestRevisions(backupDirectory, revisionCount)
	if err != nil {
		return "", err
	}

	return directory, nil
}

func generateTestRevisions(directory string, count int) ([]*os.File, error) {
	var buffer []*os.File

	for i := 0; i < count; i++ {
		timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		file, err := os.Create(filepath.Join(directory, timestamp))
		if err != nil {
			return nil, err
		}

		buffer = append(buffer, file)
	}

	return buffer, nil
}

func generateTestBackupRevision(directory string) (*os.File, error) {
	file, err := os.CreateTemp(directory, "kubeconfig-")
	if err != nil {
		return nil, err
	}

	return file, nil
}

func after(directory string) error {
	err := os.RemoveAll(directory)
	if err != nil {
		return err
	}
	return nil
}

func Test_reconcileRevisions(t *testing.T) {
	type args struct {
		revisionCount int
		backupFile    func(directory string) (*os.File, error)
		config        func(directory string) *config.Config
		state         func(directory string) *state.State
	}
	tests := []struct {
		name    string
		before  func(revisionCount int) (string, error)
		after   func(directory string) error
		args    args
		want    func(config *config.Config, revisions []state.Revision, backupFile *os.File) []state.Revision
		wantErr bool
	}{
		{
			name:   "should add a new revision to a state with no existing revisions",
			before: before,
			after:  after,
			args: args{
				revisionCount: 0,
				backupFile:    generateTestBackupRevision,
				config: func(directory string) *config.Config {
					return &config.Config{
						Backup: config.Backup{
							Enabled:   true,
							Revisions: 10,
							Directory: filepath.Join(directory, "backup"),
						},
					}
				},
				state: func(directory string) *state.State {
					files, err := os.ReadDir(filepath.Join(directory, "backup"))
					if err != nil {
						t.Errorf("%v", err)
					}

					return &state.State{
						Group:   state.Group{},
						Context: state.Context{},
						Backup: state.Backup{
							Revisions: lo.Map(files, func(item os.DirEntry, index int) state.Revision {
								return state.Revision(filepath.Join(filepath.Join(directory, "backup"), item.Name()))
							}),
						},
					}
				},
			},
			want: func(config *config.Config, revisions []state.Revision, backupFile *os.File) []state.Revision {
				revisions = append(revisions, state.Revision(backupFile.Name()))
				return lo.Slice(revisions, len(revisions)-config.Backup.Revisions, len(revisions))
			},
			wantErr: false,
		},
		{
			name:   "should add a new revision and ignore all existing revisions",
			before: before,
			after:  after,
			args: args{
				revisionCount: 3,
				backupFile:    generateTestBackupRevision,
				config: func(directory string) *config.Config {
					return &config.Config{
						Backup: config.Backup{
							Enabled:   true,
							Revisions: 10,
							Directory: filepath.Join(directory, "backup"),
						},
					}
				},
				state: func(directory string) *state.State {
					files, err := os.ReadDir(filepath.Join(directory, "backup"))
					if err != nil {
						t.Errorf("%v", err)
					}

					return &state.State{
						Group:   state.Group{},
						Context: state.Context{},
						Backup: state.Backup{
							Revisions: lo.Map(files, func(item os.DirEntry, index int) state.Revision {
								return state.Revision(filepath.Join(filepath.Join(directory, "backup"), item.Name()))
							}),
						},
					}
				},
			},
			want: func(config *config.Config, revisions []state.Revision, backupFile *os.File) []state.Revision {
				revisions = append(revisions, state.Revision(backupFile.Name()))
				return lo.Slice(revisions, len(revisions)-config.Backup.Revisions, len(revisions))
			},
			wantErr: false,
		},

		{
			name:   "should add a new revision and remove one revision, that exceeds the limit",
			before: before,
			after:  after,
			args: args{
				revisionCount: 10,
				backupFile:    generateTestBackupRevision,
				config: func(directory string) *config.Config {
					return &config.Config{
						Backup: config.Backup{
							Enabled:   true,
							Revisions: 10,
							Directory: filepath.Join(directory, "backup"),
						},
					}
				},
				state: func(directory string) *state.State {
					files, err := os.ReadDir(filepath.Join(directory, "backup"))
					if err != nil {
						t.Errorf("%v", err)
					}

					return &state.State{
						Group:   state.Group{},
						Context: state.Context{},
						Backup: state.Backup{
							Revisions: lo.Map(files, func(item os.DirEntry, index int) state.Revision {
								return state.Revision(filepath.Join(filepath.Join(directory, "backup"), item.Name()))
							}),
						},
					}
				},
			},
			want: func(config *config.Config, revisions []state.Revision, backupFile *os.File) []state.Revision {
				revisions = append(revisions, state.Revision(backupFile.Name()))
				return lo.Slice(revisions, len(revisions)-config.Backup.Revisions, len(revisions))
			},
			wantErr: false,
		},
		{
			name:   "should add a new revision and remove multiple revisions, that exceed the limit",
			before: before,
			after:  after,
			args: args{
				revisionCount: 50,
				backupFile:    generateTestBackupRevision,
				config: func(directory string) *config.Config {
					return &config.Config{
						Backup: config.Backup{
							Enabled:   true,
							Revisions: 10,
							Directory: filepath.Join(directory, "backup"),
						},
					}
				},
				state: func(directory string) *state.State {
					files, err := os.ReadDir(filepath.Join(directory, "backup"))
					if err != nil {
						t.Errorf("%v", err)
					}

					return &state.State{
						Group:   state.Group{},
						Context: state.Context{},
						Backup: state.Backup{
							Revisions: lo.Map(files, func(item os.DirEntry, index int) state.Revision {
								return state.Revision(filepath.Join(filepath.Join(directory, "backup"), item.Name()))
							}),
						},
					}
				},
			},
			want: func(config *config.Config, revisions []state.Revision, backupFile *os.File) []state.Revision {
				revisions = append(revisions, state.Revision(backupFile.Name()))
				return lo.Slice(revisions, len(revisions)-config.Backup.Revisions, len(revisions))
			},
			wantErr: false,
		},
		{
			name:   "should remove all backup revisions",
			before: before,
			after:  after,
			args: args{
				revisionCount: 5,
				backupFile:    generateTestBackupRevision,
				config: func(directory string) *config.Config {
					return &config.Config{
						Backup: config.Backup{
							Enabled:   true,
							Revisions: 0,
							Directory: filepath.Join(directory, "backup"),
						},
					}
				},
				state: func(directory string) *state.State {
					files, err := os.ReadDir(filepath.Join(directory, "backup"))
					if err != nil {
						t.Errorf("%v", err)
					}

					return &state.State{
						Group:   state.Group{},
						Context: state.Context{},
						Backup: state.Backup{
							Revisions: lo.Map(files, func(item os.DirEntry, index int) state.Revision {
								return state.Revision(filepath.Join(filepath.Join(directory, "backup"), item.Name()))
							}),
						},
					}
				},
			},
			want: func(config *config.Config, revisions []state.Revision, backupFile *os.File) []state.Revision {
				revisions = append(revisions, state.Revision(backupFile.Name()))
				return lo.Slice(revisions, len(revisions)-config.Backup.Revisions, len(revisions))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDirectory, err := tt.before(tt.args.revisionCount)
			if err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			// create a new revision
			backupFile, err := tt.args.backupFile(tempDirectory)
			if err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			reconciler := Reconciler{
				Config: tt.args.config(tempDirectory),
				State:  tt.args.state(tempDirectory),
				Backup: backupFile,
			}
			got, err := reconciler.Reconcile()
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got: '%v'", err)
			}

			want := tt.want(reconciler.Config, reconciler.State.Backup.Revisions, backupFile)
			if !tt.wantErr && !cmp.Equal(want, got) {
				diff := cmp.Diff(want, got)
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
