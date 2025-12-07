package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "A simple CLI application",
	Long:  "A simple CLI application built with Cobra",
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}

}
