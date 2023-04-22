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
		name string
		args args
		want *pterm.InteractiveSelectPrinter
	}{
		{
			name: "should return a printer, that sorts the given group as they are given",
			args: args{
				Config: &config.Config{
					Group: config.Group{
						Items: []config.GroupItem{
							{
								Name: "b",
							},
							{
								Name: "c",
							},
							{
								Name: "a",
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
					"b",
					"c",
					"a",
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
		},
		{
			name: "should return a printer, that sorts the given group ascending",
			args: args{
				Config: &config.Config{
					Group: config.Group{
						Selection: config.Selection{
							Sort: "asc",
						},
						Items: []config.GroupItem{
							{
								Name: "c",
							},
							{
								Name: "b",
							},
							{
								Name: "a",
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
					"a",
					"b",
					"c",
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
		},
		{
			name: "should return a printer, that sorts the given group descending",
			args: args{
				Config: &config.Config{
					Group: config.Group{
						Selection: config.Selection{
							Sort: "desc",
						},
						Items: []config.GroupItem{
							{
								Name: "a",
							},
							{
								Name: "b",
							},
							{
								Name: "c",
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
					"c",
					"b",
					"a",
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
		},
		{
			name: "should return a printer, that sorts the given group descending and set a default",
			args: args{
				Config: &config.Config{
					Group: config.Group{
						Selection: config.Selection{
							Sort:    "desc",
							Default: "c",
						},
						Items: []config.GroupItem{
							{
								Name: "a",
							},
							{
								Name: "b",
							},
							{
								Name: "c",
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
					"c",
					"b",
					"a",
				},
				OptionStyle: &pterm.Style{
					pterm.FgDefault,
					pterm.BgDefault,
				},
				DefaultOption: "c",
				MaxHeight:     MaxSelectHeight,
				Selector:      ">",
				SelectorStyle: &pterm.Style{
					pterm.FgLightMagenta,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := Client{
				Config: tt.args.Config,
				State:  tt.args.State,
			}

			got := client.buildInteractiveSelectPrinter()

			options := cmpopts.IgnoreUnexported(pterm.InteractiveSelectPrinter{})
			if !cmp.Equal(&tt.want, &got, options) {
				diff := cmp.Diff(tt.want, got, options)
				t.Errorf("group.Set() state mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
