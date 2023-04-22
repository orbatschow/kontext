package context

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/pterm/pterm"
	"k8s.io/client-go/tools/clientcmd/api"
)

func Test_Get(t *testing.T) {
	type args struct {
		ContextName string
		APIConfig   *api.Config
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*api.Context
		wantErr bool
	}{
		{
			name: "should return the context within a map",
			args: args{
				ContextName: "kind",
				APIConfig: &api.Config{
					Contexts: map[string]*api.Context{
						"kind": {},
					},
				},
			},
			want: map[string]*api.Context{
				"kind": {},
			},
			wantErr: false,
		},
		{
			name:    "should throw an error, due to missing context name",
			args:    args{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := Client{
				APIConfig: tt.args.APIConfig,
			}

			got, err := client.Get(tt.args.ContextName)
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got: '%v'", err)
			}

			if !tt.wantErr && tt.want != nil && !cmp.Equal(tt.want, got) {
				diff := cmp.Diff(&tt.want, &client.APIConfig)
				t.Errorf("client.Get() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_List(t *testing.T) {
	type args struct {
		ContextName string
		APIConfig   *api.Config
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*api.Context
		wantErr bool
	}{
		{
			name: "should return the context within a map",
			args: args{
				ContextName: "kind",
				APIConfig: &api.Config{
					Contexts: map[string]*api.Context{
						"kind":  {},
						"local": {},
					},
				},
			},
			want: map[string]*api.Context{
				"kind":  {},
				"local": {},
			},
			wantErr: false,
		},
		{
			name: "should return nil",
			args: args{
				APIConfig: &api.Config{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := Client{
				APIConfig: tt.args.APIConfig,
			}

			got := client.List()
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want: '%v', got: '%v'", tt.want, got)
			}
		})
	}
}

func Test_Set(t *testing.T) {
	type args struct {
		ContextName string
		Config      *config.Config
		State       *state.State
		APIConfig   *api.Config
	}
	tests := []struct {
		name string
		args args
		want *struct {
			state     *state.State
			apiConfig *api.Config
		}
		wantErr bool
	}{
		{
			name: "should change the api config and state to the given context",
			args: args{
				ContextName: "kind",
				Config: &config.Config{
					State: config.State{
						History: config.History{
							Size: state.DefaultMaximumHistorySize,
						},
					},
				},
				APIConfig: &api.Config{
					CurrentContext: "local",
					Contexts: map[string]*api.Context{
						"kind":  {},
						"local": {},
					},
				},
				State: &state.State{
					Context: state.Context{
						Active: "local",
					},
				},
			},
			want: &struct {
				state     *state.State
				apiConfig *api.Config
			}{
				state: &state.State{
					Context: state.Context{
						Active: "kind",
						History: []state.History{
							"kind",
						},
					},
				},
				apiConfig: &api.Config{
					CurrentContext: "kind",
					Contexts: map[string]*api.Context{
						"kind":  {},
						"local": {},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should change the api config, state and history to the given context",
			args: args{
				ContextName: "kind",
				Config: &config.Config{
					State: config.State{
						History: config.History{
							Size: state.DefaultMaximumHistorySize,
						},
					},
				},
				APIConfig: &api.Config{
					CurrentContext: "local",
					Contexts: map[string]*api.Context{
						"kind":  {},
						"local": {},
					},
				},
				State: &state.State{
					Context: state.Context{
						Active: "local",
						History: []state.History{
							"local",
						},
					},
				},
			},
			want: &struct {
				state     *state.State
				apiConfig *api.Config
			}{
				state: &state.State{
					Context: state.Context{
						Active: "kind",
						History: []state.History{
							"local",
							"kind",
						},
					},
				},
				apiConfig: &api.Config{
					CurrentContext: "kind",
					Contexts: map[string]*api.Context{
						"kind":  {},
						"local": {},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := Client{
				Config:    tt.args.Config,
				State:     tt.args.State,
				APIConfig: tt.args.APIConfig,
			}

			err := client.Set(tt.args.ContextName)
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected err: '%v'", err)
			}

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got: '%v'", err)
			}

			if !tt.wantErr && tt.want != nil && !cmp.Equal(tt.want.apiConfig, client.APIConfig) {
				diff := cmp.Diff(&tt.want, &client.APIConfig)
				t.Errorf("client.Get() apiConfig mismatch (-want +got):\n%s", diff)
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.want.state, client.State) {
				diff := cmp.Diff(&tt.want, &client.APIConfig)
				t.Errorf("client.Get() state mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_Print(t *testing.T) {
	type args struct {
		State    *state.State
		Contexts map[string]*api.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *pterm.TablePrinter
		wantErr bool
	}{
		{
			name: "should print two contexts, one of them being active",
			args: args{
				State: &state.State{
					Context: state.Context{
						Active: "local",
					},
				},
				Contexts: map[string]*api.Context{
					"kind":  {},
					"local": {},
				},
			},
			want: &pterm.TablePrinter{
				Style: &pterm.Style{
					pterm.FgDefault,
				},
				HasHeader: true,
				HeaderStyle: &pterm.Style{
					pterm.FgLightCyan,
				},
				HeaderRowSeparator: "",
				HeaderRowSeparatorStyle: &pterm.Style{
					pterm.FgGray,
				},
				Separator: " | ",
				SeparatorStyle: &pterm.Style{
					pterm.FgGray,
				},
				RowSeparator: "",
				RowSeparatorStyle: &pterm.Style{
					pterm.FgGray,
				},
				Data: pterm.TableData{
					[]string{"Active", "Name", "Cluster", "AuthInfo"},
					[]string{"", "kind", "", ""},
					[]string{"*", "local", "", ""},
				},
				Boxed:          false,
				LeftAlignment:  true,
				RightAlignment: false,
				Writer:         nil,
			},
			wantErr: false,
		},
		{
			name: "should print two contexts, none of them being active",
			args: args{
				State: &state.State{},
				Contexts: map[string]*api.Context{
					"kind":  {},
					"local": {},
				},
			},
			want: &pterm.TablePrinter{
				Style: &pterm.Style{
					pterm.FgDefault,
				},
				HasHeader: true,
				HeaderStyle: &pterm.Style{
					pterm.FgLightCyan,
				},
				HeaderRowSeparator: "",
				HeaderRowSeparatorStyle: &pterm.Style{
					pterm.FgGray,
				},
				Separator: " | ",
				SeparatorStyle: &pterm.Style{
					pterm.FgGray,
				},
				RowSeparator: "",
				RowSeparatorStyle: &pterm.Style{
					pterm.FgGray,
				},
				Data: pterm.TableData{
					[]string{"Active", "Name", "Cluster", "AuthInfo"},
					[]string{"", "kind", "", ""},
					[]string{"", "local", "", ""},
				},
				Boxed:          false,
				LeftAlignment:  true,
				RightAlignment: false,
				Writer:         nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := Client{
				State: tt.args.State,
			}

			got := client.BuildTablePrinter(tt.args.Contexts)
			options := cmpopts.IgnoreUnexported(pterm.InteractiveSelectPrinter{})
			if !cmp.Equal(&tt.want, &got, options) {
				diff := cmp.Diff(tt.want, got, options)
				t.Errorf("context.BuildTablePrinter() state mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
