package com.tourism.tours.service;

import com.tourism.tours.dto.CreateTourDTO;
import com.tourism.tours.dto.ReviewDTO;
import com.tourism.tours.model.KeyPoint;
import com.tourism.tours.model.Review;
import com.tourism.tours.model.Tour;
import com.tourism.tours.repository.TourRepository;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;

@Service
public class TourService {

    private final TourRepository tourRepository;

    public TourService(TourRepository tourRepository) {
        this.tourRepository = tourRepository;
    }

    public Tour createTour(CreateTourDTO dto, String authorId) {
        Tour tour = new Tour();
        tour.setAuthorId(authorId);
        tour.setName(dto.getName());
        tour.setDescription(dto.getDescription());
        tour.setDifficulty(dto.getDifficulty());
        tour.setTags(dto.getTags());
        tour.setStatus("DRAFT");
        tour.setPrice(0.0);
        tour.setCreatedAt(LocalDateTime.now());
        return tourRepository.save(tour);
    }

    public List<Tour> getMyTours(String authorId) {
        return tourRepository.findByAuthorId(authorId);
    }

    public Tour addKeyPoint(String tourId, KeyPoint keyPoint, String authorId) {
        Tour tour = tourRepository.findById(tourId)
                .orElseThrow(() -> new RuntimeException("Tour not found"));

        if (!tour.getAuthorId().equals(authorId)) {
            throw new RuntimeException("Forbidden: you are not the author");
        }

        tour.getKeyPoints().add(keyPoint);
        return tourRepository.save(tour);
    }

    public Tour updateKeyPoint(String tourId, int index, KeyPoint updated, String authorId) {
        Tour tour = tourRepository.findById(tourId)
                .orElseThrow(() -> new RuntimeException("Tour not found"));

        if (!tour.getAuthorId().equals(authorId)) {
            throw new RuntimeException("Forbidden: you are not the author");
        }

        tour.getKeyPoints().set(index, updated);
        return tourRepository.save(tour);
    }

    public Tour deleteKeyPoint(String tourId, int index, String authorId) {
        Tour tour = tourRepository.findById(tourId)
                .orElseThrow(() -> new RuntimeException("Tour not found"));

        if (!tour.getAuthorId().equals(authorId)) {
            throw new RuntimeException("Forbidden: you are not the author");
        }

        tour.getKeyPoints().remove(index);
        return tourRepository.save(tour);
    }

    public Tour addReview(String tourId, ReviewDTO dto, String touristId) {
        Tour tour = tourRepository.findById(tourId)
                .orElseThrow(() -> new RuntimeException("Tour not found"));

        Review review = new Review();
        review.setTouristId(touristId);
        review.setTouristUsername(dto.getTouristUsername());
        review.setRating(dto.getRating());
        review.setComment(dto.getComment());
        review.setVisitedDate(dto.getVisitedDate());
        review.setCreatedAt(LocalDateTime.now());
        review.setImages(dto.getImages());

        tour.getReviews().add(review);
        return tourRepository.save(tour);
    }

    public Optional<Tour> getById(String tourId) {
        return tourRepository.findById(tourId);
    }

    public List<Tour> getPublishedTours() {
        return tourRepository.findByStatus("PUBLISHED");
    }
}