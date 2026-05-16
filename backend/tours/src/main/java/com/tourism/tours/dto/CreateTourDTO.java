package com.tourism.tours.dto;

import com.tourism.tours.model.Difficulty;
import lombok.Data;

import java.util.List;

@Data
public class CreateTourDTO {

    private String title;

    private String description;

    private Difficulty difficulty;

    private List<String> tags;
}