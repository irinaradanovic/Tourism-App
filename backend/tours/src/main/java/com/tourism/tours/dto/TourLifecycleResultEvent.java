package com.tourism.tours.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

@Data
public class TourLifecycleResultEvent {
    @JsonProperty("tour_id")
    private String tourId;

    @JsonProperty("action")
    private String action;

    @JsonProperty("success")
    private Boolean success;

    @JsonProperty("reason")
    private String reason;
}
