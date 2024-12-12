/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize an ease.yaml file in the current directory",
	Long: `The init command creates an ease.yaml configuration file
in the current working directory, using a default template.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Define the output file name
		outputFile := "ease.yaml"

		// Check if the file already exists
		if _, err := os.Stat(outputFile); err == nil {
			return fmt.Errorf("file %s already exists", outputFile)
		}

		// Get the path to the template file
		templatePath := filepath.Join("internal", "static", "ease-example.yaml")

		// Load and parse the template
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		// Create a buffer to hold the template output
		var buf bytes.Buffer

		// Execute the template (no data needed for this example)
		if err := tmpl.Execute(&buf, nil); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		// Write the generated file
		if err := os.WriteFile(outputFile, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		fmt.Printf("Successfully created %s\n", outputFile)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
