package com.tourism.tours.service;

import com.tourism.tours.dto.StartTourDTO;
import com.tourism.tours.model.KeyPoint;
import com.tourism.tours.model.Tour;
import com.tourism.tours.model.TourExecution;
import com.tourism.tours.repository.TourExecutionRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import java.time.LocalDateTime;
import java.util.List;

@Service
@RequiredArgsConstructor
public class TourExecutionService {

    private final TourExecutionRepository executionRepository;
    private final TourService tourService;

    public TourExecution startTour(String tourId, Long touristId, StartTourDTO dto) {
        Tour tour = tourService.getTourById(tourId);

        if (!"PUBLISHED".equals(tour.getStatus().toString()) && !"ARCHIVED".equals(tour.getStatus().toString())) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Tour must be published or archived to start");
        }

        TourExecution execution = new TourExecution();
        execution.setTourId(tourId);
        execution.setTouristId(touristId);
        execution.setStatus("ACTIVE");
        execution.setStartLat(dto.getLat());
        execution.setStartLon(dto.getLon());
        execution.setStartedAt(LocalDateTime.now());
        execution.setLastActivity(LocalDateTime.now());

        return executionRepository.save(execution);
    }

    public TourExecution abandonTour(String executionId, Long touristId) {
        TourExecution execution = executionRepository.findById(executionId)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Execution not found"));

        if (!execution.getTouristId().equals(touristId)) {
            throw new ResponseStatusException(HttpStatus.FORBIDDEN, "Not your tour execution");
        }

        execution.setStatus("ABANDONED");
        execution.setFinishedAt(LocalDateTime.now());
        execution.setLastActivity(LocalDateTime.now());

        return executionRepository.save(execution);
    }

    public TourExecution completeTour(String executionId, Long touristId) {
        TourExecution execution = executionRepository.findById(executionId)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Execution not found"));

        if (!execution.getTouristId().equals(touristId)) {
            throw new ResponseStatusException(HttpStatus.FORBIDDEN, "Not your tour execution");
        }

        execution.setStatus("COMPLETED");
        execution.setFinishedAt(LocalDateTime.now());
        execution.setLastActivity(LocalDateTime.now());

        return executionRepository.save(execution);
    }

    public TourExecution checkProximity(String executionId, Long touristId, double lat, double lon) {
        TourExecution execution = executionRepository.findById(executionId)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Execution not found"));

        if (!execution.getTouristId().equals(touristId)) {
            throw new ResponseStatusException(HttpStatus.FORBIDDEN, "Not your tour execution");
        }

        execution.setLastActivity(LocalDateTime.now());

        Tour tour = tourService.getTourById(execution.getTourId());
        List<KeyPoint> keyPoints = tour.getKeyPoints();

        for (int i = 0; i < keyPoints.size(); i++) {
            if (execution.getCompletedKeyPoints().containsKey(i)) {
                continue;
            }

            KeyPoint kp = keyPoints.get(i);
            double distance = calculateDistance(lat, lon, kp.getLatitude(), kp.getLongitude());

            if (distance <= 0.2) {
                execution.getCompletedKeyPoints().put(i, LocalDateTime.now());
            }
        }

        return executionRepository.save(execution);
    }

    public List<TourExecution> getMyExecutions(Long touristId) {
        return executionRepository.findByTouristId(touristId);
    }

    private double calculateDistance(double lat1, double lon1, double lat2, double lon2) {
        final int R = 6371;
        double dLat = Math.toRadians(lat2 - lat1);
        double dLon = Math.toRadians(lon2 - lon1);
        double a = Math.sin(dLat / 2) * Math.sin(dLat / 2)
                + Math.cos(Math.toRadians(lat1)) * Math.cos(Math.toRadians(lat2))
                * Math.sin(dLon / 2) * Math.sin(dLon / 2);
        double c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
        return R * c;
    }
}