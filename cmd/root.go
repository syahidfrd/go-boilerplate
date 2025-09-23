package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Execute() {
	var command = &cobra.Command{
		Use:   "go-boilerplate",
		Short: "Run service",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.AddCommand(serverCmd())

	if err := command.Execute(); err != nil {
		log.Fatal().Msgf("failed run app: %s", err.Error())
	}
}
