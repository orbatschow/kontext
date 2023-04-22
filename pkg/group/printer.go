package group

import (
	"strings"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/pterm/pterm"
)

func (c *Client) BuildTablePrinter(groups ...config.GroupItem) *pterm.TablePrinter {
	table := pterm.TableData{
		{"Active", "Name", "Source(s)"},
	}

	for _, value := range groups {
		active := ""
		if value.Name == c.State.Group.Active {
			active = "*"
		}
		table = append(table, []string{
			active, value.Name, strings.Join(value.Sources, "\n"),
		})
	}

	return pterm.DefaultTable.WithHasHeader().WithData(table)
}
