package backup

import (
	"testing"

	"github.com/orbatschow/kontext/pkg/config"
)

func Test_Create(t *testing.T) {
	type args struct {
		Config *config.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should skip the backup, as it is disabled",
			args: args{
				Config: &config.Config{
					Backup: config.Backup{
						Enabled:  true,
						Location: "",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Create(tt.args.Config)
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got: '%v'", err)
			}

			// load kubeconfig from file and compare with local file

		})
	}
}
