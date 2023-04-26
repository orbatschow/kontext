package group

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/pterm/pterm"
)

func Test_BuildTablePrinter(t *testing.T) {
	type args struct {
		State  *state.State
		Groups []config.GroupItem
	}
	tests := []struct {
		name    string
		args    args
		want    *pterm.TablePrinter
		wantErr bool
	}{
		{
			name: "should print a group with one source, that is active",
			args: args{
				State: &state.State{
					Group: state.Group{
						Active: "dev",
					},
				},
				Groups: []config.GroupItem{
					{
						Name: "dev",
						Sources: []string{
							"dev-local",
						},
					},
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
					[]string{"Active", "Name", "Source(s)"},
					[]string{"*", "dev", "dev-local"},
				},
				Boxed:          false,
				LeftAlignment:  true,
				RightAlignment: false,
				Writer:         nil,
			},
			wantErr: false,
		},
		{
			name: "should print a group with two sources, that are split by \n",
			args: args{
				State: &state.State{
					Group: state.Group{},
				},
				Groups: []config.GroupItem{
					{
						Name: "dev",
						Sources: []string{
							"dev-local",
							"dev-kind",
						},
					},
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
					[]string{"Active", "Name", "Source(s)"},
					[]string{"", "dev", "dev-local\ndev-kind"},
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

			got := client.BuildTablePrinter(tt.args.Groups...)
			options := cmpopts.IgnoreUnexported(pterm.InteractiveSelectPrinter{})
			if !cmp.Equal(&tt.want, &got, options) {
				diff := cmp.Diff(tt.want, got, options)
				t.Errorf("context.BuildTablePrinter() state mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
