package group

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/pterm/pterm"
)

func Test_buildInteractiveSelectPrinter(t *testing.T) {
	type args struct {
		Config *config.Config
		State  *state.State
	}
	tests := []struct {
		name    string
		args    args
		want    *pterm.InteractiveSelectPrinter
		wantErr bool
	}{
		{
			name: "should return a printer, that sorts the given group as they are given",
			args: args{
				Config: &config.Config{
					Group: config.Group{
						Items: []config.GroupItem{
							{
								Name: "group-b",
							},
							{
								Name: "group-c",
							},
							{
								Name: "group-a",
							},
						},
					},
				},
				State: &state.State{},
			},
			want: &pterm.InteractiveSelectPrinter{
				TextStyle: &pterm.Style{
					pterm.FgLightCyan,
				},
				DefaultText: "Please select an option",
				Options: []string{
					"group-b",
					"group-c",
					"group-a",
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
			name: "should return a printer, that sorts the given group ascending",
			args: args{
				Config: &config.Config{
					Group: config.Group{
						Selection: config.Selection{
							Sort: SortAsc,
						},
						Items: []config.GroupItem{
							{
								Name: "group-c",
							},
							{
								Name: "group-b",
							},
							{
								Name: "group-a",
							},
						},
					},
				},
				State: &state.State{},
			},
			want: &pterm.InteractiveSelectPrinter{
				TextStyle: &pterm.Style{
					pterm.FgLightCyan,
				},
				DefaultText: "Please select an option",
				Options: []string{
					"group-a",
					"group-b",
					"group-c",
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
			name: "should return a printer, that sorts the given group descending",
			args: args{
				Config: &config.Config{
					Group: config.Group{
						Selection: config.Selection{
							Sort: SortDesc,
						},
						Items: []config.GroupItem{
							{
								Name: "group-a",
							},
							{
								Name: "group-b",
							},
							{
								Name: "group-c",
							},
						},
					},
				},
				State: &state.State{},
			},
			want: &pterm.InteractiveSelectPrinter{
				TextStyle: &pterm.Style{
					pterm.FgLightCyan,
				},
				DefaultText: "Please select an option",
				Options: []string{
					"group-c",
					"group-b",
					"group-a",
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
			name: "should return a printer, that sorts the given group descending and set a default",
			args: args{
				Config: &config.Config{
					Group: config.Group{
						Selection: config.Selection{
							Sort:    SortDesc,
							Default: "group-c",
						},
						Items: []config.GroupItem{
							{
								Name: "group-a",
							},
							{
								Name: "group-b",
							},
							{
								Name: "group-c",
							},
						},
					},
				},
				State: &state.State{},
			},
			want: &pterm.InteractiveSelectPrinter{
				TextStyle: &pterm.Style{
					pterm.FgLightCyan,
				},
				DefaultText: "Please select an option",
				Options: []string{
					"group-c",
					"group-b",
					"group-a",
				},
				OptionStyle: &pterm.Style{
					pterm.FgDefault,
					pterm.BgDefault,
				},
				DefaultOption: "group-c",
				MaxHeight:     MaxSelectHeight,
				Selector:      ">",
				SelectorStyle: &pterm.Style{
					pterm.FgLightMagenta,
				},
			},
			wantErr: false,
		},
		{
			name: "should return a printer, that sets the default selection to the active group",
			args: args{
				Config: &config.Config{
					Group: config.Group{
						Selection: config.Selection{
							Default: "-",
						},
						Items: []config.GroupItem{
							{
								Name: "group-a",
							},
							{
								Name: "group-b",
							},
							{
								Name: "group-c",
							},
						},
					},
				},
				State: &state.State{
					Group: state.Group{
						Active: "group-b",
					},
				},
			},
			want: &pterm.InteractiveSelectPrinter{
				TextStyle: &pterm.Style{
					pterm.FgLightCyan,
				},
				DefaultText: "Please select an option",
				Options: []string{
					"group-a",
					"group-b",
					"group-c",
				},
				OptionStyle: &pterm.Style{
					pterm.FgDefault,
					pterm.BgDefault,
				},
				DefaultOption: "group-b",
				MaxHeight:     MaxSelectHeight,
				Selector:      ">",
				SelectorStyle: &pterm.Style{
					pterm.FgLightMagenta,
				},
			},
			wantErr: false,
		},
		{
			name: "should throw an error, as the default selection group does not exist",
			args: args{
				Config: &config.Config{
					Group: config.Group{
						Selection: config.Selection{
							Default: "dev",
						},
						Items: []config.GroupItem{},
					},
				},
				State: &state.State{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := Client{
				Config: tt.args.Config,
				State:  tt.args.State,
			}

			got, err := client.buildInteractiveSelectPrinter()
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error, err: '%v'", err)
			}

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got: '%v'", err)
			}

			options := cmpopts.IgnoreUnexported(pterm.InteractiveSelectPrinter{})
			if !tt.wantErr && !cmp.Equal(&tt.want, &got, options) {
				diff := cmp.Diff(tt.want, got, options)
				t.Errorf("group.buildInteractiveSelectPrinter() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
