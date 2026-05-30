package com.tourism.tours.controller;

import com.tourism.tours.config.JwtUtil;
import com.tourism.tours.dto.StartTourDTO;
import com.tourism.tours.model.TourExecution;
import com.tourism.tours.service.TourExecutionService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/executions")
@RequiredArgsConstructor
@CrossOrigin(origins = "*")
public class TourExecutionController {

    private final TourExecutionService executionService;
    private final JwtUtil jwtUtil;

    @PostMapping("/start/{tourId}")
    public ResponseEntity<TourExecution> startTour(
            @PathVariable String tourId,
            @RequestBody StartTourDTO dto,
            @RequestHeader("Authorization") String authHeader
    ) {
        Long userId = jwtUtil.extractUserId(authHeader.substring(7));
        TourExecution execution = executionService.startTour(tourId, userId, dto);
        return ResponseEntity.ok(execution);
    }

    @PutMapping("/{executionId}/abandon")
    public ResponseEntity<TourExecution> abandonTour(
            @PathVariable String executionId,
            @RequestHeader("Authorization") String authHeader
    ) {
        Long userId = jwtUtil.extractUserId(authHeader.substring(7));
        return ResponseEntity.ok(executionService.abandonTour(executionId, userId));
    }

    @PutMapping("/{executionId}/complete")
    public ResponseEntity<TourExecution> completeTour(
            @PathVariable String executionId,
            @RequestHeader("Authorization") String authHeader
    ) {
        Long userId = jwtUtil.extractUserId(authHeader.substring(7));
        return ResponseEntity.ok(executionService.completeTour(executionId, userId));
    }

    @PostMapping("/{executionId}/proximity")
    public ResponseEntity<TourExecution> checkProximity(
            @PathVariable String executionId,
            @RequestBody Map<String, Double> body,
            @RequestHeader("Authorization") String authHeader
    ) {
        Long userId = jwtUtil.extractUserId(authHeader.substring(7));
        double lat = body.get("lat");
        double lon = body.get("lon");
        return ResponseEntity.ok(executionService.checkProximity(executionId, userId, lat, lon));
    }

    @GetMapping("/my")
    public ResponseEntity<List<TourExecution>> getMyExecutions(
            @RequestHeader("Authorization") String authHeader
    ) {
        Long userId = jwtUtil.extractUserId(authHeader.substring(7));
        return ResponseEntity.ok(executionService.getMyExecutions(userId));
    }
}