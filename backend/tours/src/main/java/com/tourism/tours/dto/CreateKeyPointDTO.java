package com.tourism.tours.dto;

import lombok.Data;

@Data
public class CreateKeyPointDTO {

    private String name;

    private String description;

    private Double latitude;

    private Double longitude;
}