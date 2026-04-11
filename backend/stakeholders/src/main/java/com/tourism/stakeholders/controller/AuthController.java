package com.tourism.stakeholders.controller;

import com.tourism.stakeholders.model.User;
import com.tourism.stakeholders.service.UserService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import com.tourism.stakeholders.dto.LoginRequestDTO;
import com.tourism.stakeholders.dto.UserResponseDTO;

@RestController
@RequestMapping("/api/auth")
@RequiredArgsConstructor
public class AuthController {

    private final UserService userService;

    @PostMapping("/register")
    public ResponseEntity<User> register(@RequestBody User user) {
        User saved = userService.register(user);
        return ResponseEntity.ok(saved);
    }

    @PostMapping("/login")
    public ResponseEntity<UserResponseDTO> login(@RequestBody LoginRequestDTO request) {
        User user = userService.login(request.getUsername(), request.getPassword());
        return ResponseEntity.ok(UserResponseDTO.fromUser(user));
    }
}