package group

import (
	"fmt"
	"sort"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/pterm/pterm"
	"github.com/samber/lo"
)

// start an interactive group selection
func (c *Client) buildInteractiveSelectPrinter() (*pterm.InteractiveSelectPrinter, error) {
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
		// get the default selection group
		_, ok := lo.Find(c.Config.Group.Items, func(item config.GroupItem) bool {
			return item.Name == c.Config.Group.Selection.Default
		})
		if !ok {
			return nil, fmt.Errorf("could not find default selection group: '%s'", c.State.Group.Active)
		}

		return pterm.DefaultInteractiveSelect.
			WithMaxHeight(MaxSelectHeight).
			WithOptions(keys).
			WithDefaultOption(c.Config.Group.Selection.Default), nil
	}
	return pterm.DefaultInteractiveSelect.
		WithMaxHeight(MaxSelectHeight).
		WithOptions(keys), nil
}
