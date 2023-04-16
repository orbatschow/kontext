package source

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/samber/lo"
)

func Test_ComputeFiles(t *testing.T) {
	type args struct {
		Source *config.Source
	}
	tests := []struct {
		name    string
		args    args
		want    []*os.File
		wantErr bool
	}{
		{
			name: "should build the absolute paths for the glob successfully, without duplicates",
			args: args{
				Source: func() *config.Source {
					_, caller, _, _ := runtime.Caller(0)
					include := filepath.Join(caller, "..", "testdata", "*merge*.yaml")
					exclude := filepath.Join(caller, "..", "testdata", "*merge-invalid.yaml")

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
			want: func() []*os.File {
				var buffer []*os.File

				_, caller, _, _ := runtime.Caller(0)
				kubeconfigFilePathOne := filepath.Join(caller, "..", "testdata", "01-kontext-merge-1.yaml")
				kubeconfigFilePathTwo := filepath.Join(caller, "..", "testdata", "02-kontext-merge-2.yaml")
				kubeconfigFilePathThree := filepath.Join(caller, "..", "testdata", "03-kontext-merge-3.yaml")

				kubeconfigFileOne, err := os.Open(kubeconfigFilePathOne)
				if err != nil {
					t.Errorf("%v", err)
				}
				buffer = append(buffer, kubeconfigFileOne)

				kubeconfigFileTwo, err := os.Open(kubeconfigFilePathTwo)
				if err != nil {
					t.Errorf("%v", err)
				}
				buffer = append(buffer, kubeconfigFileTwo)

				kubeconfigFileThree, err := os.Open(kubeconfigFilePathThree)
				if err != nil {
					t.Errorf("%v", err)
				}
				buffer = append(buffer, kubeconfigFileThree)

				return buffer
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ComputeFiles(tt.args.Source)
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			if tt.wantErr && err == nil {
				t.Errorf("expected error, err: '%v'", err)
			}

			gotNames := lo.Map(got, func(item *os.File, index int) string {
				return item.Name()
			})
			wantNames := lo.Map(tt.want, func(item *os.File, index int) string {
				return item.Name()
			})

			if !tt.wantErr && !cmp.Equal(wantNames, gotNames) {
				diff := cmp.Diff(wantNames, gotNames)
				t.Errorf("source.ComputeFiles() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
