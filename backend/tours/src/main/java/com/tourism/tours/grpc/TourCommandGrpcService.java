package com.tourism.tours.grpc;

import com.tourism.tours.dto.CreateTourDTO;
import com.tourism.tours.service.TourService;
import io.grpc.Status;
import io.grpc.stub.StreamObserver;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import net.devh.boot.grpc.server.service.GrpcService;

@GrpcService
@RequiredArgsConstructor
@Slf4j
public class TourCommandGrpcService extends TourCommandServiceGrpc.TourCommandServiceImplBase {

    private final TourService tourService;

    @Override
    public void createTour(CreateTourRequest request, StreamObserver<CreateTourResponse> responseObserver) {
        log.info("========== CreateTour gRPC Request ==========");
        log.info("Accepted create tour request for author: {}", request.getAuthorId());

        try {
            CreateTourDTO dto = new CreateTourDTO();
            dto.setTitle(request.getTitle());
            dto.setDescription(request.getDescription());

            if (request.getDifficulty() != Difficulty.DIFFICULTY_UNSPECIFIED) {
                dto.setDifficulty(com.tourism.tours.model.Difficulty.valueOf(request.getDifficulty().name()));
            }

            dto.setTags(request.getTagsList());
            dto.setDurations(request.getDurationsList().stream()
                    .filter(d -> d.getTransportType() != TransportType.TRANSPORT_UNSPECIFIED)
                    .map(d -> {
                        com.tourism.tours.model.TourDuration duration = new com.tourism.tours.model.TourDuration();
                        duration.setTransportType(com.tourism.tours.model.TransportType.valueOf(d.getTransportType().name()));
                        duration.setMinutes(d.getMinutes());
                        return duration;
                    })
                    .toList());

            com.tourism.tours.model.Tour created = tourService.createTour(dto, request.getAuthorId());

            responseObserver.onNext(CreateTourResponse.newBuilder()
                    .setTour(toGrpcTour(created))
                    .build());
            responseObserver.onCompleted();
            log.info("========== CreateTour gRPC Response Sent ==========");
        } catch (Exception e) {
            log.error("Error in CreateTour gRPC service: {}", e.getMessage(), e);
            responseObserver.onError(Status.INTERNAL
                    .withDescription("Failed to create tour")
                    .withCause(e)
                    .asRuntimeException());
        }
    }

    @Override
    public void publishTour(PublishTourRequest request, StreamObserver<TourCommandResponse> responseObserver) {
        log.info("========== PublishTour gRPC Request ==========");
        log.info("Accepted publish tour request for tour: {}", request.getTourId());

        try {
            com.tourism.tours.model.Tour tour = tourService.requestPublish(
                    request.getTourId(),
                    request.getAuthorId(),
                    request.getRole()
            );

            responseObserver.onNext(TourCommandResponse.newBuilder()
                    .setTour(toGrpcTour(tour))
                    .setMessage("Publish saga started")
                    .build());
            responseObserver.onCompleted();
            log.info("========== PublishTour gRPC Response Sent ==========");
        } catch (Exception e) {
            log.error("Error in PublishTour gRPC service: {}", e.getMessage(), e);
            responseObserver.onError(Status.INTERNAL
                    .withDescription(e.getMessage())
                    .withCause(e)
                    .asRuntimeException());
        }
    }

    private Tour toGrpcTour(com.tourism.tours.model.Tour tour) {
        Tour.Builder builder = Tour.newBuilder()
                .setId(tour.getId() == null ? "" : tour.getId())
                .setAuthorId(tour.getAuthorId() == null ? 0L : tour.getAuthorId())
                .setTitle(tour.getTitle() == null ? "" : tour.getTitle())
                .setDescription(tour.getDescription() == null ? "" : tour.getDescription())
                .setDifficulty(tour.getDifficulty() == null
                        ? Difficulty.DIFFICULTY_UNSPECIFIED
                        : Difficulty.valueOf(tour.getDifficulty().name()))
                .addAllTags(tour.getTags() == null ? java.util.Collections.emptyList() : tour.getTags())
                .setPrice(tour.getPrice() == null ? 0.0 : tour.getPrice())
                .setDistanceKm(tour.getDistanceKm() == null ? 0.0 : tour.getDistanceKm())
                .setStatus(tour.getStatus() == null ? "" : tour.getStatus().name())
                .setPublishedAt(tour.getPublishedAt() == null ? "" : tour.getPublishedAt().toString())
                .setArchivedAt(tour.getArchivedAt() == null ? "" : tour.getArchivedAt().toString());

        if (tour.getDurations() != null) {
            for (com.tourism.tours.model.TourDuration duration : tour.getDurations()) {
                if (duration.getTransportType() != null && duration.getMinutes() != null) {
                    builder.addDurations(TourDuration.newBuilder()
                            .setTransportType(TransportType.valueOf(duration.getTransportType().name()))
                            .setMinutes(duration.getMinutes())
                            .build());
                }
            }
        }

        return builder.build();
    }
}
