package com.tourism.tours.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class StartTourRequestedEvent {
    private String correlationId;
    private String tourId;
    private Long touristId;
    private Double lat;
    private Double lon;
}