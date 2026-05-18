package com.tourism.tours.controller;

import com.tourism.tours.dto.CreateKeyPointDTO;
import com.tourism.tours.dto.CreateTourDTO;
import com.tourism.tours.model.KeyPoint;
import com.tourism.tours.model.Tour;
import com.tourism.tours.config.JwtUtil;
import com.tourism.tours.service.TourService;
import com.tourism.tours.dto.ReviewDTO;
import com.tourism.tours.model.Review;
import com.tourism.tours.dto.AddReviewDTO;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.util.StringUtils;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.io.File;
import java.io.IOException;
import java.time.LocalDateTime;
import java.util.ArrayList;
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

    @GetMapping("/{id}")
    public ResponseEntity<Tour> getTourById(@PathVariable String id) {
                System.out.println("GETTOUR hit");
        return ResponseEntity.ok(tourService.getTourById(id));
    }

    @PostMapping("/{tourId}/key-points")
    public ResponseEntity<Tour> addKeyPoint(
            @PathVariable String tourId,
            @RequestPart("data") CreateKeyPointDTO dto,
            @RequestPart(value = "image", required = false) MultipartFile image // required = false sprečava pucanje ako nema slike
    ) throws IOException {

        Tour tour = tourService.getTourById(tourId);

        // 1. BEZBEDNOSNA PROVERA: Inicijalizuj listu ako je null da izbegneš NullPointerException
        if (tour.getKeyPoints() == null) {
            tour.setKeyPoints(new java.util.ArrayList<>());
        }

        KeyPoint keyPoint = new KeyPoint();
        keyPoint.setName(dto.getName());
        keyPoint.setDescription(dto.getDescription());
        keyPoint.setLatitude(dto.getLatitude());
        keyPoint.setLongitude(dto.getLongitude());

        // 2. Obrada slike samo ako je korisnik zapravo izabrao fajl
        if (image != null && !image.isEmpty()) {
            String fileName = UUID.randomUUID() + "_" +
                    org.springframework.util.StringUtils.cleanPath(image.getOriginalFilename());

            File uploadDir = new File("uploads");
            if (!uploadDir.exists()) {
                uploadDir.mkdirs();
            }

            File destination = new File(uploadDir, fileName);
            image.transferTo(destination);

            keyPoint.setImage(fileName);
        } else {
            keyPoint.setImage(null); // Ili neku default sliku ako želiš
        }

        tour.getKeyPoints().add(keyPoint);

        return ResponseEntity.ok(
                tourService.save(tour)
        );
    }
    @PostMapping("/{tourId}/reviews")
    public ResponseEntity<Tour> addReview(
            @PathVariable String tourId,
            @RequestPart("data") AddReviewDTO dto,
            @RequestPart(value = "images", required = false) List<MultipartFile> images,
            @RequestHeader("Authorization") String authHeader
    ) throws IOException {

        String token = authHeader.substring(7);
        Long userId = jwtUtil.extractUserId(token);
        System.out.println("addReview hit");

        Tour updated = tourService.addReview(tourId, dto, userId, images);
        //Tour updated = tourService.addReview(tourId, dto, userId);
        return ResponseEntity.ok(updated);
    }
    @PutMapping("/{tourId}/key-points/{index}")
    public ResponseEntity<Tour> updateKeyPoint(
            @PathVariable String tourId,
            @PathVariable int index,
            @RequestPart("data") CreateKeyPointDTO dto,
            @RequestPart(value = "image", required = false) MultipartFile image,
            @RequestHeader("Authorization") String authHeader
    ) throws IOException {

        String token = authHeader.substring(7);
        Long userId = jwtUtil.extractUserId(token);
        String role = jwtUtil.extractRole(token);

        Tour updated = tourService.updateKeyPoint(tourId, index, dto, image, userId, role);
        return ResponseEntity.ok(updated);
    }
    @DeleteMapping("/{tourId}/key-points/{index}")
    public ResponseEntity<Tour> deleteKeyPoint(
            @PathVariable String tourId,
            @PathVariable int index,
            @RequestHeader("Authorization") String authHeader
    ) {
        String token = authHeader.substring(7);
        Long userId = jwtUtil.extractUserId(token);
        String role = jwtUtil.extractRole(token);

        Tour deleted = tourService.deleteKeyPoint(tourId,index,userId,role);
        return ResponseEntity.ok(deleted);
    }
}