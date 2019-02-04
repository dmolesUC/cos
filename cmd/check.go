package cmd

import (
	"fmt"

	"github.com/dmolesUC3/cos/internal/objects"

	"github.com/dmolesUC3/cos/pkg"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/internal/streaming"

	"github.com/spf13/cobra"
)

// ------------------------------------------------------------
// Constants: Help Text

const (
	usageCheck = "check <OBJECT-URL>"

	shortDescCheck = "check: verify the digest of an object"

	longDescCheck = shortDescCheck + `

		Verifies the digest of an object in cloud object storage, using SHA-256 (by
		default) or MD5 (optionally). 
	`

	exampleCheck = ` 
		cos check s3://www.dmoles.net/images/fa/archive.svg --endpoint https://s3.us-west-2.amazonaws.com/
		cos check s3://www.dmoles.net/images/fa/archive.svg -e https://s3.us-west-2.amazonaws.com/ -x c99ad299fa53d5d9688909164cf25b386b33bea8d4247310d80f615be29978f5
		cos check s3://mrt-test/inusitatum.png -e http://127.0.0.1:9000/ -a md5 -x cadf871cd4135212419f488f42c62482
	    SWIFT_API_USER=<user> SWIFT_API_KEY=<key> cos check 'swift://distrib.stage.9001.__c5e/ark:/99999/fk4kw5kc1z|1|producer/6GBZeroFile.txt' -e http://cloud.sdsc.edu/auth/v1.0
    `
)

// ------------------------------------------------------------
// checkFlags type

type checkFlags struct {
	CosFlags

	Expected  []byte
	Algorithm string
}

func (f checkFlags) Pretty() string {
	format := `
		verbose: %v
		expected: %x
		algorithm: '%v'
		endpoint: '%v'
		region: '%v'`
	format = logging.Untabify(format, "  ")
	return fmt.Sprintf(format, f.Verbose, f.Expected, f.Algorithm, f.Endpoint, f.Region)
}

func (f checkFlags) String() string {
	return fmt.Sprintf(
		"checkFlags{ verbose: %v, expected: %x, algorithm: '%v', endpoint: '%v', region: '%v'}",
		f.Verbose, f.Expected, f.Algorithm, f.Endpoint, f.Region,
	)
}

// ------------------------------------------------------------
// Functions

func check(objURLStr string, f checkFlags) error {
	logger := logging.DefaultLoggerWithLevel(f.LogLevel())
	logger.Tracef("flags: %v\n", f)
	logger.Tracef("object URL: %v\n", objURLStr)

	objURL, err := streaming.ValidAbsURL(objURLStr)
	if err != nil {
		return err
	}

	endpointURL, err := streaming.ValidAbsURL(f.Endpoint)
	if err != nil {
		return err
	}

	obj, err := objects.NewObject(objURL, endpointURL, f.Region)
	if err != nil {
		return err
	}
	logger.Tracef("object: %v\n", obj)

	var check = pkg.Check{
		Object:    obj,
		Expected:  f.Expected,
		Algorithm: f.Algorithm,
	}
	digest, err := check.VerifyDigest()
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
		Use:           usageCheck,
		Short:         shortDescCheck,
		Long:          logging.Untabify(longDescCheck, ""),
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       logging.Untabify(exampleCheck, "  "),
		RunE: func(cmd *cobra.Command, args []string) error {
			return check(args[0], flags)
		},
	}
	cmdFlags := cmd.Flags()
	flags.AddTo(cmdFlags)

	cmdFlags.StringVarP(&flags.Algorithm, "algorithm", "a", "sha256", "digest algorithm (md5 or sha256)")
	cmdFlags.BytesHexVarP(&flags.Expected, "expected", "x", nil, "expected digest value (exit with error if not matched)")

	rootCmd.AddCommand(cmd)
}
