package cmd

import (
	"fmt"
	"io"
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
	CosFlags

	Key      string
	Size     string
	Seed     int64
	Zero     bool
	Keep     bool
}

func (f crvdFlags) ContentLength() (int64, error) {
	if f.Zero {
		return 0, nil
	}
	sizeIsNumeric := strings.IndexFunc(f.Size, unicode.IsLetter) == -1
	if sizeIsNumeric {
		return strconv.ParseInt(f.Size, 10, 64)
	}

	bytes, err := bytefmt.ToBytes(f.Size)
	if err == nil && bytes > math.MaxInt64 {
		return 0, fmt.Errorf("specified size %d bytes exceeds maximum %d", bytes, math.MaxInt64)
	}
	return int64(bytes), err
}

func (f crvdFlags) Pretty() string {
	format := `
		log level: %v
		region:   '%v'
		endpoint: '%v'
        key:      '%v'
		size:      %v (%d bytes)
        seed:      %d
        zero:      %v
        keep:      %v`
	format = logging.Untabify(format, "  ")

	contentLength, _ := f.ContentLength()

	return fmt.Sprintf(format, f.LogLevel(), f.Region, f.Endpoint, f.Key, f.Size, contentLength, f.Seed, f.Zero, f.Keep)
}

func crvd(bucketStr string, f crvdFlags) (err error) {

	var logger = logging.NewLogger(f.LogLevel())
	logger.Tracef("flags: %v\n", f)
	logger.Tracef("bucket URL: %v\n", bucketStr)

	if f.Key == "" {
		f.Key = fmt.Sprintf("cos-crvd-%d.bin", time.Now().Unix())
	}
	if f.Endpoint == "" {
		return fmt.Errorf("endpoint URL must be specified")
	}

	bucketUrl, err := objects.ValidAbsURL(bucketStr)
	if err != nil {
		return err
	}

	obj, err := objects.NewObjectBuilder().
		WithEndpointStr(f.Endpoint).
		WithRegion(f.Region).
		WithProtocolUri(bucketUrl, logger).
		WithKey(f.Key).
		Build(logger)
	if err != nil {
		return err
	}

	contentLength, err := f.ContentLength()
	if err != nil {
		return err
	}

	var crvd = pkg.Crvd{Object: obj}
	random := rand.New(rand.NewSource(f.Seed))
	body := io.LimitReader(random, contentLength)

	if f.Keep {
		err = crvd.CreateRetrieveVerify(body, contentLength)
	} else {
		err = crvd.CreateRetrieveVerifyDelete(body, contentLength)
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
	cmdFlags := cmd.Flags()
	flags.AddTo(cmdFlags)

	cmdFlags.StringVarP(&flags.Size, "size", "s", "1K", "size in bytes of object to create, if --zero not set")
	cmdFlags.BoolVarP(&flags.Zero, "zero", "z", false, "whether to generate a zero-byte file") // TODO: replace with NoOptDefaultVal
	cmdFlags.Int64VarP(&flags.Seed, "random-seed", "", 0, "seed for random-number generator (default 0)")
	cmdFlags.StringVarP(&flags.Key, "key", "k", "", "key to create (defaults to cos-crvd-TIMESTAMP.bin)")
	cmdFlags.BoolVarP(&flags.Keep, "keep", "", false, "keep object after verification (defaults to false)")

	rootCmd.AddCommand(cmd)
}

