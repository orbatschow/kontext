package context

import (
	"sort"

	"github.com/pterm/pterm"
	"k8s.io/client-go/tools/clientcmd/api"
)

func (c *Client) BuildTablePrinter(contexts map[string]*api.Context) *pterm.TablePrinter {
	table := pterm.TableData{
		{"Active", "Name", "Cluster", "AuthInfo"},
	}

	// sort table data ascending
	var keys []string
	for key := range contexts {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		active := ""
		if key == c.State.Context.Active {
			active = "*"
		}
		table = append(table, []string{
			active, key, contexts[key].Cluster, contexts[key].AuthInfo,
		})
	}
	return pterm.DefaultTable.WithHasHeader().WithData(table)
}
