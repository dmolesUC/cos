package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dmolesUC3/cos/internal/objects"

	"github.com/dmolesUC3/cos/internal/logging"
)

const (
	shortDescRoot = "cos: a cloud object storage tool"
	longDescRoot  = shortDescRoot + `

        cos is a tool for testing and validating cloud object storage.

		For S3, cos uses AWS credentials from ~/.aws/config, or any alternate config
        file specified with the AWS_CONFIG_FILE environment variable. Credentials can
        also be specified with the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment
		variables. When running in an EC2 environment, cos can also access IAM role
        credentials.

        Note that for OpenStack Swift, the API username and key must be specified
        with the ` + objects.SwiftUserEnvVar + ` and ` + objects.SwiftKeyEnvVar + ` environment variables.
    `
)

var rootCmd = &cobra.Command{
	Use:   "cos",
	Short: shortDescRoot,
	Long:  logging.Untabify(longDescRoot, ""),
}

func Execute() error {
	return rootCmd.Execute()
}