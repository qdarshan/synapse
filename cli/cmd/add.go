package cmd

import (
	"fmt"
	"synapse/database"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [note text]",
	Short: "Adds a new note and generates its embedding.",
	RunE: func(cmd *cobra.Command, args []string) error {
		noteContent := args[0]

		dummyFloats := []float64{0.5, 0.5, 0.5}
		embeddingBytes, err := database.FloatSliceToBytes(dummyFloats)
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
