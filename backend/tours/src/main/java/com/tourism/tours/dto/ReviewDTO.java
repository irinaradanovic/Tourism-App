package com.tourism.tours.dto;

import lombok.Data;

import java.time.LocalDate;
import java.util.List;

@Data
public class ReviewDTO {
    private String touristUsername;
    private int rating;
    private String comment;
    private LocalDate visitedDate;
    private List<String> images;
}