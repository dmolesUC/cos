package cmd

import (
	"fmt"
	"math"

	"code.cloudfoundry.org/bytefmt"
	"github.com/dmolesUC3/cos/pkg"

	. "github.com/dmolesUC3/cos/internal/suite"

	"github.com/spf13/cobra"

	"github.com/dmolesUC3/cos/internal/logging"
)

type SuiteFlags struct {
	CosFlags

	Size bool
	SizeMax   string

	Count bool
	CountMax  uint64

	Unicode bool

	DryRun    bool
}

const (
	suiteLongDesc = `
		Run a suite of test cases investigating various possible limitations of a
		cloud storage service:

		- maximum file size (--size)
		- maximum number of files per key prefix (--count)
		- Unicode key support (--unicode)

		If none of --size, --count, etc. is specified, all test cases are run.

		The maximum size may be specified as an exact number of bytes, or using
		human-readable quantities such as "5K" (4 KiB or 4096 bytes), "3.5M" (3.5
		MiB or 3670016 bytes), etc. The units supported are bytes (B), binary
		kilobytes (K, KB, KiB), binary megabytes (M, MB, MiB), binary gigabytes (G,
		GB, GiB), and binary terabytes (T, TB, TiB). If no unit is specified, bytes
		are assumed.
	`
)

func init() {
	f := SuiteFlags{}
	cmd := &cobra.Command{
		Use:   "suite <BUCKET-URL>",
		Short: "run a suite of tests",
		Long: logging.Untabify(suiteLongDesc, ""),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSuite(args[0], f)
		},
	}
	cmdFlags := cmd.Flags()
	f.AddTo(cmdFlags)

	cmdFlags.BoolVarP(&f.Size, "size", "s", false, "test file sizes")
	cmdFlags.StringVar(&f.SizeMax, "size-max", bytefmt.ByteSize(SizeMaxDefault), "max file size to create")

	cmdFlags.BoolVarP(&f.Count, "count", "c", false, "test file counts")
	cmdFlags.Uint64Var(&f.CountMax, "count-max", CountMaxDefault, "max number of files to create, or -1 for no limit")

	cmdFlags.BoolVarP(&f.Unicode, "unicode", "u", false, "test Unicode keys")

	cmdFlags.BoolVarP(&f.DryRun, "dry-run", "n", false, "dry run; list all tests that would be run, but don't create any files")
	rootCmd.AddCommand(cmd)
}

func runSuite(bucketStr string, f SuiteFlags) error {
	// TODO: figure out some sensible way to log while spinning
	// logger := logging.DefaultLoggerWithLevel(f.LogLevel())
	// logger.Tracef("flags: %v\n", f)
	// logger.Tracef("bucket URL: %v\n", bucketStr)

	sizeMax, err := ParseSizeMax(f.SizeMax)
	if err != nil {
		return err
	}

	var countMax uint64
	if f.CountMax < 0 {
		countMax = math.MaxUint64
	} else {
		countMax = uint64(f.CountMax)
	}

	target, err := f.Target(bucketStr)
	if err != nil {
		return err
	}

	logLevel := f.LogLevel()
	if logLevel > logging.Detail {
		_ = logging.DefaultLoggerWithLevel(logLevel)
	}

	var cases []Case
	runAllCases := !(f.Size || f.Count || f.Unicode)
	if runAllCases || f.Size {
		cases = append(cases, FileSizeCases(sizeMax)...)
	}
	if runAllCases || f.Count {
		cases = append(cases, FileCountCases(countMax)...)
	}
	if runAllCases || f.Unicode {
		cases = append(cases, AllUnicodeCases()...)
	}

	// sanity check
	fmt.Println("Checking server connection…")
	if !f.DryRun {
		crvd := pkg.NewDefaultCrvd(target, "")
		err := crvd.CreateRetrieveVerifyDelete()
		if err != nil {
			return fmt.Errorf("connection check failed: %v", err)
		}
	}

	//noinspection GoPrintFunctions
	fmt.Printf("Starting test suite (%d cases)…\n\n", len(cases))
	suite := NewSuite(cases, target, logLevel, f.DryRun)
	elapsedAll := suite.Execute()
	fmt.Printf("\n…test complete (%v).\n", logging.FormatNanos(elapsedAll))

	return nil
}
