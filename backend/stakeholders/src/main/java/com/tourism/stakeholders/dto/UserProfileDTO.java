package com.tourism.stakeholders.dto;

import com.tourism.stakeholders.model.User;
import lombok.AllArgsConstructor;
import lombok.Data;

@Data
@AllArgsConstructor
public class UserProfileDTO {

    private Long id;

    private String username;

    private String email;

    private String firstName;

    private String lastName;

    private String biography;

    private String motto;

    private String profileImage;

    private String role;

    public static UserProfileDTO fromUser(User user) {

        return new UserProfileDTO(
                user.getId(),
                user.getUsername(),
                user.getEmail(),
                user.getFirstName(),
                user.getLastName(),
                user.getBiography(),
                user.getMotto(),
                user.getProfileImage(),
                user.getRole().name()
        );
    }
}