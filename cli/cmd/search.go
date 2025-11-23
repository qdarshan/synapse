package cmd

import (
	"fmt"
	"os"
	"strconv"
	"synapse/client"
	"synapse/database"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var searchById bool
var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search notes by semantic similarity or by ID.",
	Long: `Search your notes database in two ways:

Semantic Search (default):
  Converts your query text to an embedding and finds the top 10 most similar notes.
  Results are sorted by distance (lower distance = higher similarity).

ID Search (with --id flag):
  Retrieves a specific note using its numeric ID. When using this flag, provide
  a numeric argument instead of text.

Examples:
  synapse search "quantum mechanics"           # Semantic search
  synapse search "deep learning"               # Semantic search
  synapse search 42 --id                       # Get note with ID 42
  synapse search 7 -i                          # Short flag: get note with ID 7`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]

		if searchById {
			noteId, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("error: Id search requires a numeric argument, received '%s'", input)
			}

			note, err := dbManager.GetNoteById(noteId)
			if err != nil {
				return fmt.Errorf("failed to retrieve note by ID %d: %w", noteId, err)
			}

			if note != nil {
				fmt.Printf("Note Found (ID: %d)\nContent:\n%s\n", note.Id, note.Content)
			} else {
				fmt.Printf("Note with ID %d not found.\n", noteId)
			}
			return nil
		}

		embeddingFloats, err := client.GenerateEmbedding(input)
		if err != nil {
			return fmt.Errorf("failed to get embeddings from LMStudio: %w", err)
		}

		embeddingBytes, err := database.FloatSliceToBytes(embeddingFloats)
		if err != nil {
			return fmt.Errorf("failed to encode embedding vector: %w", err)
		}

		notes, err := dbManager.SearchNotes(embeddingBytes)
		if err != nil {
			return fmt.Errorf("failed to search the note: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tDISTANCE\tCONTENT")
		fmt.Fprintln(w, "--\t--------\t-------")

		for _, note := range notes {
			fmt.Fprintf(w, "%d\t%.4f\t%s\n", note.Id, note.Distance, note.Content)
		}
		w.Flush()

		return nil
	},
}

func init() {
	searchCmd.Flags().BoolVarP(&searchById, "id", "i", false, "Search by exact Note ID instead of content.")
	rootCmd.AddCommand(searchCmd)
}
