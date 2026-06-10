package com.tourism.tours.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class TourLifecycleRequestedEvent {
    @JsonProperty("tour_id")
    private String tourId;

    @JsonProperty("author_id")
    private Long authorId;

    @JsonProperty("action")
    private String action;
}
