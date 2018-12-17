// Copyright © 2018 David Moles <david.moles@ucop.edu>


package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/url"
	"strings"
)

const shortDescription = "Verify the digest of an object"

const longDescription = shortDescription + `
[TODO: long description]
`

const example = `
  fixity s3://mybucket/myprefix/myobject
  fixity s3://mybucket/myprefix/myobject -e 9f1d6f7b77f74a091c7734ae15f394be3f67c73a8568d623197c01ba4f05e9ff
  fixity s3://mybucket/myprefix/myobject -a md5 f7f647f7d9338fb16ba3a1da8453fee3
`

var expectedDigest *[]byte
var digestAlgorithm *string

var fixityCmd = &cobra.Command{
	Use: "fixity <OBJECT-URL>",
	Short: shortDescription,
	Long: strings.TrimSpace(longDescription),
	Args: cobra.ExactArgs(1),
	SilenceUsage: true,
	SilenceErrors: true,
	Example: "  " + strings.TrimSpace(example),
	RunE: func(cmd *cobra.Command, args []string) error {
		return checkFixity(args[0])
	},
}

func checkFixity(objUrlStr string) error {
	objUrl, err := url.Parse(objUrlStr)
	if err == nil {
		return checkFixityUrl(objUrl)
	}
	return err
}

func checkFixityUrl(objUrl *url.URL) error {
	fmt.Printf("Object URL: %v\n", objUrl)
	fmt.Printf("expectedDigest  : %x\n", *expectedDigest)
	fmt.Printf("digestAlgorithm : %s\n", *digestAlgorithm)
	return nil
}

func init() {
	// TODO: add download option
	// TODO: validate algorithm
	flags := fixityCmd.Flags()
	expectedDigest = flags.BytesHexP("expected", "e", nil, "Expected digest value")
	digestAlgorithm = flags.StringP("algorithm", "a", "sha256", "Digest algorithm (md5, sha256; default is sha256)")

	rootCmd.AddCommand(fixityCmd)
}
