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
	"k8s.io/client-go/tools/clientcmd/api"
)

func Test_create(t *testing.T) {
	type args struct {
		kubeconfig *api.Config
		config     *config.Config
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
