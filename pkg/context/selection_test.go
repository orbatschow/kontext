package context

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/pterm/pterm"
	"k8s.io/client-go/tools/clientcmd/api"
)

func Test_buildInteractiveSelectPrinter(t *testing.T) {
	type args struct {
		Config    *config.Config
		APIConfig *api.Config
		State     *state.State
	}
	tests := []struct {
		name    string
		args    args
		want    *pterm.InteractiveSelectPrinter
		wantErr bool
	}{
		{
			name: "should return a printer, that sorts the given contexts by default (ascending)",
			args: args{
				APIConfig: &api.Config{
					Contexts: map[string]*api.Context{
						"kind-b": nil,
						"kind-c": nil,
						"kind-a": nil,
					},
				},
				Config: &config.Config{
					Group: config.Group{
						Items: []config.GroupItem{
							{
								Name: "group-a",
							},
						},
					},
				},
				State: &state.State{
					Group: state.Group{
						Active: "group-a",
					},
				},
			},
			want: &pterm.InteractiveSelectPrinter{
				TextStyle: &pterm.Style{
					pterm.FgLightCyan,
				},
				DefaultText: "Please select an option",
				Options: []string{
					"kind-a",
					"kind-b",
					"kind-c",
				},
				OptionStyle: &pterm.Style{
					pterm.FgDefault,
					pterm.BgDefault,
				},
				DefaultOption: "",
				MaxHeight:     MaxSelectHeight,
				Selector:      ">",
				SelectorStyle: &pterm.Style{
					pterm.FgLightMagenta,
				},
			},
			wantErr: false,
		},
		{
			name: "should return a printer, that sorts the given contexts ascending",
			args: args{
				APIConfig: &api.Config{
					Contexts: map[string]*api.Context{
						"kind-b": nil,
						"kind-c": nil,
						"kind-a": nil,
					},
				},
				Config: &config.Config{
					Group: config.Group{
						Items: []config.GroupItem{
							{
								Name: "group-a",
								Context: config.Context{
									Selection: config.Selection{
										Sort: SortAsc,
									},
								},
							},
						},
					},
				},
				State: &state.State{
					Group: state.Group{
						Active: "group-a",
					},
				},
			},
			want: &pterm.InteractiveSelectPrinter{
				TextStyle: &pterm.Style{
					pterm.FgLightCyan,
				},
				DefaultText: "Please select an option",
				Options: []string{
					"kind-a",
					"kind-b",
					"kind-c",
				},
				OptionStyle: &pterm.Style{
					pterm.FgDefault,
					pterm.BgDefault,
				},
				DefaultOption: "",
				MaxHeight:     MaxSelectHeight,
				Selector:      ">",
				SelectorStyle: &pterm.Style{
					pterm.FgLightMagenta,
				},
			},
			wantErr: false,
		},
		{
			name: "should return a printer, that sorts the given contexts ascending",
			args: args{
				APIConfig: &api.Config{
					Contexts: map[string]*api.Context{
						"kind-b": nil,
						"kind-c": nil,
						"kind-a": nil,
					},
				},
				Config: &config.Config{
					Group: config.Group{
						Items: []config.GroupItem{
							{
								Name: "group-a",
								Context: config.Context{
									Selection: config.Selection{
										Sort: SortDesc,
									},
								},
							},
						},
					},
				},
				State: &state.State{
					Group: state.Group{
						Active: "group-a",
					},
				},
			},
			want: &pterm.InteractiveSelectPrinter{
				TextStyle: &pterm.Style{
					pterm.FgLightCyan,
				},
				DefaultText: "Please select an option",
				Options: []string{
					"kind-c",
					"kind-b",
					"kind-a",
				},
				OptionStyle: &pterm.Style{
					pterm.FgDefault,
					pterm.BgDefault,
				},
				DefaultOption: "",
				MaxHeight:     MaxSelectHeight,
				Selector:      ">",
				SelectorStyle: &pterm.Style{
					pterm.FgLightMagenta,
				},
			},
			wantErr: false,
		},
		{
			name: "should return a printer, that sorts the given contexts descending and set a default",
			args: args{
				APIConfig: &api.Config{
					Contexts: map[string]*api.Context{
						"kind-b": nil,
						"kind-c": nil,
						"kind-a": nil,
					},
				},
				Config: &config.Config{
					Group: config.Group{
						Items: []config.GroupItem{
							{
								Name: "group-a",
								Context: config.Context{
									Default: "kind-b",
									Selection: config.Selection{
										Sort: SortDesc,
									},
								},
							},
						},
					},
				},
				State: &state.State{
					Group: state.Group{
						Active: "group-a",
					},
				},
			},
			want: &pterm.InteractiveSelectPrinter{
				TextStyle: &pterm.Style{
					pterm.FgLightCyan,
				},
				DefaultText: "Please select an option",
				Options: []string{
					"kind-c",
					"kind-b",
					"kind-a",
				},
				OptionStyle: &pterm.Style{
					pterm.FgDefault,
					pterm.BgDefault,
				},
				DefaultOption: "",
				MaxHeight:     MaxSelectHeight,
				Selector:      ">",
				SelectorStyle: &pterm.Style{
					pterm.FgLightMagenta,
				},
			},
			wantErr: false,
		},
		{
			name: "should return a printer, that sets the default selection to the current context",
			args: args{
				APIConfig: &api.Config{
					Contexts: map[string]*api.Context{
						"kind-b": nil,
						"kind-c": nil,
						"kind-a": nil,
					},
					CurrentContext: "kind-b",
				},
				Config: &config.Config{
					Group: config.Group{
						Items: []config.GroupItem{
							{
								Name: "group-a",
								Context: config.Context{
									Selection: config.Selection{
										Default: "-",
									},
								},
							},
						},
					},
				},
				State: &state.State{
					Group: state.Group{
						Active: "group-a",
					},
					Context: state.Context{
						Active: "kind-b",
					},
				},
			},
			want: &pterm.InteractiveSelectPrinter{
				TextStyle: &pterm.Style{
					pterm.FgLightCyan,
				},
				DefaultText: "Please select an option",
				Options: []string{
					"kind-a",
					"kind-b",
					"kind-c",
				},
				OptionStyle: &pterm.Style{
					pterm.FgDefault,
					pterm.BgDefault,
				},
				DefaultOption: "kind-b",
				MaxHeight:     MaxSelectHeight,
				Selector:      ">",
				SelectorStyle: &pterm.Style{
					pterm.FgLightMagenta,
				},
			},
			wantErr: false,
		},

		{
			name: "should return an error, as the default selection context does not exist",
			args: args{
				APIConfig: &api.Config{},
				Config: &config.Config{
					Group: config.Group{
						Items: []config.GroupItem{
							{
								Name: "group-a",
								Context: config.Context{
									Selection: config.Selection{
										Default: "kind-b",
										Sort:    SortDesc,
									},
								},
							},
						},
					},
				},
				State: &state.State{
					Group: state.Group{
						Active: "group-a",
					},
				},
			},
			want:    nil,
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

			got, err := client.buildInteractiveSelectPrinter()
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got: '%v'", err)
			}

			options := cmpopts.IgnoreUnexported(pterm.InteractiveSelectPrinter{})
			if !cmp.Equal(&tt.want, &got, options) {
				diff := cmp.Diff(tt.want, got, options)
				t.Errorf("group.Set() state mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
