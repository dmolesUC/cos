package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/cobra"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/internal/objects"
	"github.com/dmolesUC3/cos/internal/protocols"
	"github.com/dmolesUC3/cos/pkg"
)

// ------------------------------------------------------------
// Constants: Help Text

const (
	usageCrvd = "crvd <BUCKET-URL>"

	shortDescCrvd = "crvd: create, retrieve, verify, and delete an object"

	longDescCrvd = shortDescCrvd + `
        Creates, retrieves, verifies, and deletes an object in a cloud storage bucket.
        The object consists of a stream of random bytes of the specified size.
    `

	exampleCrvd = `
        cos crvd s3://www.dmoles.net/ --endpoint https://s3.us-west-2.amazonaws.com/
        cos crvd cos check swift://distrib.stage.9001.__c5e/ -e http://cloud.sdsc.edu/auth/v1.0
    `
)

// ------------------------------------------------------------
// crvdFlags type

type crvdFlags struct {
	Verbose  bool
	Region   string
	Endpoint string
	Key      string
	Size     int64
	Seed     int64
	Zero     bool
	Keep     bool
}

func (f crvdFlags) Pretty() string {
	format := `
		verbose:   %v
		region:   '%v'
		endpoint: '%v'
        key:      '%v'
		size:      %d
        seed:      %d
        zero:      %v`
	format = logging.Untabify(format, "  ")
	return fmt.Sprintf(format, f.Verbose, f.Region, f.Endpoint, f.Key, f.Size, f.Seed, f.Zero)
}

func crvd(bucketStr string, flags crvdFlags) (err error) {
	if flags.Key == "" {
		flags.Key = fmt.Sprintf("cos-crvd-%d.bin", time.Now().Unix())
	}
	if flags.Zero {
		flags.Size = 0
	}

	var logger = logging.NewLogger(flags.Verbose)
	logger.Detailf("flags: %v\n", flags)
	logger.Detailf("bucket URL: %v\n", bucketStr)

	if flags.Endpoint == "" {
		return fmt.Errorf("endpoint URL must be specified")
	}

	bucketUrl, err := objects.ValidAbsURL(bucketStr)
	if err != nil {
		return err
	}

	obj, err := objects.NewObjectBuilder().
		WithEndpointStr(flags.Endpoint).
		WithRegion(flags.Region).
		WithProtocolUri(bucketUrl, logger).
		WithKey(flags.Key).
		Build(logger)
	if err != nil {
		return err
	}

	var crvd = pkg.Crvd{Object: obj}
	random := rand.New(rand.NewSource(flags.Seed))

	var digest []byte
	if flags.Keep {
		digest, err = crvd.CreateRetrieveValidate(random, flags.Size)
		if err == nil {
			logger.Infof("created %v (%d bytes, SHA-256 digest %x)\n", objects.ProtocolUriStr(obj), flags.Size, digest)
		}
	} else {
		digest, err = crvd.CreateRetrieveValidateDelete(random, flags.Size)
		if err == nil {
			logger.Infof("verified and deleted %v (%d bytes, SHA-256 digest %x)\n", objects.ProtocolUriStr(obj), flags.Size, digest)
		}
	}
	return err
}

func init() {
	flags := crvdFlags{}
	cmd := &cobra.Command{
		Use:           usageCrvd,
		Short:         shortDescCrvd,
		Long:          logging.Untabify(longDescCrvd, ""),
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       logging.Untabify(exampleCrvd, "  "),
		RunE: func(cmd *cobra.Command, args []string) error {
			return crvd(args[0], flags)
		},
	}
	cmd.Flags().BoolVarP(&flags.Keep, "keep", "", false, "keep object after verification (defaults to false)")
	cmd.Flags().StringVarP(&flags.Endpoint, "endpoint", "e", "", "endpoint: HTTP(S) URL (required)")
	cmd.Flags().StringVarP(&flags.Key, "key", "k", "", "key to create (defaults to cos-crvd-TIMESTAMP.bin)")
	cmd.Flags().Int64VarP(&flags.Size, "random-seed", "", 0, "seed for random-number generator")
	cmd.Flags().StringVarP(&flags.Region, "region", "r", "", "S3 region (if not in endpoint URL; default \""+protocols.DefaultAwsRegion+"\")")
	// TODO: support human-readable sizes (5K, 6MB, 3.2GB, 1TiB, etc.)
	cmd.Flags().Int64VarP(&flags.Size, "size", "s", 1024, "size in bytes of object to create, if --zero not set")
	cmd.Flags().BoolVarP(&flags.Verbose, "verbose", "v", false, "verbose output")
	cmd.Flags().BoolVarP(&flags.Zero, "zero", "z", false, "whether to generate a zero-byte file")

	rootCmd.AddCommand(cmd)
}
