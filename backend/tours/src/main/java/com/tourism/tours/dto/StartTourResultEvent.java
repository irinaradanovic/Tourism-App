package com.tourism.tours.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class StartTourResultEvent {
    private String correlationId;
    private Boolean purchased;
    private String reason;
}