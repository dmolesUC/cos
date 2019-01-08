package cmd

import (
	"fmt"
	"regexp"

	. "github.com/dmolesUC3/cos/pkg"
	. "github.com/dmolesUC3/cos/internal"

	"github.com/spf13/cobra"
)

// ------------------------------------------------------------
// Constants: Help Text

const (
	usage = "check <OBJECT-URL>"

	shortDescription = "check: verify the digest of an object"

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
		cos check https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg
		cos check https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg -x c99ad299fa53d5d9688909164cf25b386b33bea8d4247310d80f615be29978f5
		cos check https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg -a md5 -x eac8a75e3b3023e98003f1c24137ebbd
		cos check s3://www.dmoles.net/images/fa/archive.svg -e https://s3.us-west-2.amazonaws.com/ -a md5 -x eac8a75e3b3023e98003f1c24137ebbd
		cos check s3://mrt-test/inusitatum.png --endpoint http://127.0.0.1:9000/ --algorithm md5 --expected cadf871cd4135212419f488f42c62482`
)

// ------------------------------------------------------------
// checkFlags type

type checkFlags struct {
	Verbose   bool
	Expected  []byte
	Algorithm string
	Endpoint  string
	Region    string
}

func (f checkFlags) logTo(logger Logger) {
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

func runWith(objURLStr string, f checkFlags) error {
	var logger = NewLogger(f.Verbose)
	f.logTo(logger)
	logger.Detailf("object URL: %v\n", objURLStr)

	// TODO: look up default endpoint in S3 config / environment variables?
	objLoc, err := NewObjectLocationFromStrings(&objURLStr, &f.Endpoint)
	if err != nil {
		return err
	}
	logger.Detailf("ObjectLocation: %v\n", objLoc)

	var check = Check{
		Logger:    logger,
		ObjLoc:    *objLoc,
		Expected:  f.Expected,
		Algorithm: f.Algorithm,
		Region:    f.Region,
	}

	digest, err := check.GetDigest()
	if err != nil {
		return err
	}
	fmt.Printf("%x\n", digest)
	return nil
}

// ------------------------------------------------------------
// Command initialization

func init() {
	check := checkFlags{}

	cmd := &cobra.Command{
		Use:           usage,
		Short:         shortDescription,
		Long:          formatHelp(longDescription, ""),
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       formatHelp(example, "  "),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWith(args[0], check)
		},
	}
	cmd.Flags().BoolVarP(&check.Verbose, "verbose", "v", false, "Verbose output")
	cmd.Flags().BytesHexVarP(&check.Expected, "expected", "x", nil, "Expected digest value (exit with error if not matched)")
	cmd.Flags().StringVarP(&check.Algorithm, "algorithm", "a", "sha256", "Algorithm: md5 or sha256")
	cmd.Flags().StringVarP(&check.Endpoint, "endpoint", "e", "", "S3 endpoint: HTTP(S) URL")
	cmd.Flags().StringVarP(&check.Endpoint, "region", "r", "", "S3 region (if not in endpoint URL; default \""+DefaultAwsRegion+"\")")

	rootCmd.AddCommand(cmd)
}
