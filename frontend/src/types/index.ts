// User types
export interface User {
    uid: string
    email: string
    displayName?: string
    photoURL?: string
    createdAt: Date
    lastLoginAt?: Date
}

// Trip types
export interface Trip {
    id: string
    userId: string
    title: string
    destination: string
    startDate: string
    endDate: string
    status: 'draft' | 'planned' | 'ongoing' | 'completed' | 'cancelled'
    travelers: number
    budget: {
        total: number
        currency: string
        breakdown?: {
            accommodation: number
            transportation: number
            food: number
            activities: number
            other: number
        }
    }
    preferences: TripPreferences
    itinerary?: Itinerary
    createdAt: Date
    updatedAt: Date
}

export interface TripPreferences {
    travelStyle: string[]
    interests: string[]
    accommodation: string
    transportation: string
    foodPreferences?: string[]
    accessibility?: string[]
    language?: string
}

// Itinerary types
export interface Itinerary {
    id: string
    tripId: string
    days: DayPlan[]
    totalActivities: number
    estimatedCost: number
    generatedAt: Date
}

export interface DayPlan {
    date: string
    activities: Activity[]
    meals: Meal[]
    accommodation?: Accommodation
    transportation?: Transportation[]
    notes?: string
}

export interface Activity {
    id: string
    name: string
    description: string
    type: ActivityType
    location: Location
    duration: number // in minutes
    cost: number
    rating?: number
    bookingRequired: boolean
    bookingUrl?: string
    openingHours?: OpeningHours
    tips?: string[]
    images?: string[]
}

export interface Meal {
    id: string
    name: string
    type: 'breakfast' | 'lunch' | 'dinner' | 'snack'
    location: Location
    cuisine?: string
    cost: number
    rating?: number
    reservationRequired?: boolean
    dietaryOptions?: string[]
}

export interface Accommodation {
    id: string
    name: string
    type: 'hotel' | 'resort' | 'apartment' | 'hostel' | 'villa' | 'boutique'
    location: Location
    pricePerNight: number
    rating?: number
    amenities: string[]
    images?: string[]
    bookingUrl?: string
    checkIn?: string
    checkOut?: string
}

export interface Transportation {
    id: string
    type: 'flight' | 'train' | 'bus' | 'car' | 'taxi' | 'metro' | 'walk'
    from: Location
    to: Location
    departureTime?: string
    arrivalTime?: string
    duration?: number
    cost?: number
    bookingUrl?: string
    notes?: string
}

// Common types
export interface Location {
    name: string
    address?: string
    city: string
    country: string
    coordinates?: {
        latitude: number
        longitude: number
    }
    timezone?: string
}

export interface OpeningHours {
    monday?: TimeSlot[]
    tuesday?: TimeSlot[]
    wednesday?: TimeSlot[]
    thursday?: TimeSlot[]
    friday?: TimeSlot[]
    saturday?: TimeSlot[]
    sunday?: TimeSlot[]
}

export interface TimeSlot {
    open: string
    close: string
}

export type ActivityType =
    | 'sightseeing'
    | 'museum'
    | 'adventure'
    | 'cultural'
    | 'entertainment'
    | 'shopping'
    | 'nature'
    | 'food'
    | 'nightlife'
    | 'sports'
    | 'religious'
    | 'historical'

// API Response types
export interface ApiResponse<T> {
    success: boolean
    data?: T
    error?: string
    message?: string
}

export interface PaginatedResponse<T> {
    data: T[]
    pagination: {
        page: number
        limit: number
        total: number
        totalPages: number
    }
}

// Form types
export interface TripPlanningForm {
    destination: string
    startDate: string
    endDate: string
    travelers: number
    budget: string
    travelStyle: string[]
    interests: string[]
    accommodation: string
    transportation: string
    additionalRequests?: string
}

export interface UserProfile {
    uid: string
    email: string
    firstName: string
    lastName: string
    dateOfBirth?: string
    phone?: string
    nationality?: string
    preferredCurrency: string
    preferredLanguage: string
    travelPreferences: {
        budgetRange: string
        travelStyles: string[]
        interests: string[]
        accommodationTypes: string[]
        transportationModes: string[]
    }
    emergencyContact?: {
        name: string
        relationship: string
        phone: string
        email?: string
    }
    notifications: {
        email: boolean
        sms: boolean
        push: boolean
        tripReminders: boolean
        priceAlerts: boolean
        recommendations: boolean
    }
}

// Search and filter types
export interface SearchFilters {
    destination?: string
    dateRange?: {
        start: string
        end: string
    }
    budget?: {
        min: number
        max: number
    }
    travelers?: number
    activities?: string[]
    rating?: number
}

export interface Recommendation {
    id: string
    type: 'destination' | 'activity' | 'restaurant' | 'accommodation'
    title: string
    description: string
    location: Location
    images: string[]
    rating: number
    price?: number
    tags: string[]
    reasons: string[]
    confidence: number
}

// Weather types
export interface Weather {
    date: string
    temperature: {
        min: number
        max: number
        unit: 'celsius' | 'fahrenheit'
    }
    condition: string
    description: string
    humidity: number
    windSpeed: number
    precipitation: number
    icon: string
}

// Currency types
export interface CurrencyRate {
    from: string
    to: string
    rate: number
    lastUpdated: Date
}

// Error types
export interface AppError {
    code: string
    message: string
    details?: any
    timestamp: Date
}

// Theme types
export interface ThemeConfig {
    primary: string
    secondary: string
    accent: string
    background: string
    surface: string
    text: string
    error: string
    warning: string
    success: string
    info: string
}