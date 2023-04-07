package backup

import (
	"os"
	"testing"
)

func Test_validate(t *testing.T) {
	type args struct {
		Config *Config
	}
	tests := []struct {
		name    string
		before  func(t *testing.T)
		args    args
		wantErr bool
	}{
		{
			name: "should validate the kontext configuration successfully",
			args: args{
				Config: &Config{
					Global: Global{
						Kubeconfig: "fake-path",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should validate the kontext configuration successfully",
			before: func(t *testing.T) {
				t.Setenv("KUBECONFIG", "fake-kubeconfig")
			},
			args: args{
				Config: &Config{},
			},
			wantErr: false,
		},
		{
			name: "should throw an error, because no kubeconfig is set",
			before: func(t *testing.T) {
				err := os.Unsetenv("KUBECONFIG")
				if err != nil {
					t.Errorf("%v", err)
				}
			},
			args: args{
				Config: &Config{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}

			err := validate(tt.args.Config)
			if tt.wantErr == false && err != nil {
				t.Errorf("config.validate() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
