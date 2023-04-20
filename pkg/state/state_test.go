package state

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/samber/lo"
)

func Test_initialize(t *testing.T) {
	type args struct {
		Config *config.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should create the state file from the config successfully",
			args: args{
				Config: &config.Config{
					State: config.State{
						File: func() string {
							tempDir := os.TempDir()
							seed := lo.RandomString(10, lo.LowerCaseLettersCharset)
							targetDir := filepath.Join(tempDir, seed)
							t.Logf("temporary state directory: '%s'", targetDir)
							return filepath.Join(targetDir, "state.json")
						}(),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Init(tt.args.Config)
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got: '%v'", err)
			}

			if _, err := os.Stat(tt.args.Config.State.File); errors.Is(err, os.ErrNotExist) {
				t.Errorf("state file does not exist")
			}

			dir, _ := filepath.Split(tt.args.Config.State.File)
			err = os.RemoveAll(dir)
			if !tt.wantErr && err != nil {
				t.Errorf("could not remove state directory, err: '%v'", err)
			}
		})
	}
}

func Test_Read(t *testing.T) {
	type args struct {
		Config *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *State
		wantErr bool
	}{
		{
			name: "should read the state successfully",
			args: args{
				Config: &config.Config{
					State: config.State{
						File: func() string {
							_, caller, _, _ := runtime.Caller(0)
							stateFile := filepath.Join(caller, "..", "testdata", "01-valid-state.json")
							return stateFile
						}(),
					},
				},
			},
			want: &State{
				Group: Group{
					Active: "dev",
					History: []History{
						"local",
						"dev",
					},
				},
				Context: Context{
					Active: "kind-dev",
					History: []History{
						"kind-dev",
						"kind-local",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should read an empty state successfully",
			args: args{
				Config: &config.Config{
					State: config.State{
						File: func() string {
							_, caller, _, _ := runtime.Caller(0)
							stateFile := filepath.Join(caller, "..", "testdata", "02-valid-empty-state.json")
							return stateFile
						}(),
					},
				},
			},
			want:    &State{},
			wantErr: false,
		},
		{
			name: "should throw an error due to invalid file",
			args: args{
				Config: &config.Config{
					State: config.State{
						File: func() string {
							_, caller, _, _ := runtime.Caller(0)
							stateFile := filepath.Join(caller, "..", "testdata", "03-invalid-state.json")
							return stateFile
						}(),
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Read(tt.args.Config)
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			if tt.wantErr && err == nil {
				t.Errorf("expected error, err: '%v'", err)
			}

			if !tt.wantErr && !cmp.Equal(tt.want, got) {
				diff := cmp.Diff(&tt.want, got)
				t.Errorf("state.Read() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_Write(t *testing.T) {
	type args struct {
		Config *config.Config
		State  *State
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "should write the state successfully",
			args: args{
				Config: &config.Config{
					State: config.State{
						File: func() string {
							tempDir := os.TempDir()
							seed := lo.RandomString(10, lo.LowerCaseLettersCharset)
							targetDir := filepath.Join(tempDir, seed)
							t.Logf("temporary state directory: '%s'", targetDir)
							return filepath.Join(targetDir, "state.json")
						}(),
					},
				},
				State: &State{
					Group: Group{
						Active:  "dev",
						History: nil,
					},
					Context: Context{
						Active:  "kind-dev",
						History: nil,
					},
				},
			},
			want:    []byte(`{"group":{"active":"dev"},"context":{"active":"kind-dev"},"backup":{}}`),
			wantErr: false,
		},
		{
			name: "should write the state successfully, even without a context",
			args: args{
				Config: &config.Config{
					State: config.State{
						File: func() string {
							tempDir := os.TempDir()
							seed := lo.RandomString(10, lo.LowerCaseLettersCharset)
							targetDir := filepath.Join(tempDir, seed)
							t.Logf("temporary state directory: '%s'", targetDir)
							return filepath.Join(targetDir, "state.json")
						}(),
					},
				},
				State: &State{
					Group: Group{
						Active: "dev",
						History: []History{
							"dev",
						},
					},
				},
			},
			want:    []byte(`{"group":{"active":"dev","history":["dev"]},"context":{},"backup":{}}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// init the state
			err := Init(tt.args.Config)
			// log every error here, as we do not want to test init
			if err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			// write the state to the target file
			err = Write(tt.args.Config, tt.args.State)
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got: '%v'", err)
			}

			// compare written file and want
			got, err := os.ReadFile(tt.args.Config.State.File)
			if err != nil {
				t.Errorf("%v", err)
			}
			if !tt.wantErr && !bytes.Equal(tt.want, got) {
				t.Errorf("%v", err)
			}

			// cleanup
			dir, _ := filepath.Split(tt.args.Config.State.File)
			err = os.RemoveAll(dir)
			if !tt.wantErr && err != nil {
				t.Errorf("could not remove state directory, err: '%v'", err)
			}
		})
	}
}

func Test_ComputeHistory(t *testing.T) {
	type args struct {
		Config  *config.Config
		Entry   History
		History []History
	}
	var tests = []struct {
		name string
		args args
		want []History
	}{
		{
			name: "should add a new entry to the history successfully",
			args: args{
				Config: &config.Config{
					State: config.State{
						History: config.History{
							Size: DefaultMaximumHistorySize,
						},
					},
				},
				Entry: "local",
				History: []History{
					"dev",
					"prod",
				},
			},
			want: []History{
				"dev",
				"prod",
				"local",
			},
		},
		{
			name: "should skip adding a duplicate entry history",
			args: args{
				Config: &config.Config{
					State: config.State{
						History: config.History{
							Size: DefaultMaximumHistorySize,
						},
					},
				},
				Entry: "dev",
				History: []History{
					"dev",
				},
			},
			want: []History{
				"dev",
			},
		},
		{
			name: "should add an entry to an empty history",
			args: args{
				Config: &config.Config{
					State: config.State{
						History: config.History{
							Size: DefaultMaximumHistorySize,
						},
					},
				},
				Entry:   "dev",
				History: []History{},
			},
			want: []History{
				"dev",
			},
		},
		{
			name: "should not exceed the default maximum history size",
			args: args{
				Config: &config.Config{
					State: config.State{
						History: config.History{
							Size: DefaultMaximumHistorySize,
						},
					},
				},
				Entry: "private",
				History: []History{
					"dev",
					"local",
					"prod",
					"dev",
					"local",
					"prod",
					"dev",
					"local",
					"prod",
					"dev",
				},
			},
			want: []History{
				"local",
				"prod",
				"dev",
				"local",
				"prod",
				"dev",
				"local",
				"prod",
				"dev",
				"private",
			},
		},
		{
			name: "should not exceed the maximum history size, that was defined by the user",
			args: args{
				Config: &config.Config{
					State: config.State{
						History: config.History{
							Size: 3,
						},
					},
				},
				Entry: "private",
				History: []History{
					"dev",
					"local",
					"prod",
				},
			},
			want: []History{
				"local",
				"prod",
				"private",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeHistory(tt.args.Config, tt.args.Entry, tt.args.History)

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want: '%v', got: '%v'", tt.want, tt.args.Config)
			}
		})
	}
}
