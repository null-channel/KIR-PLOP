package com.enterprise.app.core.port.in;

public interface CommandUseCase<T> {
    void execute(T command);
} 