package com.tourism.tours.model;

import lombok.Data;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.index.Indexed;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.LocalDateTime;

@Data
@Document(collection = "tourist_positions")
public class TouristPosition {

    @Id
    private String id;

    @Indexed(unique = true)
    private Long touristId;

    private double lat;
    private double lon;

    private LocalDateTime updatedAt;
}