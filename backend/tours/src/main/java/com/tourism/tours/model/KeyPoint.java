package com.tourism.tours.model;

import lombok.Data;

@Data
public class KeyPoint {

    private String name;
    private String description;
    private double lat;
    private double lon;
    private String imageUrl;
}