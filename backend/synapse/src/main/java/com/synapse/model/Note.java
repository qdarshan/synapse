package com.synapse.model;

import lombok.Data;

import java.time.ZonedDateTime;

@Data
public class Note {
    private String id;
    private String userID;
    private String url;
    private String title;
    private String content;
    private ZonedDateTime savedAt;

    public Note(String id, String userID, String url, String title, String content){
        this.id = id;
        this.userID = userID;
        this.url = url;
        this.title = title;
        this.content = content;
        this.savedAt = ZonedDateTime.now();
    }
}
