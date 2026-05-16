package com.tourism.tours.model;

import lombok.Data;

@Data
public class KeyPoint {

    private String name;

    private String description;

    private Double latitude;

    private Double longitude;

    private String image;
}