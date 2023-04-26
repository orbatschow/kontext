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
		SourceItem *config.SourceItem
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
				SourceItem: func() *config.SourceItem {
					_, caller, _, _ := runtime.Caller(0)
					include := filepath.Join(caller, "..", "testdata", "*merge*.yaml")
					exclude := filepath.Join(caller, "..", "testdata", "*merge-invalid.yaml")

					return &config.SourceItem{
						Name: "dev",
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
				kubeconfigs := []string{
					filepath.Join(caller, "..", "testdata", "01-kontext-merge-1.yaml"),
					filepath.Join(caller, "..", "testdata", "02-kontext-merge-2.yaml"),
					filepath.Join(caller, "..", "testdata", "03-kontext-merge-3.yaml"),
				}

				for _, kubeconfig := range kubeconfigs {
					file, err := os.Open(kubeconfig)
					if err != nil {
						t.Errorf("%v", err)
					}
					buffer = append(buffer, file)
				}

				return buffer
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ComputeFiles(tt.args.SourceItem)
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

func Test_computeIncludes(t *testing.T) {
	type args struct {
		SourceItem *config.SourceItem
	}
	tests := []struct {
		name    string
		args    args
		want    []*os.File
		wantErr bool
	}{
		{
			name: "should include glob successfully and remove all duplicates",
			args: args{
				SourceItem: func() *config.SourceItem {
					_, caller, _, _ := runtime.Caller(0)
					include := filepath.Join(caller, "..", "testdata", "*merge*.yaml")

					return &config.SourceItem{
						Name: "dev",
						Include: []string{
							// duplication
							include,
							include,
						},
					}
				}(),
			},
			want: func() []*os.File {
				var buffer []*os.File

				_, caller, _, _ := runtime.Caller(0)
				kubeconfigs := []string{
					filepath.Join(caller, "..", "testdata", "01-kontext-merge-1.yaml"),
					filepath.Join(caller, "..", "testdata", "02-kontext-merge-2.yaml"),
					filepath.Join(caller, "..", "testdata", "03-kontext-merge-3.yaml"),
					filepath.Join(caller, "..", "testdata", "04-kontext-merge-invalid.yaml"),
				}

				for _, kubeconfig := range kubeconfigs {
					file, err := os.Open(kubeconfig)
					if err != nil {
						t.Errorf("%v", err)
					}
					buffer = append(buffer, file)
				}

				return buffer
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := computeIncludes(tt.args.SourceItem)
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
				t.Errorf("source.computeIncludes() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_computeExcludes(t *testing.T) {
	type args struct {
		SourceItem *config.SourceItem
	}
	tests := []struct {
		name    string
		args    args
		want    []*os.File
		wantErr bool
	}{
		{
			name: "should include glob successfully and remove all duplicates",
			args: args{
				SourceItem: func() *config.SourceItem {
					_, caller, _, _ := runtime.Caller(0)
					include := filepath.Join(caller, "..", "testdata", "*invalid*.yaml")

					return &config.SourceItem{
						Name: "dev",
						Include: []string{
							// duplication
							include,
							include,
						},
					}
				}(),
			},
			want: func() []*os.File {
				var buffer []*os.File

				_, caller, _, _ := runtime.Caller(0)
				kubeconfigs := []string{
					filepath.Join(caller, "..", "testdata", "04-kontext-merge-invalid.yaml"),
				}

				for _, kubeconfig := range kubeconfigs {
					file, err := os.Open(kubeconfig)
					if err != nil {
						t.Errorf("%v", err)
					}
					buffer = append(buffer, file)
				}

				return buffer
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := computeIncludes(tt.args.SourceItem)
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
				t.Errorf("source.computeIncludes() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_difference(t *testing.T) {
	type args struct {
		SourceItem *config.SourceItem
	}
	tests := []struct {
		name    string
		args    args
		want    []*os.File
		wantErr bool
	}{
		{
			name: "should build the difference for the given includes and excludes",
			args: args{
				SourceItem: func() *config.SourceItem {
					_, caller, _, _ := runtime.Caller(0)
					include := filepath.Join(caller, "..", "testdata", "*merge*.yaml")
					exclude := filepath.Join(caller, "..", "testdata", "*merge-invalid.yaml")

					return &config.SourceItem{
						Name: "dev",
						Include: []string{
							include,
						},
						Exclude: []string{
							exclude,
						},
					}
				}(),
			},
			want: func() []*os.File {
				var buffer []*os.File

				_, caller, _, _ := runtime.Caller(0)
				kubeconfigs := []string{
					filepath.Join(caller, "..", "testdata", "01-kontext-merge-1.yaml"),
					filepath.Join(caller, "..", "testdata", "02-kontext-merge-2.yaml"),
					filepath.Join(caller, "..", "testdata", "03-kontext-merge-3.yaml"),
				}

				for _, kubeconfig := range kubeconfigs {
					file, err := os.Open(kubeconfig)
					if err != nil {
						t.Errorf("%v", err)
					}
					buffer = append(buffer, file)
				}

				return buffer
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			includes, err := computeIncludes(tt.args.SourceItem)
			if err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}
			excludes, err := computeExcludes(tt.args.SourceItem)
			if err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			if tt.wantErr && err == nil {
				t.Errorf("expected error, err: '%v'", err)
			}

			got, _ := difference(includes, excludes)

			gotNames := lo.Map(got, func(item *os.File, index int) string {
				return item.Name()
			})
			wantNames := lo.Map(tt.want, func(item *os.File, index int) string {
				return item.Name()
			})

			if !tt.wantErr && !cmp.Equal(wantNames, gotNames) {
				diff := cmp.Diff(wantNames, gotNames)
				t.Errorf("source.computeIncludes() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
