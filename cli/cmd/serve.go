package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

type NoteResponse struct {
	ID       int     `json:"id"`
	Content  string  `json:"content"`
	Distance float64 `json:"distance,omitempty"`
}

type AddNoteRequest struct {
	Content string `json:"input"`
}

type SemanticSearchRequest struct {
	Content string `json:"input"`
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Synapse API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		mux := http.NewServeMux()

		mux.HandleFunc("POST /api/notes", handleAddNote)
		mux.HandleFunc("GET /api/notes", handleGetAllNotes)
		mux.HandleFunc("GET /api/notes/{id}", handleGetNoteById)
		mux.HandleFunc("DELETE /api/notes/{id}", handleDeleteNoteById)
		mux.HandleFunc("POST /api/search", handleSemanticSearch)

		port := ":8080"
		slog.Info("Server starting...", "addr", port)

		return http.ListenAndServe(port, mux)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func handleAddNote(w http.ResponseWriter, r *http.Request) {
	var req AddNoteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := noteService.CreateNote(r.Context(), req.Content); err != nil {
		slog.Error("Create failed", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "Note saved successfully"})
}

func handleGetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := noteService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func handleGetNoteById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID must be an integer", http.StatusBadRequest)
		return
	}

	note, err := noteService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if note == nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func handleDeleteNoteById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID must be an integer", http.StatusBadRequest)
		return
	}

	if err := noteService.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Deleted successfully"})
}

func handleSemanticSearch(w http.ResponseWriter, r *http.Request) {
	var req SemanticSearchRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	notes, err := noteService.SemanticSearch(r.Context(), req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}
