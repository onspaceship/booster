package cmd

import (
	"github.com/onspaceship/booster/pkg/controller"
	"github.com/onspaceship/booster/pkg/socket"

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Start both the server connection and controller",
	Run: func(cmd *cobra.Command, args []string) {
		exit := make(chan bool)

		go socket.StartListener(exit)
		go controller.StartController(exit)

		<-exit
		log.Info("Done")
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
}
