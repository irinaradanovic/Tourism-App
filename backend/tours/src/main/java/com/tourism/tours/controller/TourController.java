package com.tourism.tours.controller;

import com.tourism.tours.dto.CreateKeyPointDTO;
import com.tourism.tours.dto.CreateTourDTO;
import com.tourism.tours.model.KeyPoint;
import com.tourism.tours.model.Tour;
import com.tourism.tours.security.JwtUtil;
import com.tourism.tours.service.TourService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.util.StringUtils;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.io.File;
import java.io.IOException;
import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/api/tours")
@RequiredArgsConstructor
@CrossOrigin(origins = "*")
public class TourController {

    private final TourService tourService;
    private final JwtUtil jwtUtil;

    @PostMapping
    public ResponseEntity<Tour> createTour(
            @RequestBody CreateTourDTO dto,
            @RequestHeader("Authorization") String authHeader
    ) {

        String token = authHeader.substring(7);

        Long userId = jwtUtil.extractUserId(token);

        Tour created = tourService.createTour(dto, userId);

        return ResponseEntity.ok(created);
    }

    @GetMapping("/my")
    public ResponseEntity<List<Tour>> getMyTours(
            @RequestHeader("Authorization") String authHeader
    ) {

        String token = authHeader.substring(7);

        Long userId = jwtUtil.extractUserId(token);

        return ResponseEntity.ok(
                tourService.getMyTours(userId)
        );
    }

    @PostMapping("/{tourId}/key-points")
    public ResponseEntity<Tour> addKeyPoint(
            @PathVariable String tourId,
            @RequestPart("data") CreateKeyPointDTO dto,
            @RequestPart("image") MultipartFile image
    ) throws IOException {

        Tour tour = tourService.getTourById(tourId);

        String fileName = UUID.randomUUID() + "_" +
                StringUtils.cleanPath(image.getOriginalFilename());

        File uploadDir = new File("uploads");

        if (!uploadDir.exists()) {
            uploadDir.mkdirs();
        }

        File destination = new File(uploadDir, fileName);

        image.transferTo(destination);

        KeyPoint keyPoint = new KeyPoint();

        keyPoint.setName(dto.getName());
        keyPoint.setDescription(dto.getDescription());
        keyPoint.setLatitude(dto.getLatitude());
        keyPoint.setLongitude(dto.getLongitude());
        keyPoint.setImage(fileName);

        tour.getKeyPoints().add(keyPoint);

        return ResponseEntity.ok(
                tourService.save(tour)
        );
    }
}