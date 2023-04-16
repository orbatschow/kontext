package group

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/state"
	"k8s.io/client-go/tools/clientcmd/api"
)

func Test_Get(t *testing.T) {
	type args struct {
		GroupName string
		Config    *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *config.Group
		wantErr bool
	}{
		{
			name: "should get a group successfully",
			args: args{
				GroupName: "dev",
				Config: &config.Config{
					Groups: []config.Group{
						{
							Name:    "dev",
							Context: "kind-dev",
							Sources: []string{
								"dev",
								"prod",
							},
						},
					},
				},
			},
			want: &config.Group{
				Name:    "dev",
				Context: "kind-dev",
				Sources: []string{
					"dev",
					"prod",
				},
			},
			wantErr: false,
		},
		{
			name: "should return an error due to missing group",
			args: args{
				GroupName: "dev",
				Config: &config.Config{
					Groups: []config.Group{},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := Client{
				Config: tt.args.Config,
			}

			got, err := client.Get(tt.args.GroupName)
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error, err: '%s'", err)
			}

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got: '%v'", err)
			}

			if !tt.wantErr && !cmp.Equal(tt.want, got) {
				diff := cmp.Diff(&tt.want, &client.State)
				t.Errorf("group.Set() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_Set(t *testing.T) {
	type args struct {
		GroupName string
		Config    *config.Config
		State     *state.State
		APIConfig *api.Config
	}
	tests := []struct {
		name string
		args args
		want struct {
			APIConfig *api.Config
			State     *state.State
		}
		wantErr bool
	}{
		{
			name: "should change the state to the given group",
			args: args{
				GroupName: "dev",
				Config: &config.Config{
					Groups: []config.Group{
						{
							Name: "dev",
							Sources: []string{
								"dev",
							},
						},
					},
					Sources: []config.Source{
						{
							Name:    "dev",
							Include: nil,
							Exclude: nil,
						},
					},
				},
				State: &state.State{},
			},
			want: struct {
				APIConfig *api.Config
				State     *state.State
			}{
				APIConfig: &api.Config{},
				State: &state.State{
					Group: state.Group{
						Active: "dev",
						History: []state.History{
							"dev",
						},
					},
					Context: state.Context{},
				},
			},
		},
		{
			name: "should throw an error due to non existing group",
			args: args{
				GroupName: "kind",
				Config:    &config.Config{},
				State:     &state.State{},
			},
			want: struct {
				APIConfig *api.Config
				State     *state.State
			}{
				APIConfig: nil,
				State:     nil,
			},
			wantErr: true,
		},
		{
			name: "should change the state to the given group and set the default context",
			args: args{
				GroupName: "dev",
				Config: &config.Config{
					Groups: []config.Group{
						{
							Name:    "dev",
							Context: "kind-dev",
							Sources: []string{
								"dev",
							},
						},
					},
					Sources: []config.Source{
						{
							Name: "dev",
							Include: func() []string {
								var buffer []string
								_, caller, _, _ := runtime.Caller(0)
								kubeConfigFile := filepath.Join(caller, "..", "testdata", "01-valid-kubeconfig.yaml")

								buffer = append(buffer, kubeConfigFile)

								return buffer
							}(),
							Exclude: nil,
						},
					},
				},
				State: &state.State{},
			},
			want: struct {
				APIConfig *api.Config
				State     *state.State
			}{
				APIConfig: &api.Config{
					CurrentContext: "kind-dev",
				},
				State: &state.State{
					Group: state.Group{
						Active: "dev",
						History: []state.History{
							"dev",
						},
					},
					Context: state.Context{
						Active: "kind-dev",
						History: []state.History{
							"kind-dev",
						},
					},
				},
			},
		},
		{
			name: "should throw an error, as default context does not exist",
			args: args{
				GroupName: "dev",
				Config: &config.Config{
					Groups: []config.Group{
						{
							Name:    "dev",
							Context: "invalid",
							Sources: []string{
								"dev",
							},
						},
					},
					Sources: []config.Source{
						{
							Name: "dev",
							Include: func() []string {
								var buffer []string
								_, caller, _, _ := runtime.Caller(0)
								kubeConfigFile := filepath.Join(caller, "..", "testdata", "01-valid-kubeconfig.yaml")

								buffer = append(buffer, kubeConfigFile)

								return buffer
							}(),
							Exclude: nil,
						},
					},
				},
				State: &state.State{},
			},
			want: struct {
				APIConfig *api.Config
				State     *state.State
			}{
				APIConfig: nil,
				State:     nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := Client{
				Config:    tt.args.Config,
				State:     tt.args.State,
				APIConfig: tt.args.APIConfig,
			}

			err := client.Set(tt.args.GroupName)
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got: '%v'", err)
			}

			if !tt.wantErr && !cmp.Equal(tt.want.State, client.State) {
				diff := cmp.Diff(&tt.want.State, &client.State)
				t.Errorf("group.Set() state mismatch (-want +got):\n%s", diff)
			}

			ignored := cmpopts.IgnoreFields(api.Config{},
				"Preferences",
				"Clusters",
				"AuthInfos",
				"Contexts",
				"Extensions",
			)

			if !tt.wantErr && tt.want.APIConfig != nil && !cmp.Equal(tt.want.APIConfig, client.APIConfig, ignored) {
				diff := cmp.Diff(&tt.want.APIConfig, &client.APIConfig, ignored)
				t.Errorf("group.Set() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_Reload(t *testing.T) {
	type args struct {
		GroupName string
		Config    *config.Config
		State     *state.State
		APIConfig *api.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *state.State
		wantErr bool
	}{
		{
			name: "should reload the current group without an error",
			args: args{
				GroupName: "dev",
				Config: &config.Config{
					Groups: []config.Group{
						{
							Name: "dev",
							Sources: []string{
								"dev",
							},
						},
					},
					Sources: []config.Source{
						{
							Name:    "dev",
							Include: nil,
							Exclude: nil,
						},
					},
				},
				State: &state.State{
					Group: state.Group{
						Active: "dev",
					},
				},
			},
			want: &state.State{
				Group: state.Group{
					Active: "dev",
					History: []state.History{
						"dev",
					},
				},
				Context: state.Context{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := Client{
				Config:    tt.args.Config,
				State:     tt.args.State,
				APIConfig: tt.args.APIConfig,
			}

			err := client.Reload()
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got: '%v'", err)
			}

			if !tt.wantErr && !cmp.Equal(tt.want, client.State) {
				diff := cmp.Diff(&tt.want, &client.State)
				t.Errorf("group.Set() state mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_Print(t *testing.T) {
	type args struct {
		State *state.State
		Group []config.Group
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should print without an error",
			args: args{
				State: &state.State{
					Group: state.Group{
						Active: "dev",
					},
				},
				Group: []config.Group{
					{
						Name:    "dev",
						Context: "kind-dev",
						Sources: []string{
							"dev",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should print without an error, even when there is no group",
			args: args{
				State: &state.State{},
				Group: []config.Group{},
			},
			wantErr: false,
		},
		{
			name: "should print without an error, even when there is no active group",
			args: args{
				State: &state.State{},
				Group: []config.Group{
					{
						Name: "dev",
						Sources: []string{
							"dev",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should print without an error, even when there is no source",
			args: args{
				State: &state.State{},
				Group: []config.Group{
					{
						Name: "dev",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := Client{
				State: tt.args.State,
			}

			err := client.Print(tt.args.Group...)
			if !tt.wantErr && err != nil {
				t.Errorf("err: '%v'", err)
			}
		})
	}
}
