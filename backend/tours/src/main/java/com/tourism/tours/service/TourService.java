package com.tourism.tours.service;

import com.tourism.tours.dto.AddReviewDTO;
import com.tourism.tours.dto.CreateTourDTO;
import com.tourism.tours.dto.ReviewDTO;
import com.tourism.tours.model.Review;
import com.tourism.tours.model.Tour;
import com.tourism.tours.repository.TourRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;
import org.springframework.web.multipart.MultipartFile;

import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.io.File;
import java.io.IOException;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class TourService {

    private final TourRepository tourRepository;

    public Tour createTour(CreateTourDTO dto, Long authorId) {

        Tour tour = new Tour();

        tour.setAuthorId(authorId);
        tour.setTitle(dto.getTitle());
        tour.setDescription(dto.getDescription());
        tour.setDifficulty(dto.getDifficulty());
        tour.setTags(dto.getTags());

        return tourRepository.save(tour);
    }

    public List<Tour> getMyTours(Long authorId) {
        return tourRepository.findByAuthorId(authorId);
    }

    public Tour getTourById(String id) {
        return tourRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Tour not found"));
    }

    public Tour save(Tour tour) {
        return tourRepository.save(tour);
    }

    public Tour addReview(
            String tourId,
            AddReviewDTO dto,
            Long userId,
            List<MultipartFile> images
    ) throws IOException {

        if (dto.getRating() < 1 || dto.getRating() > 5) {
            throw new IllegalArgumentException("Rating must be between 1 and 5");
        }
        if (dto.getComment() == null || dto.getComment().isBlank()) {
            throw new IllegalArgumentException("Comment is required");
        }
        if (dto.getVisitedDate() == null) {
            throw new IllegalArgumentException("Visited date is required");
        }

        Tour tour = getTourById(tourId);

        if (tour.getReviews() == null) {
            tour.setReviews(new ArrayList<>());
        }

        Review review = new Review();
        review.setTouristId(String.valueOf(userId));
        review.setRating(dto.getRating());
        review.setComment(dto.getComment());
        review.setVisitedDate(dto.getVisitedDate());
        review.setCreatedAt(LocalDateTime.now());


        List<String> storedImages = new ArrayList<>();
        if (images != null && !images.isEmpty()) {
        Path uploadDir = Paths.get("uploads", "blogs");
        Files.createDirectories(uploadDir);

        for (MultipartFile file : images) {
            if (file != null && !file.isEmpty()) {
                String fileName = UUID.randomUUID() + "_" +
                        StringUtils.cleanPath(file.getOriginalFilename());
                Path destination = uploadDir.resolve(fileName);
                file.transferTo(destination);
                storedImages.add(fileName);
            }
        }
}
        review.setImages(storedImages);

        tour.getReviews().add(review);

        return save(tour);
    }
}