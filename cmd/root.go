package cmd

import (
	"fmt"
	"os"

	"github.com/dmolesUC3/cos/internal/objects"
	"github.com/spf13/cobra"

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
        with the `+objects.SwiftUserEnvVar+` and `+objects.SwiftKeyEnvVar+` environment variables.
    ` // TODO: use ST_USER, ST_KEY, whatever the endpoint variable was
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cos",
	Short: shortDescRoot,
	Long:  logging.Untabify(longDescRoot, ""),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}
