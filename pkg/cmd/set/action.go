package set

import "github.com/rsteube/carapace"

func buildSetContextAction() (carapace.Action, error) {
	var values []string
	values = append(values, "hello")
	values = append(values, "world")
	values = append(values, "foo")
	values = append(values, "bar")
	values = append(values, "test")
	values = append(values, "beskar")
	values = append(values, "iridium")
	values = append(values, "vader")
	return carapace.ActionValues(values...), nil
}
