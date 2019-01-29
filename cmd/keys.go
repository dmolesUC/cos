package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/pkg"
)

const (
	usageKeys = "keys <BUCKET-URL>"

	shortDescKeys = "keys: test creating, retrieving, verifying and deleting potentially problematic keys"

	longDescKeys = shortDescKeys + `

		Creates, retrieves, verifies, and deletes a small object for each 
        value in the Big List of Naughty Strings
        (https://github.com/minimaxir/big-list-of-naughty-strings).
	`
	
	exampleKeys = `
		cos keys s3://www.dmoles.net/ --endpoint https://s3.us-west-2.amazonaws.com/
	`
)

func init() {
	flags := CosFlags{}
	cmd := &cobra.Command{
		Use:           usageKeys,
		Short:         shortDescKeys,
		Long:          logging.Untabify(longDescKeys, ""),
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       logging.Untabify(exampleKeys, "  "),
		RunE: func(cmd *cobra.Command, args []string) error {
			return keys(args[0], flags)
		},
	}
	cmdFlags := cmd.Flags()
	flags.AddTo(cmdFlags)
	rootCmd.AddCommand(cmd)
}

func keys(bucketStr string, f CosFlags) error {
	logger := f.NewLogger()
	logger.Tracef("flags: %v\n", f)
	logger.Tracef("bucket URL: %v\n", bucketStr)

	k := pkg.NewKeys(f.Endpoint, f.Region, bucketStr, logger)
	failures, err := k.CheckAll()
	if err != nil {
		return err
	}
	failureCount := len(failures)
	if failureCount > 0 {
		return fmt.Errorf("%d of %d keys failed", failureCount, k.TotalKeys())
	}
	return nil
}