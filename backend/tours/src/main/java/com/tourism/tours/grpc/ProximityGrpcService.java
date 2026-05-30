package com.tourism.tours.grpc;

import com.tourism.tours.model.TourExecution;
import com.tourism.tours.service.TourExecutionService;
import io.grpc.stub.StreamObserver;
import net.devh.boot.grpc.server.service.GrpcService;
import lombok.RequiredArgsConstructor;

import java.util.Map;

@GrpcService
@RequiredArgsConstructor
public class ProximityGrpcService extends ProximityServiceGrpc.ProximityServiceImplBase {

    private final TourExecutionService executionService;

    @Override
    public void checkProximity(ProximityRequest request, StreamObserver<ProximityResponse> responseObserver) {
        TourExecution execution = executionService.checkProximity(
                request.getExecutionId(),
                request.getTouristId(),
                request.getLat(),
                request.getLon()
        );

        ProximityResponse.Builder builder = ProximityResponse.newBuilder()
                .setExecutionId(execution.getId())
                .setStatus(execution.getStatus())
                .setLastActivity(execution.getLastActivity().toString());

        for (Map.Entry<Integer, java.time.LocalDateTime> entry : execution.getCompletedKeyPoints().entrySet()) {
            builder.addCompletedKeyPoints(
                    CompletedKeyPoint.newBuilder()
                            .setIndex(entry.getKey())
                            .setCompletedAt(entry.getValue().toString())
                            .build()
            );
        }

        responseObserver.onNext(builder.build());
        responseObserver.onCompleted();
    }
}