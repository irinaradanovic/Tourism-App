package com.tourism.tours.repository;

import com.tourism.tours.model.TouristPosition;
import org.springframework.data.mongodb.repository.MongoRepository;

import java.util.Optional;

public interface TouristPositionRepository extends MongoRepository<TouristPosition, String> {

    Optional<TouristPosition> findByTouristId(Long touristId);
}