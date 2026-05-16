package com.tourism.tours.controller;

import com.tourism.tours.model.TouristPosition;
import com.tourism.tours.service.PositionService;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RestController
@RequestMapping("/api/position")
public class PositionController {

    private final PositionService positionService;

    public PositionController(PositionService positionService) {
        this.positionService = positionService;
    }

    @PostMapping
    public ResponseEntity<TouristPosition> savePosition(
            Authentication auth,
            @RequestBody Map<String, Double> body) {

        String touristId = auth.getName();
        double lat = body.get("lat");
        double lon = body.get("lon");

        TouristPosition saved = positionService.savePosition(touristId, lat, lon);
        return ResponseEntity.ok(saved);
    }

    @GetMapping
    public ResponseEntity<TouristPosition> getPosition(Authentication auth) {
        String touristId = auth.getName();

        return positionService.getPosition(touristId)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
}