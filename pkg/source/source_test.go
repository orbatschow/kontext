package source

import (
	"path"
	"reflect"
	"runtime"
	"testing"

	"github.com/orbatschow/kontext/pkg/config"
)

func Test_Expand(t *testing.T) {
	type args struct {
		Source *config.Source
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "should build the absolute paths for the glob successfully, without duplicates",
			args: args{
				Source: func() *config.Source {
					_, caller, _, _ := runtime.Caller(0)
					include := path.Join(caller, "..", "testdata", "*merge*.yaml")
					exclude := path.Join(caller, "..", "testdata", "*merge-invalid.yaml")

					return &config.Source{
						Include: []string{
							// duplication
							include,
							include,
						},
						Exclude: []string{
							// duplication
							exclude,
							exclude,
						},
					}
				}(),
			},
			want: func() []string {
				var buffer []string

				_, caller, _, _ := runtime.Caller(0)
				mergeFileOne := path.Join(caller, "..", "testdata", "01-kontext-merge-1.yaml")
				mergeFileTwo := path.Join(caller, "..", "testdata", "02-kontext-merge-2.yaml")
				mergeFileThree := path.Join(caller, "..", "testdata", "03-kontext-merge-3.yaml")

				buffer = append(buffer, mergeFileOne)
				buffer = append(buffer, mergeFileTwo)
				buffer = append(buffer, mergeFileThree)

				return buffer
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Expand(tt.args.Source)
			if !tt.wantErr && err != nil {
				t.Errorf("expected error")
			}

			// check if the current context is equal
			if !tt.wantErr && !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want: '%s', got: '%s'", tt.want, got)
			}
		})
	}
}

func Test_computeGlob(t *testing.T) {
	type args struct {
		Glob string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "should build the absolute paths for the glob successfully",
			args: args{
				Glob: func() string {
					_, caller, _, _ := runtime.Caller(0)
					glob := path.Join(caller, "..", "testdata", "*merge*.yaml")
					return glob
				}(),
			},
			want: func() []string {
				var buffer []string

				_, caller, _, _ := runtime.Caller(0)
				mergeFileOne := path.Join(caller, "..", "testdata", "01-kontext-merge-1.yaml")
				mergeFileTwo := path.Join(caller, "..", "testdata", "02-kontext-merge-2.yaml")
				mergeFileThree := path.Join(caller, "..", "testdata", "03-kontext-merge-3.yaml")
				mergeFileFour := path.Join(caller, "..", "testdata", "04-kontext-merge-invalid.yaml")

				buffer = append(buffer, mergeFileOne)
				buffer = append(buffer, mergeFileTwo)
				buffer = append(buffer, mergeFileThree)
				buffer = append(buffer, mergeFileFour)

				return buffer
			}(),
			wantErr: false,
		},
		{
			name: "should not fail on non existing paths",
			args: args{
				Glob: "",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := computeGlob(tt.args.Glob)
			if !tt.wantErr && err != nil {
				t.Errorf("expected error")
			}

			// check if the current context is equal
			if !tt.wantErr && !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want: '%s', got: '%s'", tt.want, got)
			}
		})
	}
}
