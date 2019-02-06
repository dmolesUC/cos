package cmd

import (
	"github.com/dmolesUC3/cos/internal/suite"

	"fmt"

	"github.com/janeczku/go-spinner"
	"github.com/spf13/cobra"

	"github.com/dmolesUC3/cos/internal/logging"
)

func init() {
	f := CosFlags{}
	cmd := &cobra.Command{
		Use:   "suite <BUCKET-URL>",
		Short: "run a suite of tests",
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSuite(args[0], f)
		},
	}
	cmdFlags := cmd.Flags()
	f.AddTo(cmdFlags)
	rootCmd.AddCommand(cmd)
}

func runSuite(bucketStr string, f CosFlags) error {
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

	allTasks := suite.AllTasks()
	for index, task := range allTasks {
		title := fmt.Sprintf("%d. %v", index+1, task.Title())

		s := spinner.StartNew(title)
		ok, err := task.Invoke(target)
		s.Stop()

		if ok {
			fmt.Printf("\u2705 %v: successful\n", title)
		} else {
			fmt.Printf("\u274C %v: FAILED\n", title)
		}

		if err != nil && f.LogLevel() > logging.Info {
			fmt.Println(err.Error())
		}
	}
	fmt.Println("\n…test complete.")

	return nil
}
