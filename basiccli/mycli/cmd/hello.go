package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var name string

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Prints Hello, World!",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(" james bonds & superm an batn wonderwomen !")
	},
}

func init() {

	rootCmd.AddCommand(helloCmd)
	helloCmd.Flags().StringVarP(&name, "name", "n", "World", "Name to greet")
}
