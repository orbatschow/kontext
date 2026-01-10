package completion

import (
	"github.com/carapace-sh/carapace"
	"github.com/orbatschow/kontext/pkg/context"
)

func Context() carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		client, err := context.New()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		contexts := client.List()
		var values []string
		for name := range contexts {
			values = append(values, name)
		}

		return carapace.ActionValues(values...)
	})
}
