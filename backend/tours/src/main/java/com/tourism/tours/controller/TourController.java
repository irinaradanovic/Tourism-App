package com.tourism.tours.controller;

import com.tourism.tours.dto.CreateTourDTO;
import com.tourism.tours.dto.ReviewDTO;
import com.tourism.tours.model.KeyPoint;
import com.tourism.tours.model.Tour;
import com.tourism.tours.service.TourService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/tours")
public class TourController {

    private final TourService tourService;

    public TourController(TourService tourService) {
        this.tourService = tourService;
    }

    @PostMapping
    public ResponseEntity<Tour> createTour(@RequestBody CreateTourDTO dto, Authentication auth) {
        Tour created = tourService.createTour(dto, auth.getName());
        return ResponseEntity.status(HttpStatus.CREATED).body(created);
    }

    @GetMapping("/my")
    public ResponseEntity<List<Tour>> getMyTours(Authentication auth) {
        return ResponseEntity.ok(tourService.getMyTours(auth.getName()));
    }

    @GetMapping
    public ResponseEntity<List<Tour>> getPublishedTours() {
        return ResponseEntity.ok(tourService.getPublishedTours());
    }

    @GetMapping("/{id}")
    public ResponseEntity<Tour> getById(@PathVariable String id) {
        return tourService.getById(id)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }

    @PostMapping("/{id}/keypoints")
    public ResponseEntity<?> addKeyPoint(
            @PathVariable String id,
            @RequestBody KeyPoint keyPoint,
            Authentication auth) {
        try {
            Tour updated = tourService.addKeyPoint(id, keyPoint, auth.getName());
            return ResponseEntity.ok(updated);
        } catch (RuntimeException e) {
            return ResponseEntity.status(HttpStatus.FORBIDDEN).body(Map.of("error", e.getMessage()));
        }
    }

    @PutMapping("/{id}/keypoints/{index}")
    public ResponseEntity<?> updateKeyPoint(
            @PathVariable String id,
            @PathVariable int index,
            @RequestBody KeyPoint keyPoint,
            Authentication auth) {
        try {
            Tour updated = tourService.updateKeyPoint(id, index, keyPoint, auth.getName());
            return ResponseEntity.ok(updated);
        } catch (RuntimeException e) {
            return ResponseEntity.status(HttpStatus.FORBIDDEN).body(Map.of("error", e.getMessage()));
        }
    }

    @DeleteMapping("/{id}/keypoints/{index}")
    public ResponseEntity<?> deleteKeyPoint(
            @PathVariable String id,
            @PathVariable int index,
            Authentication auth) {
        try {
            Tour updated = tourService.deleteKeyPoint(id, index, auth.getName());
            return ResponseEntity.ok(updated);
        } catch (RuntimeException e) {
            return ResponseEntity.status(HttpStatus.FORBIDDEN).body(Map.of("error", e.getMessage()));
        }
    }

    @PostMapping("/{id}/reviews")
    public ResponseEntity<?> addReview(
            @PathVariable String id,
            @RequestBody ReviewDTO dto,
            Authentication auth) {
        try {
            Tour updated = tourService.addReview(id, dto, auth.getName());
            return ResponseEntity.ok(updated);
        } catch (RuntimeException e) {
            return ResponseEntity.badRequest().body(Map.of("error", e.getMessage()));
        }
    }
}