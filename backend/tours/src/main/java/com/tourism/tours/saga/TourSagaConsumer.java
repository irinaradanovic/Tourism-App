package com.tourism.tours.saga;

import com.tourism.tours.dto.OrderCreatedEvent;
import com.tourism.tours.dto.StartTourResultEvent;
import com.tourism.tours.dto.TourLifecycleResultEvent;
import com.tourism.tours.dto.ToursValidationResultEvent;
import com.tourism.tours.model.Tour;
import com.tourism.tours.model.TourStatus;
import com.tourism.tours.service.TourExecutionService;
import com.tourism.tours.service.TourService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.context.annotation.Lazy;
import org.springframework.stereotype.Component;

@Component
@Slf4j
public class TourSagaConsumer {

    private final TourService tourService;
    private final RabbitTemplate rabbitTemplate;
    private final TourExecutionService tourExecutionService;

    public TourSagaConsumer(TourService tourService,
                            RabbitTemplate rabbitTemplate,
                            @Lazy TourExecutionService tourExecutionService) {
        this.tourService = tourService;
        this.rabbitTemplate = rabbitTemplate;
        this.tourExecutionService = tourExecutionService;
    }

    @RabbitListener(queues = "order.created")
    public void handleOrderCreated(OrderCreatedEvent event) {
        log.info("[SAGA-JAVA] Accepting request for tour validation for cart: {}", event.getCartId());

        try {
            for (String tourId : event.getTourIds()) {
                Tour tour = tourService.getTourById(tourId);

                if (tour.getStatus() != TourStatus.PUBLISHED) {
                    log.warn("[SAGA-JAVA] Tour {} is not PUBLISHED! Status: {}. Failed.", tourId, tour.getStatus());

                    ToursValidationResultEvent failedEvent = new ToursValidationResultEvent(
                            event.getCartId(), event.getTouristId(), "Tour " + tourId + " is not PUBLISHED"
                    );
                    rabbitTemplate.convertAndSend("tours.failed", failedEvent);
                    return;
                }
            }
            log.info("[SAGA-JAVA] All tours for tourist {} are valid (PUBLISHED). Sending success.", event.getTouristId());
            ToursValidationResultEvent successEvent = new ToursValidationResultEvent(
                    event.getCartId(), event.getTouristId(), "All tours valid"
            );
            rabbitTemplate.convertAndSend("tours.validated", successEvent);

        } catch (Exception e) {
            log.error("[SAGA-JAVA] Error occurred while searching for tours: {}. Initiating failure.", e.getMessage());
            ToursValidationResultEvent failedEvent = new ToursValidationResultEvent(
                    event.getCartId(), event.getTouristId(), "Tour not found in MongoDB"
            );
            rabbitTemplate.convertAndSend("tours.failed", failedEvent);
        }
    }

    @RabbitListener(queues = "tour.publish.completed")
    public void handleTourPublishCompleted(TourLifecycleResultEvent event) {
        log.info("[SAGA-JAVA] Publish result received for tour {} success={}", event.getTourId(), event.getSuccess());
        if (Boolean.TRUE.equals(event.getSuccess())) {
            tourService.finalizePublish(event.getTourId());
            log.info("[SAGA-JAVA] Tour {} moved to PUBLISHED", event.getTourId());
        } else {
            log.warn("[SAGA-JAVA] Publish saga failed for tour {}: {}", event.getTourId(), event.getReason());
        }
    }

    @RabbitListener(queues = "tour.archive.completed")
    public void handleTourArchiveCompleted(TourLifecycleResultEvent event) {
        log.info("[SAGA-JAVA] Archive result received for tour {} success={}", event.getTourId(), event.getSuccess());
        if (Boolean.TRUE.equals(event.getSuccess())) {
            tourService.finalizeArchive(event.getTourId());
            log.info("[SAGA-JAVA] Tour {} moved to ARCHIVED", event.getTourId());
        } else {
            log.warn("[SAGA-JAVA] Archive saga failed for tour {}: {}", event.getTourId(), event.getReason());
        }
    }

    @RabbitListener(queues = "start.tour.result")
    public void handleStartTourResult(StartTourResultEvent event) {
        log.info("[SAGA-START-TOUR] Primljen odgovor od Purchase servisa. correlationId={} purchased={}",
                event.getCorrelationId(), event.getPurchased());
        tourExecutionService.handleStartTourResult(event);
    }
}