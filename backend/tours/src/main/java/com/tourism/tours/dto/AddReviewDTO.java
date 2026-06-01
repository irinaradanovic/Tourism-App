package com.tourism.tours.dto;

import lombok.Data;
import org.springframework.cglib.core.Local;

import java.time.LocalDate;
import java.util.List;

@Data
public class AddReviewDTO {
    private int rating;
    private String comment;
    private LocalDate visitedDate;
}