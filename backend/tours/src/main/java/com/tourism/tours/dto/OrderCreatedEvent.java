package com.tourism.tours.dto;
import lombok.Data;
import java.util.List;

import com.fasterxml.jackson.annotation.JsonProperty;

@Data
public class OrderCreatedEvent {
    @JsonProperty("cart_id")
    private Long cartId;

    @JsonProperty("tourist_id")
    private Long touristId;

    @JsonProperty("tour_ids")
    private List<String> tourIds;  
}
