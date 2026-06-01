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

            com.tourism.tours.model.Tour created = tourService.createTour(dto, request.getAuthorId());

            Tour responseTour = Tour.newBuilder()
                    .setId(created.getId() == null ? "" : created.getId())
                    .setAuthorId(created.getAuthorId() == null ? 0L : created.getAuthorId())
                    .setTitle(created.getTitle() == null ? "" : created.getTitle())
                    .setDescription(created.getDescription() == null ? "" : created.getDescription())
                    .setDifficulty(created.getDifficulty() == null
                            ? Difficulty.DIFFICULTY_UNSPECIFIED
                            : Difficulty.valueOf(created.getDifficulty().name()))
                    .addAllTags(created.getTags() == null ? java.util.Collections.emptyList() : created.getTags())
                    .setPrice(created.getPrice() == null ? 0.0 : created.getPrice())
                    .setStatus(created.getStatus() == null ? "" : created.getStatus().name())
                    .build();

            responseObserver.onNext(CreateTourResponse.newBuilder()
                    .setTour(responseTour)
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
}