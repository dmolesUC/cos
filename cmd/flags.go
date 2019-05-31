package cmd

import (
	"github.com/spf13/pflag"

	"github.com/dmolesUC3/cos/internal/streaming"

	"github.com/dmolesUC3/cos/internal/objects"

	"github.com/dmolesUC3/cos/internal/logging"
)

type CosFlags struct {
	Endpoint  string
	Region    string
	Verbose int
}

func (f *CosFlags) LogLevel() logging.LogLevel {
	return logging.LogLevel(f.Verbose)
}

// Deprecated TODO: use rootCmd.PersistentFlags() instead
func (f *CosFlags) AddTo(cmdFlags *pflag.FlagSet) {
	cmdFlags.SortFlags = false

	cmdFlags.StringVarP(&f.Endpoint, "endpoint", "e", "", "HTTP(S) endpoint URL (required)")
	cmdFlags.StringVarP(&f.Region, "region", "r", "", "AWS region (if not in endpoint URL; default \""+objects.DefaultAwsRegion+"\")")
	cmdFlags.CountVarP(&f.Verbose, "verbose", "v", "verbose output (-vv for maximum verbosity)")
}

func (f *CosFlags) Target(bucketStr string) (objects.Target, error) {
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