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

	selector := pterm.DefaultInteractiveSelect.
		WithMaxHeight(MaxSelectHeight).
		WithOptions(keys)

	// check if there are defaults for the selection and set them accordingly
	switch c.Config.Group.Selection.Default {
	// if the default is empty, return without setting default option
	case "":
		return selector, nil
	// if the default select is "-", set the current group as the default selection option
	case "-":
		return selector.WithDefaultOption(c.State.Group.Active), nil
	// search for the given default selection group
	default:
		// get the default selection group
		_, ok := lo.Find(c.Config.Group.Items, func(item config.GroupItem) bool {
			return item.Name == c.Config.Group.Selection.Default
		})
		if !ok {
			return nil, fmt.Errorf("could not find default selection group: '%s'", c.State.Group.Active)
		}

		return selector.WithDefaultOption(c.Config.Group.Selection.Default), nil
	}
}
