package cmd

import (
	"fmt"
	"math"

	"code.cloudfoundry.org/bytefmt"

	. "github.com/dmolesUC3/cos/internal/suite"

	"github.com/spf13/cobra"

	"github.com/dmolesUC3/cos/internal/logging"
)

type SuiteFlags struct {
	CosFlags
	SizeMax   string
	CountMax  int64
	DryRun    bool
}

func init() {
	f := SuiteFlags{}
	cmd := &cobra.Command{
		Use:   "suite <BUCKET-URL>",
		Short: "run a suite of tests",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSuite(args[0], f)
		},
	}
	cmdFlags := cmd.Flags()
	f.AddTo(cmdFlags)

	// TODO: document these
	sizeMaxDefault := bytefmt.ByteSize(256 * bytefmt.GIGABYTE)
	cmdFlags.StringVarP(&f.SizeMax, "size-max", "s", sizeMaxDefault, "max file size to create")
	cmdFlags.Int64VarP(&f.CountMax, "count-max", "c", -1, "max number of files to create, or -1 for no limit")
	cmdFlags.BoolVarP(&f.DryRun, "dryRun", "n", false, "dry run")
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

	//noinspection GoPrintFunctions
	fmt.Println("Starting test suite…\n")
	suite := AllCases(sizeMax, countMax, target, f.LogLevel(), f.DryRun)
	elapsedAll := suite.Execute()
	fmt.Printf("\n…test complete (%v).\n", logging.FormatNanos(elapsedAll))

	return nil
}
