package com.tourism.stakeholders.controller;

import com.tourism.stakeholders.dto.UpdateProfileDTO;
import com.tourism.stakeholders.model.User;
import com.tourism.stakeholders.service.UserService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;
import com.tourism.stakeholders.dto.UserResponseDTO;
import org.springframework.web.server.ResponseStatusException;
import com.tourism.stakeholders.dto.UserProfileDTO;

import java.util.List;

@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
public class UserController {

    private final UserService userService;

    @GetMapping
    public ResponseEntity<List<UserResponseDTO>> getAllUsers() {
        List<UserResponseDTO> users = userService.getAllUsers()
                .stream()
                .map(UserResponseDTO::fromUser)
                .toList();
        return ResponseEntity.ok(users);
    }

    @GetMapping("/{id}")
    public ResponseEntity<UserResponseDTO> getUserById(@PathVariable Long id) {
        User user = userService.getById(id);
        if (user == null) {
            throw new ResponseStatusException(HttpStatus.NOT_FOUND, "User not found");
        }
        return ResponseEntity.ok(UserResponseDTO.fromUser(user));
    }
    
    @PutMapping("/update-profile")
    public ResponseEntity<UserResponseDTO> updateProfile(@RequestBody UpdateProfileDTO request,
                                                         Authentication authentication) {
        if (authentication == null || authentication.getName() == null) {
            throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "Unauthorized");
        }

        Long userId;
        try {
            userId = Long.parseLong(authentication.getName());
        } catch (NumberFormatException ex) {
            throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "Invalid token subject");
        }

        User updated = userService.updateProfile(userId, request);
        return ResponseEntity.ok(UserResponseDTO.fromUser(updated));
    }

    @GetMapping("/profile")
    public ResponseEntity<UserProfileDTO> getMyProfile(Authentication authentication) {

        if (authentication == null || authentication.getName() == null) {
            throw new ResponseStatusException(
                    HttpStatus.UNAUTHORIZED,
                    "Unauthorized"
            );
        }

        Long userId;

        try {
            userId = Long.parseLong(authentication.getName());
        } catch (NumberFormatException ex) {
            throw new ResponseStatusException(
                    HttpStatus.UNAUTHORIZED,
                    "Invalid token subject"
            );
        }

        User user = userService.getById(userId);

        return ResponseEntity.ok(
                UserProfileDTO.fromUser(user)
        );
    }

    @PutMapping("/{id}/toggle-block")
    public ResponseEntity<UserResponseDTO> toggleBlock(@PathVariable Long id) {

        User updated = userService.toggleBlockUser(id);

        return ResponseEntity.ok(UserResponseDTO.fromUser(updated));
    }

    @PostMapping("/profile-image")
    public ResponseEntity<UserProfileDTO> uploadProfileImage(@RequestParam("image") org.springframework.web.multipart.MultipartFile file,
                                                             Authentication authentication) {
        if (authentication == null || authentication.getName() == null) {
            throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "Unauthorized");
        }

        Long userId;
        try {
            userId = Long.parseLong(authentication.getName());
        } catch (NumberFormatException ex) {
            throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "Invalid token subject");
        }

        if (file.isEmpty()) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "File is empty");
        }

        User updatedUser = userService.uploadProfileImage(userId, file);

        return ResponseEntity.ok(UserProfileDTO.fromUser(updatedUser));
    }

}