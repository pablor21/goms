package generators

import (
	"github.com/pablor21/goms/app/generators/query"
	"github.com/spf13/cobra"
)

func AddGeneratorCmd(root *cobra.Command) {
	generatorCmd := &cobra.Command{
		Use: "generate",
	}

	// Add subcommands
	query.AddReposCmd(generatorCmd)

	root.AddCommand(generatorCmd)

}
