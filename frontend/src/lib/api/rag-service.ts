import { api } from './client'
import { RAGTripContext, RAGItineraryResponse, WeatherForecast } from '@/types'

export interface RAGTripRequest {
    destination: string
    start_date: string
    end_date: string
    budget: number
    travelers: number
    preferences: Record<string, any>
    interests: string[]
    user_id: string
    trip_type?: string
}

export interface RAGContextRequest {
    destination: string
    user_id?: string
    start_date: string
    end_date: string
    budget: number
    travelers: number
    interests: string[]
    preferences: Record<string, any>
}

export interface ValidationCriteria {
    budget: number
    required_rating?: number
    preferred_types?: string[]
    accessibility?: boolean
    preferences?: Record<string, any>
    availability_check?: boolean
}

export interface RankingWeights {
    rating: number
    price: number
    distance: number
    availability: number
    user_match: number
}

class RAGService {
    /**
     * Plan a trip using RAG-enhanced AI
     */
    async planTrip(request: RAGTripRequest): Promise<RAGItineraryResponse> {
        try {
            const response = await api.ai.planTrip(request)
            return response.data
        } catch (error) {
            console.error('Error planning trip with RAG:', error)
            throw new Error('Failed to plan trip. Please try again.')
        }
    }

    /**
     * Get RAG context for a destination
     */
    async getRAGContext(request: RAGContextRequest): Promise<RAGTripContext> {
        try {
            const response = await api.ai.getRagContext(request)
            return response.data
        } catch (error) {
            console.error('Error fetching RAG context:', error)
            throw new Error('Failed to fetch travel context. Please try again.')
        }
    }

    /**
     * Search for similar attractions using vector similarity
     */
    async searchSimilarAttractions(interests: string[], limit: number = 10) {
        try {
            const response = await api.vector.searchSimilarAttractions({
                interests,
                limit
            })
            return response.data
        } catch (error) {
            console.error('Error searching similar attractions:', error)
            throw new Error('Failed to search attractions. Please try again.')
        }
    }

    /**
     * Search for similar trips using vector similarity
     */
    async searchSimilarTrips(destination: string, preferences: Record<string, any>, limit: number = 5) {
        try {
            const response = await api.vector.searchSimilarTrips({
                destination,
                preferences,
                limit
            })
            return response.data
        } catch (error) {
            console.error('Error searching similar trips:', error)
            throw new Error('Failed to search similar trips. Please try again.')
        }
    }

    /**
     * Store user preferences in vector database
     */
    async storeUserPreferences(userId: string, preferences: Record<string, any>, tripHistory?: string[]) {
        try {
            const response = await api.vector.storeUserPreferences({
                user_id: userId,
                preferences,
                trip_history: tripHistory
            })
            return response.data
        } catch (error) {
            console.error('Error storing user preferences:', error)
            throw new Error('Failed to store preferences. Please try again.')
        }
    }

    /**
     * Validate real-time availability of items
     */
    async validateAvailability(type: 'attractions' | 'hotels' | 'transportation', data: any[], checkDate: string) {
        try {
            const response = await api.ai.validateAvailability({
                type,
                data,
                check_date: checkDate
            })
            return response.data
        } catch (error) {
            console.error('Error validating availability:', error)
            throw new Error('Failed to validate availability. Please try again.')
        }
    }

    /**
     * Get AI recommendations with RAG context
     */
    async getRecommendations(userId?: string, budget?: number, interests: string[] = []) {
        try {
            const response = await api.ai.getRecommendations({
                user_id: userId,
                budget,
                interests
            })
            return response.data
        } catch (error) {
            console.error('Error getting recommendations:', error)
            throw new Error('Failed to get recommendations. Please try again.')
        }
    }

    /**
     * Parse and validate RAG itinerary response
     */
    parseItineraryResponse(response: RAGItineraryResponse) {
        const parsedResponse = {
            ...response,
            dayPlans: this.extractDayPlans(response),
            hasWeatherInfo: !!response.weather_info,
            hasHotelRecommendations: !!response.recommended_hotels && response.recommended_hotels.length > 0,
            hasTransportation: !!response.transportation && response.transportation.length > 0,
            hasLocalEvents: !!response.local_events && response.local_events.length > 0,
            hasEMTServices: !!response.emt_services && response.emt_services.length > 0,
            isRAGEnhanced: response.rag_enhanced === true,
            dataSources: response.data_sources || []
        }

        return parsedResponse
    }

    /**
     * Extract day plans from response
     */
    private extractDayPlans(response: RAGItineraryResponse) {
        const dayPlans = []
        let dayIndex = 1

        while (response[`day_${dayIndex}`]) {
            const dayData = response[`day_${dayIndex}`]
            dayPlans.push({
                day: dayIndex,
                date: this.calculateDate(response.created_at, dayIndex - 1),
                morning: dayData.morning || '',
                afternoon: dayData.afternoon || '',
                evening: dayData.evening || '',
                activities: this.parseActivities(dayData),
                recommendations: this.parseDayRecommendations(dayData)
            })
            dayIndex++
        }

        return dayPlans
    }

    /**
     * Parse activities from day data
     */
    private parseActivities(dayData: any) {
        const activities = []

        if (dayData.morning) {
            activities.push({ time: 'morning', description: dayData.morning, type: 'activity' })
        }

        if (dayData.afternoon) {
            activities.push({ time: 'afternoon', description: dayData.afternoon, type: 'activity' })
        }

        if (dayData.evening) {
            activities.push({ time: 'evening', description: dayData.evening, type: 'activity' })
        }

        return activities
    }

    /**
     * Parse recommendations for a specific day
     */
    private parseDayRecommendations(dayData: any) {
        const recommendations = []

        if (dayData.restaurant_recommendations) {
            recommendations.push(...dayData.restaurant_recommendations.map((r: any) => ({ ...r, type: 'restaurant' })))
        }

        if (dayData.activity_recommendations) {
            recommendations.push(...dayData.activity_recommendations.map((r: any) => ({ ...r, type: 'activity' })))
        }

        return recommendations
    }

    /**
     * Calculate date for a specific day
     */
    private calculateDate(createdAt: string, daysToAdd: number): string {
        const baseDate = new Date(createdAt)
        const targetDate = new Date(baseDate.getTime() + (daysToAdd * 24 * 60 * 60 * 1000))
        return targetDate.toISOString().split('T')[0]
    }

    /**
     * Format budget for display
     */
    formatBudget(budget: number, currency: string = 'USD'): string {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: currency
        }).format(budget)
    }

    /**
     * Calculate trip duration in days
     */
    calculateDuration(startDate: string, endDate: string): number {
        const start = new Date(startDate)
        const end = new Date(endDate)
        const diffTime = Math.abs(end.getTime() - start.getTime())
        return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    }

    /**
     * Validate trip request data
     */
    validateTripRequest(request: RAGTripRequest): { isValid: boolean; errors: string[] } {
        const errors: string[] = []

        if (!request.destination?.trim()) {
            errors.push('Destination is required')
        }

        if (!request.start_date) {
            errors.push('Start date is required')
        }

        if (!request.end_date) {
            errors.push('End date is required')
        }

        if (request.start_date && request.end_date) {
            const startDate = new Date(request.start_date)
            const endDate = new Date(request.end_date)

            if (startDate >= endDate) {
                errors.push('End date must be after start date')
            }

            if (startDate < new Date()) {
                errors.push('Start date cannot be in the past')
            }
        }

        if (!request.budget || request.budget <= 0) {
            errors.push('Budget must be greater than 0')
        }

        if (!request.travelers || request.travelers <= 0) {
            errors.push('Number of travelers must be greater than 0')
        }

        if (!request.user_id?.trim()) {
            errors.push('User ID is required')
        }

        return {
            isValid: errors.length === 0,
            errors
        }
    }
}

export const ragService = new RAGService()
export default ragService