package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "mydemoapp",
	Short: "我的演示程序",
	Long:  "为了演示而做的程序",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello world")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
