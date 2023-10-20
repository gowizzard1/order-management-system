package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nu_device_manager",
	Short: "NU Device Manager",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	err := cmd.Usage()
	if err != nil {
		log.Error(err)
	}
}

// Execute cmd pkg entry point
func Execute() (err error) {
	return rootCmd.Execute()
}
