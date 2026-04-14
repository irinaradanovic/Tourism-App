package com.tourism.stakeholders.dto;

import lombok.Data;
import lombok.Getter;

@Data
public class UpdateProfileDTO {
    private String username;
    private String email;
    private String password;
}