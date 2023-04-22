package group

import (
	"sort"

	"github.com/pterm/pterm"
)

// start an interactive context selection
func (c *Client) buildInteractiveSelectPrinter() *pterm.InteractiveSelectPrinter {
	// compute all selection options
	var keys []string
	for _, value := range c.Config.Group.Items {
		keys = append(keys, value.Name)
	}

	// sort the selection
	switch c.Config.Group.Selection.Sort {
	case SortAsc:
		sort.Strings(keys)
	case SortDesc:
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	}

	// set the default selection option
	if len(c.Config.Group.Selection.Default) > 0 {
		return pterm.DefaultInteractiveSelect.
			WithMaxHeight(MaxSelectHeight).
			WithOptions(keys).
			WithDefaultOption(c.Config.Group.Selection.Default)
	}
	return pterm.DefaultInteractiveSelect.
		WithMaxHeight(MaxSelectHeight).
		WithOptions(keys)
}
