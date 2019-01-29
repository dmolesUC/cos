package cmd

import (
	"github.com/spf13/pflag"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/internal/protocols"
)

type CosFlags struct {
	Endpoint  string
	Region    string
	Verbose int
}

// TODO: can we share any of Pretty() or String()?

func (f *CosFlags) LogLevel() logging.LogLevel {
	return logging.LogLevel(f.Verbose)
}

func (flags *CosFlags) AddTo(cmdFlags *pflag.FlagSet) {
	cmdFlags.StringVarP(&flags.Endpoint, "endpoint", "e", "", "endpoint: HTTP(S) URL")
	cmdFlags.StringVarP(&flags.Region, "region", "r", "", "S3 region (if not in endpoint URL; default \""+protocols.DefaultAwsRegion+"\")")
	cmdFlags.CountVarP(&flags.Verbose, "verbose", "v", "verbose output (-vv for maximum verbosity)")
}