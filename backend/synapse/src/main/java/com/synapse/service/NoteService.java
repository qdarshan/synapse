package com.synapse.service;

import com.synapse.dto.NoteSaveRequest;
import com.synapse.model.Note;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

@Service
public class NoteService {

    private List<Note> notes;

    public NoteService(){
        notes = new ArrayList<>();
    }

    public void saveArticle(String userID, NoteSaveRequest noteSaveRequest){
        Note note = new Note(UUID.randomUUID().toString(), userID,
                noteSaveRequest.url(), noteSaveRequest.title(), noteSaveRequest.content());
        notes.add(note);
    }

    public List<Note> getNotes() {
        return notes;
    }
}
