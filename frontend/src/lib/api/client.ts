import axios, { AxiosInstance, AxiosResponse, AxiosError } from 'axios'

// Create axios instance with base configuration
const apiClient: AxiosInstance = axios.create({
    baseURL: process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8000/api',
    timeout: 30000,
    headers: {
        'Content-Type': 'application/json',
    }
})

// Request interceptor to add auth token
apiClient.interceptors.request.use(
    async (config) => {
        // Add auth token if available
        if (typeof window !== 'undefined') {
            const token = localStorage.getItem('authToken')
            if (token) {
                config.headers.Authorization = `Bearer ${token}`
            }
        }
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

// Response interceptor for error handling
apiClient.interceptors.response.use(
    (response: AxiosResponse) => {
        return response
    },
    (error: AxiosError) => {
        // Handle common errors
        if (error.response?.status === 401) {
            // Unauthorized - redirect to login
            if (typeof window !== 'undefined') {
                localStorage.removeItem('authToken')
                window.location.href = '/login'
            }
        }

        return Promise.reject(error)
    }
)

// API endpoint functions
export const api = {
    // Authentication
    auth: {
        login: (credentials: { email: string; password: string }) =>
            apiClient.post('/auth/login', credentials),

        register: (userData: { email: string; password: string; firstName: string; lastName: string }) =>
            apiClient.post('/auth/register', userData),

        logout: () =>
            apiClient.post('/auth/logout'),

        refreshToken: () =>
            apiClient.post('/auth/refresh'),

        forgotPassword: (email: string) =>
            apiClient.post('/auth/forgot-password', { email }),

        resetPassword: (token: string, password: string) =>
            apiClient.post('/auth/reset-password', { token, password }),
    },

    // User profile
    user: {
        getProfile: () =>
            apiClient.get('/user/profile'),

        updateProfile: (profileData: any) =>
            apiClient.put('/user/profile', profileData),

        uploadAvatar: (file: FormData) =>
            apiClient.post('/user/avatar', file, {
                headers: { 'Content-Type': 'multipart/form-data' }
            }),

        deleteAccount: () =>
            apiClient.delete('/user/account'),
    },

    // Trip planning
    trips: {
        getAll: (params?: { page?: number; limit?: number; status?: string }) =>
            apiClient.get('/trips', { params }),

        getById: (tripId: string) =>
            apiClient.get(`/trips/${tripId}`),

        create: (tripData: any) =>
            apiClient.post('/trips', tripData),

        update: (tripId: string, tripData: any) =>
            apiClient.put(`/trips/${tripId}`, tripData),

        delete: (tripId: string) =>
            apiClient.delete(`/trips/${tripId}`),

        generateItinerary: (tripId: string, preferences: any) =>
            apiClient.post(`/trips/${tripId}/generate-itinerary`, preferences),

        getItinerary: (tripId: string) =>
            apiClient.get(`/trips/${tripId}/itinerary`),

        updateItinerary: (tripId: string, itineraryData: any) =>
            apiClient.put(`/trips/${tripId}/itinerary`, itineraryData),
    },

    // AI-powered trip planning with RAG
    ai: {
        planTrip: (tripData: {
            destination: string
            start_date: string
            end_date: string
            budget: number
            travelers: number
            preferences: Record<string, any>
            trip_type?: string
            interests: string[]
            user_id: string
        }) =>
            apiClient.post('/ai/plan-trip', tripData),

        getRecommendations: (params: {
            user_id?: string
            budget?: number
            interests: string[]
        }) =>
            apiClient.get('/ai/recommendations', { params }),

        getRagContext: (params: {
            destination: string
            user_id?: string
            start_date: string
            end_date: string
            budget: number
            travelers: number
            interests: string[]
            preferences: Record<string, any>
        }) =>
            apiClient.post('/ai/rag-context', params),

        validateAvailability: (items: {
            type: 'attractions' | 'hotels' | 'transportation'
            data: any[]
            check_date: string
        }) =>
            apiClient.post('/ai/validate-availability', items),
    },

    // Vector database operations
    vector: {
        searchSimilarAttractions: (params: {
            interests: string[]
            limit?: number
        }) =>
            apiClient.post('/vector/search-attractions', params),

        searchSimilarTrips: (params: {
            destination: string
            preferences: Record<string, any>
            limit?: number
        }) =>
            apiClient.post('/vector/search-trips', params),

        storeUserPreferences: (userProfile: {
            user_id: string
            preferences: Record<string, any>
            trip_history?: string[]
        }) =>
            apiClient.post('/vector/store-preferences', userProfile),
    },

    // Recommendations
    recommendations: {
        getDestinations: (preferences: any) =>
            apiClient.post('/recommendations/destinations', preferences),

        getActivities: (destination: string, preferences: any) =>
            apiClient.post('/recommendations/activities', { destination, preferences }),

        getAccommodations: (destination: string, filters: any) =>
            apiClient.post('/recommendations/accommodations', { destination, filters }),

        getRestaurants: (destination: string, preferences: any) =>
            apiClient.post('/recommendations/restaurants', { destination, preferences }),
    },

    // Search
    search: {
        destinations: (query: string) =>
            apiClient.get(`/search/destinations?q=${encodeURIComponent(query)}`),

        activities: (destination: string, query: string) =>
            apiClient.get(`/search/activities?destination=${encodeURIComponent(destination)}&q=${encodeURIComponent(query)}`),

        accommodations: (destination: string, checkIn: string, checkOut: string, guests: number) =>
            apiClient.get(`/search/accommodations?destination=${encodeURIComponent(destination)}&checkIn=${checkIn}&checkOut=${checkOut}&guests=${guests}`),
    },

    // Weather
    weather: {
        getCurrent: (destination: string) =>
            apiClient.get(`/weather/current?destination=${encodeURIComponent(destination)}`),

        getForecast: (destination: string, days?: number) =>
            apiClient.get(`/weather/forecast?destination=${encodeURIComponent(destination)}&days=${days || 7}`),
    },

    // Currency
    currency: {
        getRates: (from: string, to: string) =>
            apiClient.get(`/currency/rates?from=${from}&to=${to}`),

        convert: (amount: number, from: string, to: string) =>
            apiClient.get(`/currency/convert?amount=${amount}&from=${from}&to=${to}`),
    },

    // Bookings (if integrated)
    bookings: {
        getAll: () =>
            apiClient.get('/bookings'),

        getById: (bookingId: string) =>
            apiClient.get(`/bookings/${bookingId}`),

        create: (bookingData: any) =>
            apiClient.post('/bookings', bookingData),

        cancel: (bookingId: string) =>
            apiClient.delete(`/bookings/${bookingId}`),
    },

    // Notifications
    notifications: {
        getAll: () =>
            apiClient.get('/notifications'),

        markAsRead: (notificationId: string) =>
            apiClient.put(`/notifications/${notificationId}/read`),

        markAllAsRead: () =>
            apiClient.put('/notifications/read-all'),

        updateSettings: (settings: any) =>
            apiClient.put('/notifications/settings', settings),
    },

    // Analytics (for admin/insights)
    analytics: {
        getTripStats: () =>
            apiClient.get('/analytics/trips'),

        getUserActivity: () =>
            apiClient.get('/analytics/user-activity'),

        getPopularDestinations: () =>
            apiClient.get('/analytics/popular-destinations'),
    },
}

// Helper functions
export const handleApiError = (error: AxiosError): string => {
    if (error.response) {
        // Server responded with error status
        const message = (error.response.data as any)?.message || 'An error occurred'
        return message
    } else if (error.request) {
        // Request was made but no response received
        return 'Network error. Please check your connection.'
    } else {
        // Something else happened
        return 'An unexpected error occurred.'
    }
}

export const isApiError = (error: any): error is AxiosError => {
    return error.isAxiosError === true
}

export default apiClient