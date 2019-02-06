package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/dmolesUC3/cos/internal/keys"

	"github.com/dmolesUC3/cos/internal/logging"

	"github.com/dmolesUC3/cos/internal/objects"
	"github.com/dmolesUC3/cos/internal/streaming"
)

type keysFlags struct {
	CosFlags

	// TODO: more output formats other than --raw and quoted-Go-literal, e.g. --ascii
	Raw      bool
	OkFile   string
	BadFile  string
	ListName string
	KeyFile  string
	Sample   int
}

func (f *keysFlags) Pretty() string {
	format := `
		raw:        %v
        okFile:     %v
        badFile:    %v
		listName:   %v
		listFile:	%v
		sample:     %d
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
		f.Sample,

		f.Region,
		f.Endpoint,
		f.LogLevel(),
	)
}

func (f *keysFlags) KeyList() (keyList keys.KeyList, err error) {
	if f.KeyFile == "" {
		keyList, err = keys.KeyListForName(f.ListName)
	} else {
		keyList, err = keys.KeyListForFile(f.KeyFile)
	}
	if err == nil && f.Sample > 0 {
		keyList, err = keys.SamplingKeyList(keyList, f.Sample)
	}
	return keyList, err
}

func (f *keysFlags) Target(bucketStr string) (objects.Target, error) {
	endpointURL, err := streaming.ValidAbsURL(f.Endpoint)
	if err != nil {
		return nil, err
	}

	bucketURL, err := streaming.ValidAbsURL(bucketStr)
	if err != nil {
		return nil, err
	}

	return objects.NewTarget(endpointURL, bucketURL, f.Region)
}

func (f *keysFlags) Outputs() (okOut io.Writer, badOut io.Writer, err error) {
	if f.OkFile != "" {
		okOut, err = os.Create(f.OkFile)
		if err != nil {
			return nil, nil, err
		}
	}
	if f.BadFile == "" {
		badOut = os.Stdout
	} else {
		badOut, err = os.Create(f.BadFile)
		if err != nil {
			// if we opened the okFile we now need to close it
			if okOutC, ok := okOut.(io.WriteCloser); ok {
				//noinspection GoUnhandledErrorResult
				defer okOutC.Close()
			}
			return nil, nil, err
		}
	}
	return
}
