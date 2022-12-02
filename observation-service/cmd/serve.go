package cmd

import (
	"github.com/spf13/cobra"

	"github.com/caraml-dev/observation-service/observation-service/log"
	"github.com/caraml-dev/observation-service/observation-service/server"
)

var cfgFile []string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts Observation Service CMux server with configured API",
	Run: func(cmd *cobra.Command, args []string) {

		server, err := server.NewServer(cfgFile)
		if err != nil {
			log.Glob().Panic(err)
		}
		server.Start()
	},
}

func init() {
	serveCmd.Flags().StringArrayVar(&cfgFile, "config", []string{},
		`Path to one or more configuration files. The flag can be set multiple times
	and the later values will take precedence.`)
	RootCmd.AddCommand(serveCmd)
}
