package com.tourism.tours.config;

import org.springframework.amqp.core.Queue;
import org.springframework.amqp.support.converter.Jackson2JsonMessageConverter;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class RabbitMqConfig {

    @Bean
    public Queue orderCreatedQueue() { return new Queue("order.created", true); }

    @Bean
    public Queue toursValidatedQueue() { return new Queue("tours.validated", true); }

    @Bean
    public Queue toursFailedQueue() { return new Queue("tours.failed", true); }

    @Bean
    public Jackson2JsonMessageConverter messageConverter() {
        return new Jackson2JsonMessageConverter();
    }
}