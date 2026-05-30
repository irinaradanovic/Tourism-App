package com.tourism.tours.service;

import com.tourism.tours.dto.AddReviewDTO;
import com.tourism.tours.dto.CreateTourDTO;
import com.tourism.tours.dto.ReviewDTO;
import com.tourism.tours.model.Review;
import com.tourism.tours.model.Tour;
import com.tourism.tours.repository.TourRepository;
import com.tourism.tours.dto.CreateKeyPointDTO;
import com.tourism.tours.model.KeyPoint;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.http.HttpStatus;
import org.springframework.web.server.ResponseStatusException;

import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.io.File;
import java.io.IOException;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.UUID;
import com.tourism.tours.model.TourStatus;

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
    private Tour checkTourOwner(
            String tourId,
            Long userId,
            String role,
            int index
    ){
        Tour tour = getTourById(tourId);

        if (!"GUIDE".equals(role) || userId == null || !userId.equals(tour.getAuthorId())) {
            throw new ResponseStatusException(HttpStatus.FORBIDDEN, "Only guide owner can update key points");
        }

        if (tour.getKeyPoints() == null || index < 0 || index >= tour.getKeyPoints().size()) {
            throw new ResponseStatusException(HttpStatus.NOT_FOUND, "Key point not found");
        }
        return tour;
    }
    public Tour updateKeyPoint(
            String tourId,
            int index,
            CreateKeyPointDTO dto,
            MultipartFile image,
            Long userId,
            String role
    ) throws IOException {

        Tour tour = checkTourOwner(tourId,userId,role,index);

        KeyPoint keyPoint = tour.getKeyPoints().get(index);
        keyPoint.setName(dto.getName());
        keyPoint.setDescription(dto.getDescription());
        keyPoint.setLatitude(dto.getLatitude());
        keyPoint.setLongitude(dto.getLongitude());

        if (image != null && !image.isEmpty()) {
            String fileName = UUID.randomUUID() + "_" +
                    StringUtils.cleanPath(image.getOriginalFilename());

            File uploadDir = new File("uploads");
            if (!uploadDir.exists()) {
                uploadDir.mkdirs();
            }

            File destination = new File(uploadDir, fileName);
            image.transferTo(destination);

            keyPoint.setImage(fileName);
        }

        return save(tour);
    }
    public Tour deleteKeyPoint(
            String tourId,
            int index,
            Long userId,
            String role
    ){
        Tour tour = checkTourOwner(tourId,userId,role,index);

        tour.getKeyPoints().remove(index);
        return save(tour);

    }

    public List<Tour> getToursByStatus(TourStatus status) {
        return tourRepository.findByStatus(status.name());
    }

    public Tour getTourPreview(String id) {
    Tour tour = getTourById(id);

    if (tour.getKeyPoints() == null || tour.getKeyPoints().isEmpty()) {
        tour.setKeyPoints(new ArrayList<>());
        return tour;
    }

    KeyPoint first = tour.getKeyPoints().get(0);
    KeyPoint preview = new KeyPoint();
    preview.setName(first.getName());
    preview.setImage(first.getImage());
    preview.setDescription("Unlock the full route after purchase");
    preview.setLatitude(first.getLatitude());
    preview.setLongitude(first.getLongitude());

    tour.setKeyPoints(Collections.singletonList(preview));
    return tour;
}
}