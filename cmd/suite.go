package cmd

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"code.cloudfoundry.org/bytefmt"

	. "github.com/dmolesUC3/cos/internal/suite"

	"github.com/janeczku/go-spinner"
	"github.com/spf13/cobra"

	"github.com/dmolesUC3/cos/internal/logging"
)

type SuiteFlags struct {
	CosFlags
	SizeMax   string
	CountMax  int64
	SizeOnly  bool
	CountOnly bool
	DryRun    bool
}

func (f *SuiteFlags) sizeMax() (int64, error) {
	sizeStr := f.SizeMax
	sizeIsNumeric := strings.IndexFunc(sizeStr, unicode.IsLetter) == -1
	if sizeIsNumeric {
		return strconv.ParseInt(sizeStr, 10, 64)
	}

	bytes, err := bytefmt.ToBytes(sizeStr)
	if err == nil && bytes > math.MaxInt64 {
		return 0, fmt.Errorf("specified size %d bytes exceeds maximum %d", bytes, math.MaxInt64)
	}
	return int64(bytes), err
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
	cmdFlags.BoolVar(&f.SizeOnly, "size-only", false, "run only file-size tests")
	cmdFlags.BoolVar(&f.CountOnly, "count-only", false, "run only files-per-prefix tests")
	cmdFlags.BoolVarP(&f.DryRun, "dryRun", "n", false, "dry run")
	rootCmd.AddCommand(cmd)
}

func runSuite(bucketStr string, f SuiteFlags) error {
	// TODO: figure out some sensible way to log while spinning
	// logger := logging.DefaultLoggerWithLevel(f.LogLevel())
	// logger.Tracef("flags: %v\n", f)
	// logger.Tracef("bucket URL: %v\n", bucketStr)

	target, err := f.Target(bucketStr)
	if err != nil {
		return err
	}

	sizeMax, err := f.sizeMax()
	if err != nil {
		return err
	}

	var countMax uint64
	if f.CountMax < 0 {
		countMax = math.MaxUint64
	} else {
		countMax = uint64(f.CountMax)
	}

	var tasks []TestTask
	if f.SizeOnly {
		if f.CountOnly {
			return fmt.Errorf("can't specify both --size-only and --count-only")
		} else {
			tasks = SizeTasks(sizeMax)
		}
	} else if f.CountOnly {
		tasks = CountTasks(countMax)
	} else {
		tasks = AllTasks(sizeMax, countMax)
	}

	//noinspection GoPrintFunctions
	fmt.Println("Starting test suite…\n")

	startAll := time.Now().UnixNano()
	for index, task := range tasks {
		title := fmt.Sprintf("%d. %v", index+1, task.Title())

		sp := spinner.NewSpinner(title)
		sp.SetCharset([]string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"})
		sp.Start()

		var ok bool
		var detail string

		start := time.Now().UnixNano()
		if f.DryRun {
			ok = true
		} else {
			ok, detail = task.Invoke(target)
		}
		elapsed := time.Now().UnixNano() - start

		// Lock() / Unlock() around Stop() needed to synchronize cursor movement
		// ..but not always enough (thus the sleep above)
		// TODO: file an issue about this
		sp.Lock()
		sp.Stop()
		sp.Unlock()

		// More framerate sync shenanigans
		time.Sleep(time.Duration(len(sp.Charset)) * sp.FrameRate)

		var msgFmt string
		if ok {
			msgFmt = "\u2705 %v: successful (%v)"
		} else {
			msgFmt = "\u274C %v: FAILED (%v)"
		}
		msg := fmt.Sprintf(msgFmt, title, logging.FormatNanos(elapsed))
		fmt.Println(msg)

		if detail != "" && f.LogLevel() > logging.Info {
			fmt.Println(detail)
		}
	}
	elapsedAll := time.Now().UnixNano() - startAll
	fmt.Printf("\n…test complete (%v).\n", logging.FormatNanos(elapsedAll))

	return nil
}
