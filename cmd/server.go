package cmd

import (
	"github.com/spf13/cobra"
	"github.com/syahidfrd/go-boilerplate/internal/server"
)

func serverCmd() *cobra.Command {
	var port int
	var command = &cobra.Command{
		Use:   "server",
		Short: "Run server",
		Run: func(cmd *cobra.Command, args []string) {
			srv := server.NewServer()
			srv.Run(port)
		},
	}

	command.Flags().IntVar(&port, "port", 8080, "Listen on given port")
	return command
}
