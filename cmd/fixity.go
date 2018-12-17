// Copyright © 2018 David Moles <david.moles@ucop.edu>

package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
)

const shortDescription = "Verify the digest of an object"

const longDescription = shortDescription + `
[TODO: long description]
`

const example = `
  fixity http://s3-myregion.amazonaws.com/mybucket/myprefix/myobject
  fixity http://s3-myregion.amazonaws.com/mybucket/myprefix/myobject -e 9f1d6f7b77f74a091c7734ae15f394be3f67c73a8568d623197c01ba4f05e9ff
  fixity http://s3-myregion.amazonaws.com/mybucket/myprefix/myobject -a md5 f7f647f7d9338fb16ba3a1da8453fee3
`

var expectedDigest *[]byte
var digestAlgorithm *string

var fixityCmd = &cobra.Command{
	Use:           "fixity <OBJECT-URL>",
	Short:         shortDescription,
	Long:          strings.TrimSpace(longDescription),
	Args:          cobra.ExactArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,
	Example:       "  " + strings.TrimSpace(example),
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

	s3Config := &aws.Config{
		/*
		   coscheck fixity https://s3-us-west-2.amazonaws.com/www.dmoles.net/favicon.png \
		     --expected 8d3496edee4a8a1bb3ba54e8f762aa9a6be9ce69f84a0e8fec0e4809deaaf804
		 */
		Endpoint: aws.String(objUrl.String()), // TODO: decompose into endpoint + region + bucket + key
	}

	sess, err := session.NewSession(s3Config)
	if err != nil {
		return err
	}

	svc := s3.New(sess)
	result, err := svc.ListBuckets(nil)
	if err != nil {
		return err
	}

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n", aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

	return nil
}

func init() {
	// TODO: region, endpoint, auth?
	// TODO: add download option
	// TODO: validate algorithm
	flags := fixityCmd.Flags()
	expectedDigest = flags.BytesHexP("expected", "e", nil, "Expected digest value")
	digestAlgorithm = flags.StringP("algorithm", "a", "sha256", "Digest algorithm (md5, sha256; default is sha256)")

	rootCmd.AddCommand(fixityCmd)
}
