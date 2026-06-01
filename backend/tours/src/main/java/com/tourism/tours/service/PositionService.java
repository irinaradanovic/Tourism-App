package com.tourism.tours.service;

import com.tourism.tours.model.TouristPosition;
import com.tourism.tours.repository.TouristPositionRepository;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.Optional;

@Service
public class PositionService {

    private final TouristPositionRepository positionRepository;

    public PositionService(TouristPositionRepository positionRepository) {
        this.positionRepository = positionRepository;
    }

    public TouristPosition savePosition(Long touristId, double lat, double lon) {
        TouristPosition position = positionRepository
                .findByTouristId(touristId)
                .orElse(new TouristPosition());

        position.setTouristId(touristId);
        position.setLat(lat);
        position.setLon(lon);
        position.setUpdatedAt(LocalDateTime.now());

        return positionRepository.save(position);
    }

    public Optional<TouristPosition> getPosition(Long touristId) {
        return positionRepository.findByTouristId(touristId);
    }
}