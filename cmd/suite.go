package cmd

import (
	"time"

	"github.com/dmolesUC3/cos/internal/suite"

	"fmt"

	"github.com/janeczku/go-spinner"
	"github.com/spf13/cobra"

	"github.com/dmolesUC3/cos/internal/logging"
)

type SuiteFlags struct {
	CosFlags
	DryRun bool
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

	//noinspection GoPrintFunctions
	fmt.Println("Starting test suite…\n")

	startAll := time.Now().UnixNano()
	allTasks := suite.AllTasks()
	for index, task := range allTasks {
		title := fmt.Sprintf("%d. %v", index+1, task.Title())

		sp := spinner.StartNew(title)

		var ok bool
		var err error

		start := time.Now().UnixNano()
		if f.DryRun {
			// More framerate sync shenanigans
			time.Sleep(time.Duration(len(sp.Charset)) * sp.FrameRate)
			ok = true
		} else {
			ok, err = task.Invoke(target)
		}
		elapsed := time.Now().UnixNano() - start

		// Lock() / Unlock() around Stop() needed to synchronize cursor movement
		// ..but not always enough (thus the sleep above)
		// TODO: file an issue about this
		sp.Lock()
		sp.Stop()
		sp.Unlock()

		var msgFmt string
		if ok {
			msgFmt = "\u2705 %v: successful (%v)"
		} else {
			msgFmt = "\u274C %v: FAILED (%v)"
		}
		msg := fmt.Sprintf(msgFmt, title, logging.FormatNanos(elapsed))
		fmt.Println(msg)

		if err != nil && f.LogLevel() > logging.Info {
			fmt.Println(err.Error())
		}
	}
	elapsedAll := time.Now().UnixNano() - startAll
	fmt.Printf("\n…test complete (%v).\n", logging.FormatNanos(elapsedAll))

	return nil
}
