package com.tourism.tours.service;

import com.tourism.tours.dto.CreateTourDTO;
import com.tourism.tours.model.Tour;
import com.tourism.tours.repository.TourRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;

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
}