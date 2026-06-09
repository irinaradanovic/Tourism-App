package com.tourism.tours.saga;
import com.tourism.tours.dto.OrderCreatedEvent;
import com.tourism.tours.dto.ToursValidationResultEvent;
import com.tourism.tours.model.Tour;
import com.tourism.tours.model.TourStatus;
import com.tourism.tours.service.TourService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
@Slf4j
public class TourSagaConsumer {
private final TourService tourService;
    private final RabbitTemplate rabbitTemplate;

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
}
