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

func (f *CosFlags) AddTo(cmdFlags *pflag.FlagSet) {
	cmdFlags.StringVarP(&f.Endpoint, "endpoint", "e", "", "endpoint: HTTP(S) URL")
	cmdFlags.StringVarP(&f.Region, "region", "r", "", "S3 region (if not in endpoint URL; default \""+protocols.DefaultAwsRegion+"\")")
	cmdFlags.CountVarP(&f.Verbose, "verbose", "v", "verbose output (-vv for maximum verbosity)")
}

func (f *CosFlags) NewLogger() logging.Logger {
	return logging.NewLogger(f.LogLevel())
}