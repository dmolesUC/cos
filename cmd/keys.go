package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/dmolesUC3/cos/internal/keys"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/pkg"
)

const (
	usageKeys = "keys <BUCKET-URL>"

	shortDescKeys = "keys: test the keys supported by an object storage endpoint"

	longDescKeys = shortDescKeys + `

		Creates, retrieves, verifies, and deletes a small object for each value
		in the specified key list. By default, keys outputs only failed keys, to
		standard output, writing each key as a quoted Go string literal
		(see https://golang.org/pkg/strconv/). 

		Use the --raw option to write the keys without quoting or escaping; note
		that this may produce confusing results if any of the keys contain
		newlines.

        Use the --ok option to write successful keys to a file, and the --bad
		option (or shell redirection) to write failed keys to a file instead of to
		standard output.

		Use the --list option to select one of the built-in "standard" key lists.
        Use the --file option to specify a file containing keys to test, one key per
        file, separated by newlines (LF, \n).

        Available lists:
	`

	exampleKeys = `
		cos keys --endpoint https://s3.us-west-2.amazonaws.com/s3://www.dmoles.net/ 
		cos keys --list naughty-strings --endpoint https://s3.us-west-2.amazonaws.com/ s3://www.dmoles.net/  
		cos keys --raw --ok ok.txt --bad bad.txt --endpoint https://s3.us-west-2.amazonaws.com/ s3://www.dmoles.net/
		cos keys --file my-keys.txt --endpoint https://s3.us-west-2.amazonaws.com/ s3://www.dmoles.net/
        cos keys --sample 100 --file my-keys.txt --raw --ok ok.txt --bad bad.txt --endpoint https://s3.us-west-2.amazonaws.com/ s3://www.dmoles.net/ 
	`
)

func longDescription() string {
	listList, err := availableKeyLists()
	if err != nil {
		panic(err)
	}
	longDesc := longDescKeys + "\n" + *listList
	longDescription := logging.Untabify(longDesc, "")
	return longDescription
}

func availableKeyLists() (*string, error) {
	var sb strings.Builder
	w := tabwriter.NewWriter(&sb, 0, 0, 2, ' ', tabwriter.DiscardEmptyColumns)
	for i, list := range keys.KnownKeyLists() {
		_, err := fmt.Fprintf(w, "%d.\t%v\t%v (%d keys)\n", i+1, list.Name(), list.Desc(), list.Count())
		if err != nil {
			return nil, err
		}
	}
	err := w.Flush()
	if err != nil {
		return nil, err
	}
	listList := sb.String()
	return &listList, nil
}

func init() {
	f := keysFlags{}
	cmd := &cobra.Command{
		Use:           usageKeys,
		Short:         shortDescKeys,
		Long:          longDescription(),
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       logging.Untabify(exampleKeys, "  "),
		Run: func(cmd *cobra.Command, args []string) {
			err := checkKeys(args[0], f)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
		},
	}
	cmdFlags := cmd.Flags()
	f.AddTo(cmdFlags)

	cmdFlags.BoolVar(&f.Raw, "raw", false, "write keys in raw (unquoted) format")
	cmdFlags.StringVarP(&f.OkFile, "ok", "o", "", "write successful (\"OK\") keys to specified file")
	cmdFlags.StringVarP(&f.BadFile, "bad", "b", "", "write failed (\"bad\") keys to specified file")
	cmdFlags.StringVarP(&f.ListName, "list", "l", keys.DefaultKeyListName, "key list to check")
	cmdFlags.StringVarP(&f.KeyFile, "file", "f", "", "file of keys to check")
	cmdFlags.IntVarP(&f.Sample, "sample", "s", 0, "sample size, or 0 for all keys")

	rootCmd.AddCommand(cmd)
}

func checkKeys(bucketStr string, f keysFlags) error {
	logger := logging.DefaultLoggerWithLevel(f.LogLevel())
	logger.Tracef("flags: %v\n", f)
	logger.Tracef("bucket URL: %v\n", bucketStr)

	target, err := f.Target(bucketStr)
	if err != nil {
		return err
	}

	keyList, err := f.KeyList()
	if err != nil {
		return err
	}

	okOut, badOut, err := f.Outputs()
	if err != nil {
		return err
	}

	k := pkg.NewKeys(target, keyList)
	failures, err := k.CheckAll(okOut, badOut, f.Raw)
	if err != nil {
		return err
	}
	failureCount := len(failures)
	if failureCount > 0 {
		return fmt.Errorf("%v: %d of %d keys failed", keyList.Name(), failureCount, keyList.Count())
	}
	return nil
}
