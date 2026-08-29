package prompt

import "github.com/charmbracelet/huh"

func AsOptions(values []string) []huh.Option[string] {
	options := make([]huh.Option[string], len(values))
	for i, v := range values {
		options[i] = huh.NewOption(v, v)
	}
	return options
}
