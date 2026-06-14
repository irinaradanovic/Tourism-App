package com.tourism.tours.service;

import com.tourism.tours.dto.StartTourRequestedEvent;
import com.tourism.tours.dto.StartTourResultEvent;
import com.tourism.tours.dto.StartTourDTO;
import com.tourism.tours.model.KeyPoint;
import com.tourism.tours.model.Tour;
import com.tourism.tours.model.TourExecution;
import com.tourism.tours.repository.TourExecutionRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.UUID;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;

@Service
@RequiredArgsConstructor
@Slf4j
public class TourExecutionService {

    private final TourExecutionRepository executionRepository;
    private final TourService tourService;
    private final RabbitTemplate rabbitTemplate;

    // čuvamo pending SAGA zahtev
    private final Map<String, CompletableFuture<StartTourResultEvent>> pendingSagas =
            new ConcurrentHashMap<>();

    public TourExecution startTour(String tourId, Long touristId, StartTourDTO dto) {
        Tour tour = tourService.getTourById(tourId);

        if (!"PUBLISHED".equals(tour.getStatus().toString()) && !"ARCHIVED".equals(tour.getStatus().toString())) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Tour must be published or archived to start");
        }

        Optional<TourExecution> existingActive =
                executionRepository.findByTourIdAndTouristIdAndStatus(tourId, touristId, "ACTIVE");
        if (existingActive.isPresent()) {
            return existingActive.get();
        }

        // SAGA korak 1: šalji event ka Purchase servisu da proveri kupovinu
        String correlationId = UUID.randomUUID().toString();
        CompletableFuture<StartTourResultEvent> future = new CompletableFuture<>();
        pendingSagas.put(correlationId, future);

        StartTourRequestedEvent event = new StartTourRequestedEvent(
                correlationId, tourId, touristId, dto.getLat(), dto.getLon()
        );
        log.info("[SAGA-START-TOUR] Korak 1: Saljem proveru kupovine ka Purchase servisu. correlationId={}", correlationId);
        rabbitTemplate.convertAndSend("start.tour.requested", event);

        // čekaj odgovor max 5 sekundi
        StartTourResultEvent result;
        try {
            result = future.get(5, TimeUnit.SECONDS);
        } catch (Exception e) {
            pendingSagas.remove(correlationId);
            log.error("[SAGA-START-TOUR] Timeout - Purchase servis nije odgovorio na vreme: {}", e.getMessage());
            throw new ResponseStatusException(HttpStatus.SERVICE_UNAVAILABLE, "Purchase check timed out, try again");
        }

        // SAGA kompenzacija: ako tura nije kupljena, ne kreiramo sesiju
        if (!Boolean.TRUE.equals(result.getPurchased())) {
            log.warn("[SAGA-START-TOUR] SAGA FAILED - Tura nije kupljena. Razlog: {}", result.getReason());
            throw new ResponseStatusException(HttpStatus.FORBIDDEN, "You must purchase this tour before starting it");
        }

        // SAGA korak 2: kupovina potvrđena, kreiraj TourExecution
        log.info("[SAGA-START-TOUR] Korak 2: Kupovina potvrđena, kreiram TourExecution. correlationId={}", correlationId);
        TourExecution execution = new TourExecution();
        execution.setTourId(tourId);
        execution.setTouristId(touristId);
        execution.setStatus("ACTIVE");
        execution.setStartLat(dto.getLat());
        execution.setStartLon(dto.getLon());
        execution.setStartedAt(LocalDateTime.now());
        execution.setLastActivity(LocalDateTime.now());

        TourExecution saved = executionRepository.save(execution);
        log.info("[SAGA-START-TOUR] SAGA uspešno završena. executionId={}", saved.getId());
        return saved;
    }

    // poziva TourSagaConsumer kada Purchase servis odgovori
    public void handleStartTourResult(StartTourResultEvent event) {
        CompletableFuture<StartTourResultEvent> future = pendingSagas.remove(event.getCorrelationId());
        if (future != null) {
            future.complete(event);
        } else {
            log.warn("[SAGA-START-TOUR] Primljen odgovor za nepoznati correlationId={}", event.getCorrelationId());
        }
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
            if (distance <= 0.05) {
                execution.getCompletedKeyPoints().put(i, LocalDateTime.now());
            }
        }

        if (tour.getKeyPoints() != null && !tour.getKeyPoints().isEmpty()
                && execution.getCompletedKeyPoints().size() == tour.getKeyPoints().size()
                && !"COMPLETED".equals(execution.getStatus())) {
            execution.setStatus("COMPLETED");
            execution.setFinishedAt(LocalDateTime.now());
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