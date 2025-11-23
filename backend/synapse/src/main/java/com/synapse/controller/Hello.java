package com.synapse.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.HashMap;
import java.util.Map;

@RestController
public class Hello {
    @GetMapping("/")
    public Map<String,String> hello(){
        return new HashMap<>(){{put("response","hello world");}};
    }
}
