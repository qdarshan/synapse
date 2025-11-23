package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"synapse/database"

	"github.com/spf13/cobra"
)

var dbManager *database.SQLiteManager
var dbFilepath = "synapse.db"

var rootCmd = &cobra.Command{
	Use:   "synapse",
	Short: "Synapse: A high-performance local notes and embedding tool.",
	Long:  `Synapse allows you to capture notes, generate embeddings, and search your knowledge base semantically.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		database.RegisterCustomDriver()

		var err error
		dbManager, err = database.Initialize("synapse.db")
		if err != nil {
			return fmt.Errorf("Failed to initialize database: %w", err)
		}
		if err := dbManager.SetupSchema(); err != nil {
			return fmt.Errorf("failed to setup schema: %w", err)
		}

		slog.Debug("Database initialized and ready.", "path", dbFilepath)

		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if dbManager != nil && dbManager.DB != nil {
			dbManager.DB.Close()
			slog.Debug("Database connection closed.")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
