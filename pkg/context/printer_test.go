package context

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/pterm/pterm"
	"k8s.io/client-go/tools/clientcmd/api"
)

func Test_BuildTablePrinter(t *testing.T) {
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
