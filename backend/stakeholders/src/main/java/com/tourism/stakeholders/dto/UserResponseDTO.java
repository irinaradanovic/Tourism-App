package com.tourism.stakeholders.dto;

import com.tourism.stakeholders.model.User;
import lombok.Data;

@Data
public class UserResponseDTO {

    private Long id;
    private String username;
    private String email;
    private User.Role role;
    private boolean blocked;

    public static UserResponseDTO fromUser(User user) {
        UserResponseDTO dto = new UserResponseDTO();
        dto.setId(user.getId());
        dto.setUsername(user.getUsername());
        dto.setEmail(user.getEmail());
        dto.setRole(user.getRole());
        dto.setBlocked(user.isBlocked());
        return dto;
    }
}