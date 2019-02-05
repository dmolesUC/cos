package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var tag string
var commitHash string
var timestamp string

func init() {
	cmd := &cobra.Command{
		Use: "version",
		Short: "print cos version",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("cos " + versionStr())
		},
	}
	rootCmd.AddCommand(cmd)
}

func versionStr() string {
	// if these are blank, we were probably built with plain 'go build' or
	// 'go install', bypassing the ldflags in the magefile
	if tag == "" {
		return "(devel)"
	}
	versionStr := fmt.Sprintf("%v (%v, %v)", tag, commitHash, timestamp)
	return versionStr
}
