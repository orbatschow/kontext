package state

import (
	"bytes"
	"errors"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/orbatschow/kontext/pkg/config"
)

var charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(length int) string {
	buffer := make([]byte, length)
	for i := range buffer {
		buffer[i] = charset[rand.Intn(len(charset))]
	}
	return string(buffer)
}

func Test_computeStateFileLocation(t *testing.T) {
	type args struct {
		Config *config.Config
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should return the default location",
			args: args{
				Config: &config.Config{
					State: config.State{},
				},
			},
			want: defaultStateFile,
		},

		{
			name: "should return the default location",
			args: args{
				Config: &config.Config{
					State: config.State{
						Location: "test-path/test.json",
					},
				},
			},
			want: "test-path/test.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeStateFileLocation(tt.args.Config)

			if tt.want != got {
				t.Errorf("state.computeStateFileLocation(), want: '%s', got: '%s'", tt.want, got)
			}
		})
	}
}

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
						Location: func() string {
							tempDir := os.TempDir()
							seed := randomString(10)
							targetDir := path.Join(tempDir, seed)
							t.Logf("temporary state directory: '%s'", targetDir)
							return path.Join(targetDir, "state.json")
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
				t.Errorf("expected error")
			}

			if _, err := os.Stat(tt.args.Config.State.Location); errors.Is(err, os.ErrNotExist) {
				t.Errorf("state file does not exist")
			}

			dir, _ := filepath.Split(tt.args.Config.State.Location)
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
						Location: func() string {
							_, caller, _, _ := runtime.Caller(0)
							stateFile := path.Join(caller, "..", "testdata", "01-valid-state.json")
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
						Location: func() string {
							_, caller, _, _ := runtime.Caller(0)
							stateFile := path.Join(caller, "..", "testdata", "02-valid-empty-state.json")
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
						Location: func() string {
							_, caller, _, _ := runtime.Caller(0)
							stateFile := path.Join(caller, "..", "testdata", "03-invalid-state.json")
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
			t.Logf("%s", tt.args.Config.State.Location)
			got, err := Read(tt.args.Config)
			if !tt.wantErr && err != nil {
				t.Errorf("expected error")
			}
			t.Logf("%v", got)
			if !tt.wantErr && !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want: '%v', got: '%v'", tt.want, tt.args.Config)
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
						Location: func() string {
							tempDir := os.TempDir()
							seed := randomString(10)
							targetDir := path.Join(tempDir, seed)
							t.Logf("temporary state directory: '%s'", targetDir)
							return path.Join(targetDir, "state.json")
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
			want:    []byte(`{"group":{"active":"dev"},"context":{"active":"kind-dev"}}`),
			wantErr: false,
		},
		{
			name: "should write the state successfully, even without a context",
			args: args{
				Config: &config.Config{
					State: config.State{
						Location: func() string {
							tempDir := os.TempDir()
							seed := randomString(10)
							targetDir := path.Join(tempDir, seed)
							t.Logf("temporary state directory: '%s'", targetDir)
							return path.Join(targetDir, "state.json")
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
			want:    []byte(`{"group":{"active":"dev","history":["dev"]},"context":{}}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// init the state
			err := Init(tt.args.Config)
			if !tt.wantErr && err != nil {
				t.Errorf("expected error")
			}

			// write the state to the target file
			err = Write(tt.args.Config, tt.args.State)
			if !tt.wantErr && err != nil {
				t.Errorf("expected error")
			}

			// compare written file and want
			got, err := os.ReadFile(tt.args.Config.State.Location)
			if err != nil {
				t.Errorf("%v", err)
			}
			if !tt.wantErr && !bytes.Equal(tt.want, got) {
				t.Errorf("%v", err)
			}

			// cleanup
			dir, _ := filepath.Split(tt.args.Config.State.Location)
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
				Config: &config.Config{},
				Entry:  "local",
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
				Config: &config.Config{},
				Entry:  "dev",
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
				Config:  &config.Config{},
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
					State: config.State{},
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
							Size: &[]int{3}[0],
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
