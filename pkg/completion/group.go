package completion

import (
	"github.com/carapace-sh/carapace"
	"github.com/orbatschow/kontext/pkg/group"
)

func Group() carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		client, err := group.New()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		var values []string
		for _, g := range client.Config.Group.Items {
			values = append(values, g.Name)
		}

		return carapace.ActionValues(values...)
	})
}
