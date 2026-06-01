package com.tourism.tours.repository;

import com.tourism.tours.model.TourExecution;
import org.springframework.data.mongodb.repository.MongoRepository;

import java.util.List;
import java.util.Optional;

public interface TourExecutionRepository extends MongoRepository<TourExecution, String> {
    Optional<TourExecution> findByTourIdAndTouristIdAndStatus(String tourId, Long touristId, String status);
    List<TourExecution> findByTouristId(Long touristId);
}