package com.tourism.tours.repository;

import com.tourism.tours.model.Tour;
import org.springframework.data.mongodb.repository.MongoRepository;

import java.util.List;

public interface TourRepository extends MongoRepository<Tour, String> {

    List<Tour> findByAuthorId(String authorId);

    List<Tour> findByStatus(String status);
}