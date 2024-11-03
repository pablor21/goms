package query

import (
	"github.com/spf13/cobra"
)

func AddReposCmd(root *cobra.Command) {
	generatorCmd := &cobra.Command{
		Use: "repos",
		Run: func(cmd *cobra.Command, args []string) {
			generateRepositories()
		},
	}
	root.AddCommand(generatorCmd)
}
