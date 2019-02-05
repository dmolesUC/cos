package cmd

import (
	"fmt"
	"io"
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
	`
)

// TODO: more output formats other than --raw and quoted-Go-literal, e.g. --ascii

type keysFlags struct {
	CosFlags

	Raw      bool
	OkFile   string
	BadFile  string
	ListName string
	KeyFile  string

	MemProfile string
}

func (f keysFlags) Pretty() string {
	format := `
		raw:        %v
        okFile:     %v
        badFile:    %v
		listName:   %v
		listFile:	%v
        memprofile: %#v
		region:     %#v
		endpoint:   %#v
		log level:  %v
	`
	format = logging.Untabify(format, "  ")

	// TODO: clean up order of flags in other commands
	return fmt.Sprintf(format,
		f.Raw,
		f.OkFile,
		f.BadFile,
		f.ListName,
		f.KeyFile,

		f.MemProfile,

		f.Region,
		f.Endpoint,
		f.LogLevel(),
	)
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

	keyFile := f.KeyFile
	var keyList keys.KeyList
	if keyFile != "" {
		keyList, err = keys.KeyListForFile(keyFile)
	} else {
		listName := f.ListName
		keyList, err = keys.KeyListForName(listName)
	}
	if err != nil {
		return err
	}

	var okOut io.Writer
	if f.OkFile != "" {
		okOut, err = os.Create(f.OkFile)
		if err != nil {
			return err
		}
	}
	var badOut io.Writer
	if f.BadFile == "" {
		badOut = os.Stdout
	} else {
		badOut, err = os.Create(f.BadFile)
		if err != nil {
			return err
		}
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
