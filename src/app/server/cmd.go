package server

import (
	"github.com/pablor21/goms/pkg/logger"
	"github.com/spf13/cobra"
)

func AddServerCmd(root *cobra.Command) {
	serverCmd := &cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, args []string) {
			s := NewServer()
			err := s.Start()
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to start server")
			}
		},
	}
	root.AddCommand(serverCmd)
}
