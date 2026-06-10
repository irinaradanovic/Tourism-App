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

    private Long authorId;

    private String title;

    private String description;

    private Difficulty difficulty;

    private List<String> tags = new ArrayList<>();

    private Double price = 0.0;

    private Double distanceKm = 0.0;

    private TourStatus status = TourStatus.DRAFT;

    private LocalDateTime publishedAt;

    private LocalDateTime archivedAt;

    private List<TourDuration> durations = new ArrayList<>();

    private List<KeyPoint> keyPoints = new ArrayList<>();

    private List<Review> reviews = new ArrayList<>();
}
