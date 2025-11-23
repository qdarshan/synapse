package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <note-id>",
	Short: "Delete a note by ID.",
	Long: `Permanently remove a note and its embedding vector from the database.

This action cannot be undone. Use 'synapse search <query> --id' to find the ID
of a note before deleting it.

Examples:
  synapse delete 42
  synapse delete 7`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		noteID, err := strconv.Atoi(id)
		if err != nil {
			return fmt.Errorf("invalid ID format: %s is not a valid integer. Please provide a numeric ID", id)
		}
		if err := dbManager.DeleteNote(noteID); err != nil {
			return fmt.Errorf("failed to delete note: %w", err)
		}
		fmt.Printf("Success: Note deleted successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
