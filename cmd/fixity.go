// Copyright © 2018 David Moles <david.moles@ucop.edu>


package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/url"
	"strings"
)

const longDescription = `
Verify the digest of a file
[TODO: long description]
`

var expectedDigest *[]byte
var digestAlgorithm *string

var fixityCmd = &cobra.Command{
	Use: "fixity <URL>",
	Short: "Verify the digest of an object",
	Long: strings.TrimSpace(longDescription),
	Args: cobra.ExactArgs(1),
	Example: `
      fixity s3://mybucket/myprefix/myobject
      fixity s3://mybucket/myprefix/myobject -e 9f1d6f7b77f74a091c7734ae15f394be3f67c73a8568d623197c01ba4f05e9ff
      fixity s3://mybucket/myprefix/myobject -a md5 f7f647f7d9338fb16ba3a1da8453fee3
    `,
	Run: func(cmd *cobra.Command, args []string) {
		checkFixity(args[0])
	},
}

func checkFixity(objUrlStr string) {
	fmt.Printf("expectedDigest  : %x\n", *expectedDigest)
	fmt.Printf("digestAlgorithm : %s\n", *digestAlgorithm)
	objUrl, err := url.Parse(objUrlStr)
	if err != nil {
		log.Fatal(err) // TODO: figure out why url.Parse never fails
	}
	fmt.Println("Object URL: " + objUrl.String())
}


func init() {
	// TODO: add download option
	// TODO: validate algorithm
	flags := fixityCmd.Flags()
	expectedDigest = flags.BytesHexP("expected", "e", nil, "Expected digest value")
	digestAlgorithm = flags.StringP("algorithm", "a", "sha256", "Digest algorithm (md5, sha256; default is sha256)")

	rootCmd.AddCommand(fixityCmd)
}
