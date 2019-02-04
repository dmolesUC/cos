package cmd

import (
	"github.com/spf13/pflag"

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

func (f *CosFlags) AddTo(cmdFlags *pflag.FlagSet) {
	cmdFlags.SortFlags = false

	cmdFlags.StringVarP(&f.Endpoint, "endpoint", "e", "", "HTTP(S) endpoint URL (required)")
	cmdFlags.StringVarP(&f.Region, "region", "r", "", "AWS region (if not in endpoint URL; default \""+objects.DefaultAwsRegion+"\")")
	cmdFlags.CountVarP(&f.Verbose, "verbose", "v", "verbose output (-vv for maximum verbosity)")
}
