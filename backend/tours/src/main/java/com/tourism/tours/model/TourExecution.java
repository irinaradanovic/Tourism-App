package com.tourism.tours.model;

import lombok.Data;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.Map;

@Data
@Document(collection = "tour_executions")
public class TourExecution {

    @Id
    private String id;

    private String tourId;
    private Long touristId;

    private String status; // ACTIVE, COMPLETED, ABANDONED

    private double startLat;
    private double startLon;

    private LocalDateTime startedAt;
    private LocalDateTime finishedAt;
    private LocalDateTime lastActivity;

    private Map<Integer, LocalDateTime> completedKeyPoints = new HashMap<>();
}