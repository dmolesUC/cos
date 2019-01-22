package cmd

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"

	"code.cloudfoundry.org/bytefmt"
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

        The size may be specified as an exact number of bytes, or using human-readable
        quantities such as "5K" (4 KiB or 4096 bytes), "3.5M" (3.5 MiB or 3670016 bytes),
        etc. The units supported are bytes (B), binary kilobytes (K, KB, KiB), 
        binary megabytes (M, MB, MiB), binary gigabytes (G, GB, GiB), and binary 
        terabytes (T, TB, TiB). If no unit is specified, bytes are assumed.

        Random bytes are generated using the Go default random number generator, with
        a default seed of 0, for repeatability. An alternative seed can be specified
        with the --random-seed flag.
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
	Size     string
	Seed     int64
	Zero     bool
	Keep     bool

	SizeBytes int64
}

func (f crvdFlags) Pretty() string {
	format := `
		verbose:   %v
		region:   '%v'
		endpoint: '%v'
        key:      '%v'
		size:      %v (%d bytes)
        seed:      %d
        zero:      %v`
	format = logging.Untabify(format, "  ")
	return fmt.Sprintf(format, f.Verbose, f.Region, f.Endpoint, f.Key, f.Size, f.SizeBytes, f.Seed, f.Zero)
}

func crvd(bucketStr string, flags crvdFlags) (err error) {
	if flags.Key == "" {
		flags.Key = fmt.Sprintf("cos-crvd-%d.bin", time.Now().Unix())
	}
	if flags.Zero {
		flags.SizeBytes = 0
	} else {
		if strings.IndexFunc(flags.Size, unicode.IsLetter) == -1 {
			flags.SizeBytes, err = strconv.ParseInt(flags.Size, 10, 64)
			if err != nil {
				return err
			}
		} else {
			sizeBytes, err2 := bytefmt.ToBytes(flags.Size)
			if err2 != nil {
				return err2
			}
			if sizeBytes >= math.MaxInt64 {
				return fmt.Errorf("specified size %d bytes exceeds maximum %d", sizeBytes, math.MaxInt64)
			}
			flags.SizeBytes = int64(sizeBytes)
		}
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
		digest, err = crvd.CreateRetrieveValidate(random, flags.SizeBytes)
		if err == nil {
			logger.Infof("created %v (%d bytes, SHA-256 digest %x)\n", objects.ProtocolUriStr(obj), flags.Size, digest)
		}
	} else {
		digest, err = crvd.CreateRetrieveValidateDelete(random, flags.SizeBytes)
		if err == nil {
			logger.Infof("verified and deleted %v (%d bytes, SHA-256 digest %x)\n", objects.ProtocolUriStr(obj), flags.SizeBytes, digest)
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
	cmd.Flags().Int64VarP(&flags.Seed, "random-seed", "", 0, "seed for random-number generator")
	cmd.Flags().StringVarP(&flags.Region, "region", "r", "", "S3 region (if not in endpoint URL; default \""+protocols.DefaultAwsRegion+"\")")
	cmd.Flags().StringVarP(&flags.Size, "size", "s", "1K", "size in bytes of object to create, if --zero not set")
	cmd.Flags().BoolVarP(&flags.Verbose, "verbose", "v", false, "verbose output")
	cmd.Flags().BoolVarP(&flags.Zero, "zero", "z", false, "whether to generate a zero-byte file")

	rootCmd.AddCommand(cmd)
}
