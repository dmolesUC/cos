package cmd

import (
	"fmt"
	"regexp"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/internal/objects"
	"github.com/dmolesUC3/cos/internal/protocols"
	"github.com/dmolesUC3/cos/pkg"

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

        Note that for OpenStack Swift, the API username and key must be specified
        with the SWIFT_API_USER and SWIFT_API_KEY environment variables, and the
        endpoint URL is mandatory.
	`

	// TODO: add Minio example(s)
	example = ` 
		cos check https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg
		cos check https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg -x c99ad299fa53d5d9688909164cf25b386b33bea8d4247310d80f615be29978f5
		cos check s3://www.dmoles.net/images/fa/archive.svg -e https://s3.us-west-2.amazonaws.com/ -a md5 -x eac8a75e3b3023e98003f1c24137ebbd
		cos check s3://mrt-test/inusitatum.png --endpoint http://127.0.0.1:9000/ --algorithm md5 --expected cadf871cd4135212419f488f42c62482
	    SWIFT_API_USER=<user> SWIFT_API_KEY=<key> cos check 'swift://distrib.stage.9001.__c5e/ark:/99999/fk4kw5kc1z|1|producer/6GBZeroFile.txt' -e http://cloud.sdsc.edu/auth/v1.0`
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

func (f checkFlags) String() string {
	return fmt.Sprintf(
		"checkFlags{ verbose: %v, expected: %x, algorithm: '%v', endpoint: '%v', region: '%v'}",
		f.Verbose, f.Expected, f.Algorithm, f.Endpoint, f.Region,
	)
}

// ------------------------------------------------------------
// Functions

func formatHelp(text string, indent string) string {
	return regexp.MustCompile(`(?m)^[\t ]+`).ReplaceAllString(text, indent)
}

func runWith(objURLStr string, flags checkFlags) error {
	var logger = logging.NewLogger(flags.Verbose)
	logger.Detailf("flags: %v\n", flags)
	logger.Detailf("object URL: %v\n", objURLStr)

	obj, err := objects.NewObjectBuilder().
		WithObjectURLStr(objURLStr).
		WithEndpointStr(flags.Endpoint).
		WithRegion(flags.Region).
		Build(logger)
	if err != nil {
		return err
	}
	logger.Detailf("object: %v\n", obj)

	var check = pkg.Check{
		Object:    obj,
		Expected:  flags.Expected,
		Algorithm: flags.Algorithm,
	}
	digest, err := check.CalcDigest()
	if err != nil {
		return err
	}
	fmt.Printf("%x\n", digest)
	return nil
}

// ------------------------------------------------------------
// Command initialization

func init() {
	flags := checkFlags{}

	cmd := &cobra.Command{
		Use:           usage,
		Short:         shortDescription,
		Long:          formatHelp(longDescription, ""),
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       formatHelp(example, "  "),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWith(args[0], flags)
		},
	}
	cmd.Flags().BoolVarP(&flags.Verbose, "verbose", "v", false, "Verbose output")
	cmd.Flags().BytesHexVarP(&flags.Expected, "expected", "x", nil, "Expected digest value (exit with error if not matched)")
	cmd.Flags().StringVarP(&flags.Algorithm, "algorithm", "a", "sha256", "Algorithm: md5 or sha256")
	cmd.Flags().StringVarP(&flags.Endpoint, "endpoint", "e", "", "S3 endpoint: HTTP(S) URL")
	cmd.Flags().StringVarP(&flags.Endpoint, "region", "r", "", "S3 region (if not in endpoint URL; default \""+protocols.DefaultAwsRegion+"\")")

	rootCmd.AddCommand(cmd)
}
