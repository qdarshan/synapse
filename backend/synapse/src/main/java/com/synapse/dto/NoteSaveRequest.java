package com.synapse.dto;

public record NoteSaveRequest(
        String url,
        String title,
        String content
) {

}
