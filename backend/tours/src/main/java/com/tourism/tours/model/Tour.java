package com.tourism.tours.model;

import lombok.Data;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;

@Data
@Document(collection = "tours")
public class Tour {

    @Id
    private String id;

    private String authorId;
    private String name;
    private String description;
    private String difficulty;
    private List<String> tags = new ArrayList<>();

    private String status = "DRAFT";
    private double price = 0.0;

    private List<KeyPoint> keyPoints = new ArrayList<>();
    private List<Review> reviews = new ArrayList<>();

    private LocalDateTime createdAt;
}