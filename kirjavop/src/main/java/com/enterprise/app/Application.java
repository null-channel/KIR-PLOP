package com.enterprise.app;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.ConfigurationPropertiesScan;
import org.springframework.data.jpa.repository.config.EnableJpaAuditing;

@SpringBootApplication
@EnableJpaAuditing
@ConfigurationPropertiesScan
public final class Application {
    
    private Application() {
        // Private constructor to prevent instantiation
    }
    
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
} 
