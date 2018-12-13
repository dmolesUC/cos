// Copyright © 2018 David Moles <david.moles@ucop.edu>


package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const longDescription = `
Verify the checksum of a file
[TODO: long description]
`

// fixityCmd represents the fixity command
var fixityCmd = &cobra.Command{
	Use:   "fixity",
	Short: "Verify the checksum of a file",
	Long: strings.TrimSpace(longDescription),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fixity called")
	},
}

func init() {
	rootCmd.AddCommand(fixityCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fixityCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fixityCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
