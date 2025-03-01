package com.enterprise.app.core.port.in;

public interface QueryUseCase<T, R> {
    R execute(T query);
} 