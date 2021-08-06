package cmd

import (
	"github.com/onspaceship/booster/pkg/controller"

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Start the controller",
	Run: func(cmd *cobra.Command, args []string) {
		exit := make(chan bool)

		go controller.StartController(exit)

		<-exit
		log.Info("Done")
	},
}

func init() {
	rootCmd.AddCommand(controllerCmd)
}
