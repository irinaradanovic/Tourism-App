package com.tourism.tours.dto;

import lombok.Data;

import java.util.List;

@Data
public class CreateTourDTO {
    private String name;
    private String description;
    private String difficulty;
    private List<String> tags;
}