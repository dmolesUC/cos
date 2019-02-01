package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/dmolesUC3/cos/internal/keys"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/internal/objects"
	"github.com/dmolesUC3/cos/internal/streaming"
	"github.com/dmolesUC3/cos/pkg"
)

const (
	usageKeys = "keys <BUCKET-URL>"

	shortDescKeys = "keys: test creating, retrieving, verifying and deleting potentially problematic keys"

	longDescKeys = shortDescKeys + `

		Creates, retrieves, verifies, and deletes a small object for each 
        value in the specified key list.

        Available lists:
	`

	exampleKeys = `
		cos keys s3://www.dmoles.net/ --endpoint https://s3.us-west-2.amazonaws.com/
	`
)

type keysFlags struct {
	CosFlags

	From     int
	To       int
	ListName string

	MemProfile string
}

func (f keysFlags) Pretty() string {
	format := `
		log level:  %v
		region:     %#v
		endpoint:   %#v
		from:       %d
        to:         %d
        memprofile: %#v
	`
	format = logging.Untabify(format, "  ")

	return fmt.Sprintf(format, f.LogLevel(), f.Region, f.Endpoint, f.From, f.To, f.MemProfile)
}

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
	for i, list := range keys.AllKeyLists() {
		_, err := fmt.Fprintf(w, "%d.\t%v\t%v\n", i + 1, list.Name(), list.Desc())
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

	cmdFlags.IntVarP(&f.From, "from", "f", 1, "first key to check (1-indexed, inclusive)")
	cmdFlags.IntVarP(&f.To, "to", "t", -1, "last key to check (1-indexed, inclusive); -1 to check all keys in list")
	cmdFlags.StringVarP(&f.ListName, "list", "l", keys.DefaultKeyListName, "key list to check")

	cmdFlags.StringVarP(&f.MemProfile, "memprofile", "", "", "write memory profile to `file`")

	rootCmd.AddCommand(cmd)
}

func checkKeys(bucketStr string, f keysFlags) error {
	defer func() {
		if f.MemProfile != "" {
			f, err := os.Create(f.MemProfile)
			if err != nil {
				log.Fatal("could not create memory profile: ", err)
			}
			runtime.GC() // get up-to-date statistics
			if err := pprof.WriteHeapProfile(f); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
			err = f.Close()
			if err != nil {
				log.Fatal("could not close memory profile: ", err)
			}
		}
	}()

	logger := logging.DefaultLoggerWithLevel(f.LogLevel())
	logger.Tracef("flags: %v\n", f)
	logger.Tracef("bucket URL: %v\n", bucketStr)

	endpointURL, err := streaming.ValidAbsURL(f.Endpoint)
	if err != nil {
		return err
	}

	bucketURL, err := streaming.ValidAbsURL(bucketStr)
	if err != nil {
		return err
	}

	target, err := objects.NewTarget(endpointURL, bucketURL, f.Region)
	if err != nil {
		return err
	}

	listName := f.ListName
	keyList, err := keys.KeyListForName(listName)
	if err != nil {
		return err
	}

	startIndex := f.From - 1
	endIndex := f.To
	if endIndex <= 0 {
		endIndex = keyList.Count()
	}
	logger.Tracef("list: %v, startIndex: %d, endIndex: %d\n", listName, startIndex, endIndex)

	k := pkg.NewKeys(target, keyList)
	failures, err := k.CheckAll(startIndex, endIndex)
	if err != nil {
		return err
	}
	failureCount := len(failures)
	if failureCount > 0 {
		totalExpected := endIndex - startIndex
		return fmt.Errorf("%v: %d of %d keys failed", listName, failureCount, totalExpected)
	}
	return nil
}
