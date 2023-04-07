package context

import (
	"reflect"
	"testing"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/state"
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
				t.Errorf("expected error")
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want: '%v', got: '%v'", tt.want, got)
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
		name          string
		args          args
		wantState     *state.State
		wantAPIConfig *api.Config
		wantErr       bool
	}{
		{
			name: "should change the api config and state to the given context",
			args: args{
				ContextName: "kind",
				Config:      &config.Config{},
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
			wantState: &state.State{
				Context: state.Context{
					Active: "kind",
					History: []string{
						"kind",
					},
				},
			},
			wantAPIConfig: &api.Config{
				CurrentContext: "kind",
				Contexts: map[string]*api.Context{
					"kind":  {},
					"local": {},
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
				t.Errorf("err: '%v'", err)
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.wantState, client.State) {
				t.Errorf("want: '%v', got: '%v'", tt.wantState, client.State)
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.wantAPIConfig, client.APIConfig) {
				t.Errorf("want: '%v', got: '%v'", tt.wantAPIConfig, client.APIConfig)
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
		wantErr bool
	}{
		{
			name: "should print without an error",
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
			wantErr: false,
		},
		{
			name: "should print without an error, even when there is context",
			args: args{
				State:    &state.State{},
				Contexts: map[string]*api.Context{},
			},
			wantErr: false,
		},
		{
			name: "should print without an error, even when there is no active context",
			args: args{
				State: &state.State{},
				Contexts: map[string]*api.Context{
					"kind":  {},
					"local": {},
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

			err := client.Print(tt.args.Contexts)
			if !tt.wantErr && err != nil {
				t.Errorf("err: '%v'", err)
			}
		})
	}
}
