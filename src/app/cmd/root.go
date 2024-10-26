package cmd

import (
	"github.com/pablor21/goms/app/config"
	"github.com/pablor21/goms/app/server"
	"github.com/pablor21/goms/pkg/database"
	"github.com/spf13/cobra"
)

func Run() error {
	rootCmd := &cobra.Command{
		Use:   config.GetConfig().App.Name,
		Short: config.GetConfig().App.Description,
		Long:  config.GetConfig().App.Description,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
		},
	}

	database.AddDatabaseCmd(rootCmd)
	server.AddServerCmd(rootCmd)

	return rootCmd.Execute()
}
