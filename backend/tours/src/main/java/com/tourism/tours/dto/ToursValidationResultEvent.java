package com.tourism.tours.dto;

import lombok.Data;

import com.fasterxml.jackson.annotation.JsonProperty;

import lombok.AllArgsConstructor;


@Data
@AllArgsConstructor
public class ToursValidationResultEvent {
    @JsonProperty("cart_id")
    private Long cartId;

    @JsonProperty("tourist_id")
    private Long touristId;

    @JsonProperty("reason")
    private String reason;   
}
