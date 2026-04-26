// services/csrf.service.js

import api from "./api.js";

class CSRFService {
    constructor(axiosInstance) {
        this.api = axiosInstance;
        this.csrfToken = null;
        this.isInitialized = false;
        this.isRefreshing = false;
        this.pendingRequests = [];

        // Настройка interceptors
        this.setupInterceptors();
    }

    /**
     * Инициализация CSRF токена при старте приложения
     */
    async init() {
        if (this.isInitialized) {
            return this.csrfToken;
        }

        try {
            const response = await this.api.get('/api/user/profile');
            const token = response.headers['x-csrf-token'];
            if (token) {
                this.setToken(token);
                this.isInitialized = true;
                return token;
            } else {
                throw new Error('No CSRF token in response headers');
            }
        } catch (error) {
            console.error('Failed to initialize CSRF token:', error);
            throw error;
        }
    }

    /**
     * Установка токена в заголовки axios
     */
    setToken(token) {
        this.csrfToken = token;
        this.api.defaults.headers.common['X-CSRF-Token'] = token;
    }

    /**
     * Получение текущего токена
     */
    getToken() {
        return this.csrfToken;
    }

    /**
     * Сброс токена (при ошибке)
     */
    resetToken() {
        this.csrfToken = null;
        this.isInitialized = false;
        delete this.api.defaults.headers.common['X-CSRF-Token'];
    }

    /**
     * Обновление токена
     */
    async refreshToken() {
        if (this.isRefreshing) {
            // Если уже обновляется, возвращаем promise ожидания
            return new Promise((resolve, reject) => {
                this.pendingRequests.push({resolve, reject});
            });
        }

        this.isRefreshing = true;

        try {
            const response = await this.api.get('/api/user/profile');
            const newToken = response.headers['x-csrf-token'];
            if (newToken) {
                this.setToken(newToken);
                this.isInitialized = true;
                // Обрабатываем все ожидающие запросы
                this.pendingRequests.forEach(req => req.resolve(newToken));
                this.pendingRequests = [];

                return newToken;
            } else {
                throw new Error('No CSRF token in refresh response');
            }
        } catch (error) {
            console.error('Failed to refresh CSRF token:', error);

            // Обрабатываем ошибки для ожидающих запросов
            this.pendingRequests.forEach(req => req.reject(error));
            this.pendingRequests = [];

            throw error;
        } finally {
            this.isRefreshing = false;
        }
    }

    /**
     * Настройка interceptors для автоматического обновления токена
     */
    setupInterceptors() {
        // Request interceptor - добавляем токен в каждый запрос
        this.api.interceptors.request.use(
            config => {
                // Если токен есть, добавляем его в заголовки запроса
                if (this.csrfToken && this.requiresCSRF(config)) {
                    config.headers['X-CSRF-Token'] = this.csrfToken;
                }
                return config;
            },
            error => Promise.reject(error)
        );

        // Response interceptor - обрабатываем ошибки CSRF
        this.api.interceptors.response.use(
            response => response,
            async error => {
                const originalRequest = error.config;

                // Проверяем, что это ошибка CSRF (403) и запрос не повторялся
                if (error.response?.status === 403 &&
                    !originalRequest._retry &&
                    this.requiresCSRF(originalRequest)) {

                    originalRequest._retry = true;

                    try {
                        // Обновляем токен
                        const newToken = await this.refreshToken();

                        // Обновляем заголовок в оригинальном запросе
                        originalRequest.headers['X-CSRF-Token'] = newToken;

                        // Повторяем оригинальный запрос
                        return this.api(originalRequest);
                    } catch (refreshError) {
                        // Если не удалось обновить токен, сбрасываем его
                        this.resetToken();

                        // Можно перенаправить на страницу логина или показать ошибку
                        window.dispatchEvent(new CustomEvent('csrf:error', {
                            detail: refreshError
                        }));

                        return Promise.reject(refreshError);
                    }
                }

                return Promise.reject(error);
            }
        );
    }

    /**
     * Проверяет, нужен ли CSRF токен для данного запроса
     */
    requiresCSRF(config) {
        const methodsThatNeedCSRF = ['post', 'put', 'patch', 'delete'];
        return methodsThatNeedCSRF.includes(config.method?.toLowerCase());
    }
}

// Создаем и экспортируем экземпляр сервиса
let csrfInstance = null;

export function getCSRFService(axiosInstance) {
    if (!csrfInstance) {
        csrfInstance = new CSRFService(axiosInstance);
    }
    return csrfInstance;
}