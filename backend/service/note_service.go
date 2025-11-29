package service

import (
	"context"
	"fmt"
	"synapse/client"
	"synapse/database"
)

type NoteService struct {
	DBManager *database.SQLiteManager
}

func NewNoteService(dbManager *database.SQLiteManager) *NoteService {
	return &NoteService{
		DBManager: dbManager,
	}
}

func (s *NoteService) CreateNote(ctx context.Context, content string) error {
	embeddingFloats, err := client.GenerateEmbedding(ctx, content)
	if err != nil {
		return fmt.Errorf("AI generation failed: %w", err)
	}

	embeddingBytes, err := database.FloatSliceToBytes(embeddingFloats)
	if err != nil {
		return fmt.Errorf("vector encoding failed: %w", err)
	}

	note := database.Note{
		Content:         content,
		EmbeddingVector: embeddingBytes,
	}

	if err := s.DBManager.SaveNote(note); err != nil {
		return fmt.Errorf("db save failed: %w", err)
	}
	return nil
}

func (s *NoteService) SemanticSearch(ctx context.Context, query string) ([]database.Note, error) {
	embeddingFloats, err := client.GenerateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("AI generation failed: %w", err)
	}

	embeddingBytes, err := database.FloatSliceToBytes(embeddingFloats)
	if err != nil {
		return nil, fmt.Errorf("vector encoding failed: %w", err)
	}

	notes, err := s.DBManager.SearchNotes(embeddingBytes)
	if err != nil {
		return nil, fmt.Errorf("db search failed: %w", err)
	}
	return notes, nil
}

func (s *NoteService) GetAll() ([]database.Note, error) {
	return s.DBManager.GetAllNotes()
}

func (s *NoteService) GetByID(id int) (*database.Note, error) {
	return s.DBManager.GetNoteById(id)
}

func (s *NoteService) Delete(id int) error {
	return s.DBManager.DeleteNote(id)
}
