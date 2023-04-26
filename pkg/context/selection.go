package context

import (
	"fmt"
	"sort"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/pterm/pterm"
	"github.com/samber/lo"
)

// start an interactive context selection
func (c *Client) buildInteractiveSelectPrinter() (*pterm.InteractiveSelectPrinter, error) {
	// compute all selection options
	var keys []string
	for k := range c.APIConfig.Contexts {
		keys = append(keys, k)
	}

	// get the active group
	group, ok := lo.Find(c.Config.Group.Items, func(item config.GroupItem) bool {
		return item.Name == c.State.Group.Active
	})
	if !ok {
		return nil, fmt.Errorf("could not find default selection context: '%s'", c.State.Group.Active)
	}

	// sort the selection
	switch group.Context.Selection.Sort {
	case SortAsc:
		sort.Strings(keys)
	case SortDesc:
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	default:
		sort.Strings(keys)
	}

	selector := pterm.DefaultInteractiveSelect.
		WithMaxHeight(MaxSelectHeight).
		WithOptions(keys)

	// check if there are defaults for the selection and set them accordingly
	switch group.Context.Selection.Default {
	// if the default is empty, return without setting default option
	case "":
		return selector, nil
	// if the default select is "-", set the current context as the default option
	case "-":
		return selector.WithDefaultOption(c.State.Context.Active), nil
	// search for the given default selection context
	default:
		// get the default selection context
		_, ok := c.APIConfig.Contexts[group.Context.Selection.Default]
		if !ok {
			return nil, fmt.Errorf("could not find default selection context: '%s'", group.Context.Selection.Default)
		}
		return selector.WithDefaultOption(group.Context.Selection.Default), nil
	}
}
