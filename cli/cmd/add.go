package cmd

import (
	"fmt"
	"synapse/client"
	"synapse/database"

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
		noteContent := args[0]

		embeddingFloats, err := client.GenerateEmbedding(noteContent)
		if err != nil {
			return fmt.Errorf("failed to get embeddings from LMStudio: %w", err)
		}

		embeddingBytes, err := database.FloatSliceToBytes(embeddingFloats)
		if err != nil {
			return fmt.Errorf("failed to encode embedding vector: %w", err)
		}

		note := database.Note{
			Content:         noteContent,
			EmbeddingVector: embeddingBytes,
		}

		if err := dbManager.SaveNote(note); err != nil {
			return fmt.Errorf("failed to save note to database: %w", err)
		}
		fmt.Printf("Success: Note saved successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
