package com.tourism.tours.service;

import com.tourism.tours.dto.AddReviewDTO;
import com.tourism.tours.dto.CreateTourDTO;
import com.tourism.tours.dto.ReviewDTO;
import com.tourism.tours.dto.TourLifecycleRequestedEvent;
import com.tourism.tours.model.Review;
import com.tourism.tours.model.Tour;
import com.tourism.tours.model.TourDuration;
import com.tourism.tours.repository.TourRepository;
import com.tourism.tours.dto.CreateKeyPointDTO;
import com.tourism.tours.model.KeyPoint;
import lombok.RequiredArgsConstructor;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
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
    private final RabbitTemplate rabbitTemplate;

    public Tour createTour(CreateTourDTO dto, Long authorId) {

        Tour tour = new Tour();

        tour.setAuthorId(authorId);
        tour.setTitle(dto.getTitle());
        tour.setDescription(dto.getDescription());
        tour.setDifficulty(dto.getDifficulty());
        tour.setStatus(TourStatus.DRAFT);
        tour.setTags(dto.getTags() == null ? new ArrayList<>() : dto.getTags());
        tour.setDurations(normalizeDurations(dto.getDurations()));

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

    public Tour addKeyPoint(
            String tourId,
            CreateKeyPointDTO dto,
            MultipartFile image,
            Long userId,
            String role
    ) throws IOException {
        Tour tour = getTourById(tourId);
        ensureGuideOwner(tour, userId, role, "Only guide owner can add key points");

        if (tour.getKeyPoints() == null) {
            tour.setKeyPoints(new ArrayList<>());
        }

        KeyPoint keyPoint = new KeyPoint();
        keyPoint.setName(dto.getName());
        keyPoint.setDescription(dto.getDescription());
        keyPoint.setLatitude(dto.getLatitude());
        keyPoint.setLongitude(dto.getLongitude());
        keyPoint.setImage(storeImage(image));

        tour.getKeyPoints().add(keyPoint);
        recalculateDistance(tour);
        return save(tour);
    }

    public Tour updateDurations(String tourId, List<TourDuration> durations, Long userId, String role) {
        Tour tour = getTourById(tourId);
        ensureGuideOwner(tour, userId, role, "Only guide owner can update tour durations");
        tour.setDurations(normalizeDurations(durations));
        return save(tour);
    }

    public Tour requestPublish(String tourId, Long userId, String role) {
        Tour tour = getTourById(tourId);
        ensureGuideOwner(tour, userId, role, "Only guide owner can publish tour");
        validatePublishRequirements(tour);
        rabbitTemplate.convertAndSend("tour.publish.requested",
                new TourLifecycleRequestedEvent(tour.getId(), tour.getAuthorId(), "PUBLISH"));
        return tour;
    }

    public Tour finalizePublish(String tourId) {
        Tour tour = getTourById(tourId);
        if (tour.getStatus() == TourStatus.PUBLISHED) {
            return tour;
        }
        validatePublishRequirements(tour);
        tour.setStatus(TourStatus.PUBLISHED);
        tour.setPublishedAt(LocalDateTime.now());
        tour.setArchivedAt(null);
        return save(tour);
    }

    public Tour requestArchive(String tourId, Long userId, String role) {
        Tour tour = getTourById(tourId);
        ensureGuideOwner(tour, userId, role, "Only guide owner can archive tour");
        if (tour.getStatus() != TourStatus.PUBLISHED) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Only published tours can be archived");
        }
        rabbitTemplate.convertAndSend("tour.archive.requested",
                new TourLifecycleRequestedEvent(tour.getId(), tour.getAuthorId(), "ARCHIVE"));
        return tour;
    }

    public Tour finalizeArchive(String tourId) {
        Tour tour = getTourById(tourId);
        if (tour.getStatus() == TourStatus.ARCHIVED) {
            return tour;
        }
        if (tour.getStatus() != TourStatus.PUBLISHED) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Only published tours can be archived");
        }
        tour.setStatus(TourStatus.ARCHIVED);
        tour.setArchivedAt(LocalDateTime.now());
        return save(tour);
    }

    public Tour reactivateTour(String tourId, Long userId, String role) {
        Tour tour = getTourById(tourId);
        ensureGuideOwner(tour, userId, role, "Only guide owner can reactivate tour");
        if (tour.getStatus() != TourStatus.ARCHIVED) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Only archived tours can be reactivated");
        }
        validatePublishRequirements(tour);
        tour.setStatus(TourStatus.PUBLISHED);
        tour.setPublishedAt(LocalDateTime.now());
        tour.setArchivedAt(null);
        return save(tour);
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

        String fileName = storeImage(image);
        if (fileName != null) {
            keyPoint.setImage(fileName);
        }

        recalculateDistance(tour);
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
        recalculateDistance(tour);
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

    private void ensureGuideOwner(Tour tour, Long userId, String role, String message) {
        if (!"GUIDE".equals(role) || userId == null || !userId.equals(tour.getAuthorId())) {
            throw new ResponseStatusException(HttpStatus.FORBIDDEN, message);
        }
    }

    private void validatePublishRequirements(Tour tour) {
        if (!StringUtils.hasText(tour.getTitle()) ||
                !StringUtils.hasText(tour.getDescription()) ||
                tour.getDifficulty() == null ||
                tour.getTags() == null ||
                tour.getTags().stream().noneMatch(StringUtils::hasText)) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Tour must contain title, description, difficulty and tags");
        }
        if (tour.getKeyPoints() == null || tour.getKeyPoints().size() < 2) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Tour must contain at least two key points");
        }
        if (tour.getDurations() == null || tour.getDurations().stream()
                .noneMatch(d -> d.getTransportType() != null && d.getMinutes() != null && d.getMinutes() > 0)) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "At least one tour duration by transport type is required");
        }
    }

    private List<TourDuration> normalizeDurations(List<TourDuration> durations) {
        if (durations == null) {
            return new ArrayList<>();
        }
        return durations.stream()
                .filter(d -> d != null && d.getTransportType() != null && d.getMinutes() != null && d.getMinutes() > 0)
                .toList();
    }

    private String storeImage(MultipartFile image) throws IOException {
        if (image == null || image.isEmpty()) {
            return null;
        }

        String fileName = UUID.randomUUID() + "_" + StringUtils.cleanPath(image.getOriginalFilename());
        File uploadDir = new File("uploads");
        if (!uploadDir.exists()) {
            uploadDir.mkdirs();
        }

        File destination = new File(uploadDir, fileName);
        image.transferTo(destination);
        return fileName;
    }

    private void recalculateDistance(Tour tour) {
        if (tour.getKeyPoints() == null || tour.getKeyPoints().size() < 2) {
            tour.setDistanceKm(0.0);
            return;
        }

        double total = 0.0;
        for (int i = 1; i < tour.getKeyPoints().size(); i++) {
            KeyPoint previous = tour.getKeyPoints().get(i - 1);
            KeyPoint current = tour.getKeyPoints().get(i);
            if (previous.getLatitude() != null && previous.getLongitude() != null &&
                    current.getLatitude() != null && current.getLongitude() != null) {
                total += haversineKm(previous.getLatitude(), previous.getLongitude(), current.getLatitude(), current.getLongitude());
            }
        }
        tour.setDistanceKm(Math.round(total * 100.0) / 100.0);
    }

    private double haversineKm(double lat1, double lon1, double lat2, double lon2) {
        final double earthRadiusKm = 6371.0;
        double dLat = Math.toRadians(lat2 - lat1);
        double dLon = Math.toRadians(lon2 - lon1);
        double a = Math.sin(dLat / 2) * Math.sin(dLat / 2)
                + Math.cos(Math.toRadians(lat1)) * Math.cos(Math.toRadians(lat2))
                * Math.sin(dLon / 2) * Math.sin(dLon / 2);
        double c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
        return earthRadiusKm * c;
    }
}
