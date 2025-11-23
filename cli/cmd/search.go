package cmd

import (
	"fmt"
	"os"
	"synapse/database"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [note text]",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		noteContent := args[0]
		_ = noteContent
		dummyFloats := []float64{0.5, 0.5, 0.5}
		embeddingBytes, err := database.FloatSliceToBytes(dummyFloats)
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
	rootCmd.AddCommand(searchCmd)
}
