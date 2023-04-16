package backup

import (
	"regexp"
	"testing"

	"github.com/orbatschow/kontext/pkg/config"
)

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
