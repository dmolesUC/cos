package cmd

import (
	"github.com/dmolesUC3/cos/pkg"

	"fmt"

	"github.com/janeczku/go-spinner"
	"github.com/spf13/cobra"

	"github.com/dmolesUC3/cos/internal/logging"
)

func init() {
	f := CosFlags{}
	cmd := &cobra.Command{
		Use: "suite",
		Short: "run a suite of tests",
		RunE: func(cmd *cobra.Command, args []string) error {
			return suite(args[0], f)
		},
	}
	cmdFlags := cmd.Flags()
	f.AddTo(cmdFlags)
	rootCmd.AddCommand(cmd)
}

func suite(bucketStr string, f CosFlags) error {
	// TODO: figure out some sensible way to log while spinning
	// logger := logging.DefaultLoggerWithLevel(f.LogLevel())
	// logger.Tracef("flags: %v\n", f)
	// logger.Tracef("bucket URL: %v\n", bucketStr)

	target, err := f.Target(bucketStr)
	if err != nil {
		return err
	}

	fmt.Println("Starting test suite…\n")

	// TODO: track and output time for each test & total time
	title := "1. create, retrieve, verify, delete"
	s := spinner.StartNew(title)
	crvd := pkg.NewDefaultCrvd(target, "")
	err = crvd.CreateRetrieveVerifyDelete()
	s.Stop()

	var indicator string
	var result string
	if err == nil {
		indicator = "\u2705"
		result = "successful"
	} else {
		indicator = "\u274C"
		result = fmt.Sprintf("FAILED")
	}
	
	fmt.Printf("%v %v: %v\n", indicator, title, result)
	if err != nil && f.LogLevel() > logging.Info {
		fmt.Println(err.Error())
	}


	fmt.Println("\n…test complete.")

	return nil
}
