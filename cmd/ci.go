package cmd

import (
	"github.com/spf13/cobra"
)

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Management of CI configuration",
	Long: `Management of CI configuration.

You can create the project, branch and build based on a unique "CI configuration" file.

    ontrack-cli ci config

The YAML configuration files is located by default at .yontrack/ci.yaml and this can 
be changed using the --file option:

	ontrack-cli ci config --file .yontrack/ci.yaml

By default, no environment variable is passed to the configuration, for obvious security reasons.
They need to be passed explicitly using the --env options:

	ontrack-cli ci config \
	  --env GIT_URL=git@github.com:nemerosa/ontrack.git \
	  --env GIT_BRANCH=release/5.0

See the Yontrack documentation for the format of the YAML configuration and the list of needed
environment variables.
`,
}

func init() {
	rootCmd.AddCommand(ciCmd)
}
