package com.synapse.controller;

import com.synapse.dto.APIResponse;
import com.synapse.dto.NoteSaveRequest;
import com.synapse.model.Note;
import com.synapse.service.NoteService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

import static com.synapse.constants.MessageConstants.ARTICLE_SAVED;

@RestController
@CrossOrigin(origins = "*")
public class NoteController {

    private final NoteService noteService;

    public NoteController(NoteService noteService){
        this.noteService = noteService;
    }

    @PostMapping("/api/notes")
    public ResponseEntity<APIResponse<Void>> save(@RequestBody NoteSaveRequest noteSaveRequest){
        noteService.saveArticle("darshan", noteSaveRequest);
        return new ResponseEntity<>(APIResponse.success(null, ARTICLE_SAVED), HttpStatus.CREATED);
    }

    @GetMapping("/api/notes")
    public ResponseEntity<APIResponse<List<Note>>> getNotesByUserId(){
        return new ResponseEntity<>(APIResponse.success(noteService.getNotes()), HttpStatus.OK);
    }
}
