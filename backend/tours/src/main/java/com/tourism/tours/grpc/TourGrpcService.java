package com.tourism.tours.grpc;

import com.tourism.tours.model.Tour;
import com.tourism.tours.service.TourService;
import io.grpc.stub.StreamObserver;
import net.devh.boot.grpc.server.service.GrpcService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@GrpcService // annotation to register this class as a gRPC service in Spring Boot
@RequiredArgsConstructor
@Slf4j
public class TourGrpcService extends TourCheckServiceGrpc.TourCheckServiceImplBase {

    private final TourService tourService;

    @Override
    public void checkTour(TourCheckRequest request, StreamObserver<TourCheckResponse> responseObserver) {
        log.info("========== gRPC Request  ==========");
        log.info("Accepted request for tour with ID: {}", request.getTourId());

        try {
            Tour tour = tourService.getTourById(request.getTourId());
            log.info("Tour successfully found in the database: {} - Price: {}", tour.getTitle(), tour.getPrice());
            
            // packing data from MongoDB into gRPC response
            TourCheckResponse response = TourCheckResponse.newBuilder()
                    .setStatus(tour.getStatus().toString()) 
                    .setTourName(tour.getTitle())           
                    .setPrice(tour.getPrice())
                    .build();
            
            // sending response back to gRPC client, to purchase service
            responseObserver.onNext(response);
            responseObserver.onCompleted();
            log.info("========== gRPC Response Sent ==========");
            
        } catch (Exception e) {
            // if tour is not found, send NOT_FOUND status back to client
            log.error("Error in gRPC service for ID: {}. Reason: {}", request.getTourId(), e.getMessage());
            responseObserver.onError(io.grpc.Status.NOT_FOUND
                    .withDescription("Tour not found: " + request.getTourId())
                    .asRuntimeException());
        }
    }
}