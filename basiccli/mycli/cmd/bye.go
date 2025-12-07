package cmd

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

// byeCmd represents the bye command
var byeCmd = &cobra.Command{
	Use:   "bye",
	Short: "Prints a goodbye message",
	Long:  `The bye command prints a friendly goodbye message to the user.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Goodbye! Have a great day!")
	},
}

func init(){
	
	rootCmd.AddCommand(byeCmd)
}
