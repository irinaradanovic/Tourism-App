package com.tourism.tours.model;

import lombok.Data;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.List;

@Data
public class Review {

    private String touristId;
    private String touristUsername;
    private int rating;
    private String comment;
    private LocalDate visitedDate;
    private LocalDateTime createdAt;
    private List<String> images;
}