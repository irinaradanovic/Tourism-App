package com.tourism.tours.controller;

import com.tourism.tours.config.JwtUtil;
import com.tourism.tours.dto.AddReviewDTO;
import com.tourism.tours.dto.CreateKeyPointDTO;
import com.tourism.tours.dto.CreateTourDTO;
import com.tourism.tours.model.Tour;
import com.tourism.tours.model.TourDuration;
import com.tourism.tours.model.TourStatus;
import com.tourism.tours.service.TourService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestPart;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

@Slf4j
@RestController
@RequestMapping("/api/tours")
@RequiredArgsConstructor
@CrossOrigin(origins = "*")
public class TourController {

    private final TourService tourService;
    private final JwtUtil jwtUtil;
    private final RestTemplate restTemplate;

    @Value("${purchase.service.url:http://purchase:8084}")
    private String purchaseServiceUrl;

    @PostMapping
    public ResponseEntity<Tour> createTour(
            @RequestBody CreateTourDTO dto,
            @RequestHeader("Authorization") String authHeader
    ) {
        String token = authHeader.substring(7);
        Long userId = jwtUtil.extractUserId(token);

        log.info("[INFO] Create tour request received by userId={}", userId);
        Tour created = tourService.createTour(dto, userId);
        log.info("[INFO] Draft tour created successfully tourId={}, authorId={}", created.getId(), userId);

        return ResponseEntity.ok(created);
    }

    @GetMapping("/my")
    public ResponseEntity<List<Tour>> getMyTours(@RequestHeader("Authorization") String authHeader) {
        String token = authHeader.substring(7);
        Long userId = jwtUtil.extractUserId(token);
        return ResponseEntity.ok(tourService.getMyTours(userId));
    }

    @GetMapping("/{id}")
    public ResponseEntity<Tour> getTourById(
            @PathVariable String id,
            @RequestHeader(value = "Authorization", required = false) String authHeader
    ) {
        String role = null;
        Long userId = null;

        if (authHeader != null && authHeader.startsWith("Bearer ")) {
            try {
                String token = authHeader.substring(7);
                role = jwtUtil.extractRole(token);
                userId = jwtUtil.extractUserId(token);
            } catch (Exception ignored) {
            }
        }

        if ("GUIDE".equals(role) || "ADMIN".equals(role)) {
            return ResponseEntity.ok(tourService.getTourById(id));
        }

        if ("TOURIST".equals(role) && userId != null && hasPurchasedTour(id, userId)) {
            return ResponseEntity.ok(tourService.getTourById(id));
        }

        return ResponseEntity.ok(tourService.getTourPreview(id));
    }

    @PostMapping("/{tourId}/key-points")
    public ResponseEntity<Tour> addKeyPoint(
            @PathVariable String tourId,
            @RequestPart("data") CreateKeyPointDTO dto,
            @RequestPart(value = "image", required = false) MultipartFile image,
            @RequestHeader("Authorization") String authHeader
    ) throws IOException {
        String token = authHeader.substring(7);
        Long userId = jwtUtil.extractUserId(token);
        String role = jwtUtil.extractRole(token);
        return ResponseEntity.ok(tourService.addKeyPoint(tourId, dto, image, userId, role));
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
        return ResponseEntity.ok(tourService.updateKeyPoint(tourId, index, dto, image, userId, role));
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
        return ResponseEntity.ok(tourService.deleteKeyPoint(tourId, index, userId, role));
    }

    @PutMapping("/{tourId}/durations")
    public ResponseEntity<Tour> updateDurations(
            @PathVariable String tourId,
            @RequestBody List<TourDuration> durations,
            @RequestHeader("Authorization") String authHeader
    ) {
        String token = authHeader.substring(7);
        Long userId = jwtUtil.extractUserId(token);
        String role = jwtUtil.extractRole(token);
        return ResponseEntity.ok(tourService.updateDurations(tourId, durations, userId, role));
    }

    @PostMapping("/{tourId}/publish")
    public ResponseEntity<Map<String, String>> publishTour(
            @PathVariable String tourId,
            @RequestHeader("Authorization") String authHeader
    ) {
        String token = authHeader.substring(7);
        Long userId = jwtUtil.extractUserId(token);
        String role = jwtUtil.extractRole(token);
        tourService.requestPublish(tourId, userId, role);
        return ResponseEntity.accepted().body(Map.of("message", "Publish saga started"));
    }

    @PostMapping("/{tourId}/archive")
    public ResponseEntity<Map<String, String>> archiveTour(
            @PathVariable String tourId,
            @RequestHeader("Authorization") String authHeader
    ) {
        String token = authHeader.substring(7);
        Long userId = jwtUtil.extractUserId(token);
        String role = jwtUtil.extractRole(token);
        tourService.requestArchive(tourId, userId, role);
        return ResponseEntity.accepted().body(Map.of("message", "Archive saga started"));
    }

    @PostMapping("/{tourId}/reactivate")
    public ResponseEntity<Tour> reactivateTour(
            @PathVariable String tourId,
            @RequestHeader("Authorization") String authHeader
    ) {
        String token = authHeader.substring(7);
        Long userId = jwtUtil.extractUserId(token);
        String role = jwtUtil.extractRole(token);
        return ResponseEntity.ok(tourService.reactivateTour(tourId, userId, role));
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
        return ResponseEntity.ok(tourService.addReview(tourId, dto, userId, images));
    }

    @GetMapping("/published")
    public ResponseEntity<List<Tour>> getPublishedTours(
            @RequestHeader(value = "Authorization", required = false) String authHeader
    ) {
        String role = null;
        Long userId = null;

        if (authHeader != null && authHeader.startsWith("Bearer ")) {
            try {
                String token = authHeader.substring(7);
                role = jwtUtil.extractRole(token);
                userId = jwtUtil.extractUserId(token);
            } catch (Exception ignored) {
            }
        }

        List<Tour> tours = tourService.getToursByStatus(TourStatus.PUBLISHED);

        if ("GUIDE".equals(role) || "ADMIN".equals(role)) {
            return ResponseEntity.ok(tours);
        }

        List<Tour> result = new ArrayList<>();
        for (Tour tour : tours) {
            if (userId != null && hasPurchasedTour(tour.getId(), userId)) {
                result.add(tourService.getTourById(tour.getId()));
            } else {
                result.add(tourService.getTourPreview(tour.getId()));
            }
        }

        return ResponseEntity.ok(result);
    }

    private boolean hasPurchasedTour(String tourId, Long userId) {
        try {
            String url = purchaseServiceUrl + "/api/purchase/check/" + tourId + "?touristId=" + userId;
            Map response = restTemplate.getForObject(url, Map.class);
            return response != null && Boolean.TRUE.equals(response.get("purchased"));
        } catch (Exception ignored) {
            return false;
        }
    }
}
