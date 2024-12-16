package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dangrondahl/ease/internal/ease"
	"github.com/spf13/cobra"
)

// initCmd represents the init command.
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new ease.yaml configuration file",
	Long:  `The init command generates a new ease.yaml configuration file in the root directory of your project.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Get the working directory
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}

		// Define the output file path
		outputPath := filepath.Join(cwd, "ease.yaml")

		/*err = ease.New().
		WithGitOps("my-organization", "my-repo").
		WithTemplate("https://github.com/example/helmrelease-templates.git", "main", "templates/helmrelease.yaml").
		CreateEaseFile(outputPath)
		*/
		err = ease.NewFromPrompt().CreateEaseFile(outputPath)

		if err != nil {
			return err
		}

		fmt.Printf("Successfully created ease.yaml in %s\n", cwd)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
