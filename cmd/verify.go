package cmd

import (
	"github.com/legosx/gopro-media-library-verifier/verifyrun"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verifies the sync with Gopro Media Library",
	Long: `Verify

This command goes over the specified local directory recursively and outputs the files that are not yet uploaded to Gopro Media Library.
`,
	Run: func(cmd *cobra.Command, args []string) {
		runner := verifyrun.NewRunner(
			cmd.Flag("path").Value.String(),
			cmd.Flag("outputFilePath").Value.String(),
			verifyrun.TokenPromptMethod(cmd.Flag("tokenPromptMethod").Value.String()),
		)

		cobra.CheckErr(runner.Run())
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	cobra.CheckErr(verifyrun.Init(verifyCmd))
}
