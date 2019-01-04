package cmd

import (
	"regexp"

	. "github.com/dmolesUC3/coscheck/cos"
	. "github.com/dmolesUC3/coscheck/util"

	"github.com/spf13/cobra"
)

// ------------------------------------------------------------
// Constants: Help Text

const (
	usage = "fixity <OBJECT-URL>"

	shortDescription = "fixity: verify the digest of an object"

	longDescription = shortDescription + `

		Verifies the digest of an object in cloud object storage, using SHA-256 (by
		default) or MD5 (optionally). The object location can be specified either
		as a complete HTTP(S) URL, https://<endpoint>/<bucket>/<key>, or using
		separate URLs for the endpoint (HTTP(S)) and bucket/key (s3://).

		Uses AWS credentials from ~/.aws/config, or other config file specified
		with the AWS_CONFIG_FILE environment variable. Credentials can also be
		specified with the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment
		variables.
	`

	// TODO: add Minio example(s)
	example = ` 
		coscheck fixity https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg
		coscheck fixity https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg -x c99ad299fa53d5d9688909164cf25b386b33bea8d4247310d80f615be29978f5
		coscheck fixity https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg -a md5 -x eac8a75e3b3023e98003f1c24137ebbd
		coscheck fixity s3://www.dmoles.net/images/fa/archive.svg -e https://s3.us-west-2.amazonaws.com/ -a md5 -x eac8a75e3b3023e98003f1c24137ebbd
	`
)

// ------------------------------------------------------------
// fixityFlags type

type fixityFlags struct {
	Verbose   bool
	Expected  []byte
	Algorithm string
	Endpoint  string
	Region    string
}

func (f fixityFlags) logTo(logger Logger) {
	logger.Detailf("verbose   : %v\n", f.Verbose)
	logger.Detailf("algorithm : %v\n", f.Algorithm)
	logger.Detailf("expected  : %x\n", f.Expected)
	logger.Detailf("endpoint  : %v\n", f.Endpoint)
	logger.Detailf("region    : %v\n", f.Region)
}

// ------------------------------------------------------------
// Functions

func formatHelp(text string, indent string) string {
	return regexp.MustCompile(`(?m)^[\t ]+`).ReplaceAllString(text, indent)
}

func runWith(objUrlStr string, f fixityFlags) error {
	var logger = NewLogger(f.Verbose)
	f.logTo(logger)
	logger.Detailf("object URL: %v\n", objUrlStr)

	// TODO: look up default endpoint in S3 config / environment variables?
	objLoc, err := NewObjectLocationFromStrings(&objUrlStr, &f.Endpoint)
	if err != nil {
		return err
	}
	logger.Detailf("ObjectLocation: %v\n", objLoc)

	var fixity = Fixity{
		Logger:    logger,
		ObjLoc:    *objLoc,
		Expected:  f.Expected,
		Algorithm: f.Algorithm,
		Region:    f.Region,
	}

	digest, err := fixity.GetDigest()
	if err != nil {
		return err
	}
	logger.Infof("digest matched: expected: %x, actual: %x\n", f.Expected, digest)
	return nil
}

// ------------------------------------------------------------
// Command initialization

func init() {
	fixity := fixityFlags{}

	cmd := &cobra.Command{
		Use:           usage,
		Short:         shortDescription,
		Long:          formatHelp(longDescription, ""),
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       formatHelp(example, "  "),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWith(args[0], fixity)
		},
	}
	cmd.Flags().BoolVarP(&fixity.Verbose, "verbose", "v", false, "Verbose output")
	cmd.Flags().BytesHexVarP(&fixity.Expected, "expected", "x", nil, "Expected digest value (exit with error if not matched)")
	cmd.Flags().StringVarP(&fixity.Algorithm, "algorithm", "a", "sha256", "Algorithm: md5 or sha256")
	cmd.Flags().StringVarP(&fixity.Endpoint, "endpoint", "e", "", "S3 endpoint: HTTP(S) URL")
	cmd.Flags().StringVarP(&fixity.Endpoint, "region", "r", "", "S3 region (if not in endpoint URL; default \""+DefaultAwsRegion+"\")")

	rootCmd.AddCommand(cmd)
}
