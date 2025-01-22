package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/wyattcupp/codebase-tool/internal/clipboard"
	"github.com/wyattcupp/codebase-tool/internal/collector"

	"github.com/spf13/cobra"
)

var (
	outputFile  string
	copyToClip  bool
	targetDir   string
	ignoreDirs  []string
	codebaseCmd = &cobra.Command{
		Use:   "codebase",
		Short: "Collect codebase context and output to a markdown file or clipboard.",
		Long: `codebase is a CLI tool that recursively reads your project's source files,
skips files/directories listed in .codebase_ignore (and those passed via flags),
and aggregates them into a single markdown-formatted string for usage in LLMs.
You can save that markdown to a file or copy it directly to your clipboard.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if targetDir == "" {
				wd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get working directory: %v", err)
				}
				targetDir = wd
			} else {
				tdAbs, err := filepath.Abs(targetDir)
				if err != nil {
					return fmt.Errorf("failed to resolve absolute path for %s: %v", targetDir, err)
				}
				targetDir = tdAbs
			}

			// collect the entire codebase
			result, tokenCount, err := collector.CollectCodebase(targetDir, ignoreDirs)

			if err != nil {
				return fmt.Errorf("failed collecting codebase: %v", err)
			}

			// copy to clipboard
			if copyToClip {
				if err := clipboard.WriteClipboard(result); err != nil {
					return fmt.Errorf("failed copying to clipboard: %v", err)
				}
				fmt.Println("Codebase context copied to clipboard successfully!\n\nTokens:", tokenCount)
			}

			// output to file
			if outputFile != "" {
				if err := os.WriteFile(outputFile, []byte(result), 0644); err != nil {
					return fmt.Errorf("failed writing to markdown file: %v", err)
				}
				log.Printf("Codebase context saved to %s\n\nTokens:%v", outputFile, tokenCount)
			}

			if !copyToClip && outputFile == "" {
				fmt.Println("No output method selected. Use --help to see options.")
			}
			return nil
		},
	}
)

// EstimateTokens provides a ceiling estimate of the token count for a given text.
func EstimateTokens(text string) int {
	// Use a regex to approximate tokenization (split on spaces and punctuation).
	re := regexp.MustCompile(`\w+|[^\w\s]`)
	tokens := re.FindAllString(text, -1)

	// Return the length of the tokens slice as the token count.
	return len(tokens)
}

func Execute() error {
	return codebaseCmd.Execute()
}

func init() {
	codebaseCmd.Flags().StringVarP(&targetDir, "dir", "d", "", "Target directory to scan (defaults to current working directory).")
	codebaseCmd.Flags().StringVarP(&outputFile, "out", "o", "", "Markdown file to write the aggregated code context to.")
	codebaseCmd.Flags().BoolVarP(&copyToClip, "clipboard", "c", false, "Copy the aggregated code context to clipboard (overwrites it!).")
	codebaseCmd.Flags().StringArrayVarP(&ignoreDirs, "ignore", "i", []string{}, "Additional directories/files to ignore (besides .codebase_ignore).")
}
