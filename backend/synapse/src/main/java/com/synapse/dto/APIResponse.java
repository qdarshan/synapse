package com.synapse.dto;

import com.fasterxml.jackson.annotation.JsonInclude;

import java.time.ZonedDateTime;

@JsonInclude(JsonInclude.Include.NON_NULL)
public record APIResponse<T>(
        T data,
        String message,
        boolean success,
        ZonedDateTime timestamp
) {
    public APIResponse(T data, String message, boolean success) {
        this(data, message, success, ZonedDateTime.now());
    }

    public static <T> APIResponse<T> success(T data) {
        return new APIResponse<>(data, "success", true);
    }

    public static <T> APIResponse<T> success(T data, String message) {
        return new APIResponse<>(data, message, true);
    }

    public static <T> APIResponse<T> error(String message) {
        return new APIResponse<>(null, message, false);
    }
}