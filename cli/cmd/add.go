package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <note text>",
	Short: "Add a new note to your knowledge base.",
	Long: `Create a new note and automatically generate its semantic embedding.

The note will be stored in the database along with its embedding vector,
enabling semantic search across your knowledge base.

Examples:
  synapse add "Einstein's theory of relativity"
  synapse add "Machine learning is a subset of AI"`,

	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		content := args[0]

		if err := noteService.CreateNote(cmd.Context(), content); err != nil {
			return err
		}

		fmt.Printf("Success: Note saved successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
