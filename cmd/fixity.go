// Copyright © 2018 David Moles <david.moles@ucop.edu>


package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const longDescription = `
Verify the digest of a file
[TODO: long description]
`

var expected *[]byte

var fixityCmd = &cobra.Command{
	Use: "fixity URL",
	Short: "Verify the checksum of an object",
	Long: strings.TrimSpace(longDescription),
	Args: cobra.ExactArgs(1),
	Example: "fixity s3://mybucket/myprefix/myobject --expected ec57cd0008c934e61b30635efa6964cc3de8574b669175028069c459eeb01510",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("expected: %x\n", *expected);
		fmt.Println("URL: " + strings.Join(args, " "))
	},
}


func init() {
	rootCmd.AddCommand(fixityCmd)
	// TODO: add algorithm, & validate it (default SHA256, support SHA256 or MD5, but case-insensitive)
	expected = fixityCmd.Flags().BytesHexP("expected", "e", nil, "Expected digest value")
}
