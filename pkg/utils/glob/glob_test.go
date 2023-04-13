package glob

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/samber/lo"
)

func Test_Expand(t *testing.T) {
	type args struct {
		Glob Pattern
	}
	tests := []struct {
		name    string
		args    args
		want    []*os.File
		wantErr bool
	}{
		{
			name: "should build the absolute paths for the glob successfully",
			args: args{
				Glob: func() Pattern {
					_, caller, _, _ := runtime.Caller(0)
					glob := path.Join(caller, "..", "testdata", "*merge*.yaml")
					return Pattern(glob)
				}(),
			},
			want: func() []*os.File {
				var buffer []*os.File

				_, caller, _, _ := runtime.Caller(0)
				kubeconfigFilePathOne := path.Join(caller, "..", "testdata", "01-kontext-merge-1.yaml")
				kubeconfigFilePathTwo := path.Join(caller, "..", "testdata", "02-kontext-merge-2.yaml")
				kubeconfigFilePathThree := path.Join(caller, "..", "testdata", "03-kontext-merge-3.yaml")

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
			got, err := Expand(tt.args.Glob)
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
				t.Errorf("glob.Expand() mismatch (-want +got):\n%s", diff)
			}

		})
	}
}
