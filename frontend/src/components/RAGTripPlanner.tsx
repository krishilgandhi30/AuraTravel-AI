'use client'

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import {
    Brain,
    Database,
    MapPin,
    Cloud,
    Calendar,
    Users,
    DollarSign,
    Sparkles,
    CheckCircle,
    Clock,
    TrendingUp,
    Zap,
    RefreshCw,
    AlertCircle,
    Info
} from 'lucide-react'
import { ragService, RAGTripRequest } from '@/lib/api/rag-service'
import { RAGItineraryResponse, RAGTripContext } from '@/types'

interface RAGTripPlannerProps {
    initialData?: Partial<RAGTripRequest>
    onSuccess?: (itinerary: RAGItineraryResponse) => void
    onError?: (error: string) => void
}

interface RAGStep {
    id: string
    title: string
    description: string
    icon: React.ReactNode
    status: 'pending' | 'loading' | 'completed' | 'error'
    data?: any
}

export default function RAGTripPlanner({ initialData, onSuccess, onError }: RAGTripPlannerProps) {
    const [isPlanning, setIsPlanning] = useState(false)
    const [currentStep, setCurrentStep] = useState(0)
    const [ragContext, setRagContext] = useState<RAGTripContext | null>(null)
    const [finalItinerary, setFinalItinerary] = useState<RAGItineraryResponse | null>(null)
    const [error, setError] = useState<string | null>(null)

    const [ragSteps, setRagSteps] = useState<RAGStep[]>([
        {
            id: 'context',
            title: 'Gathering Context',
            description: 'Retrieving real-time data about your destination',
            icon: <Database className="w-5 h-5" />,
            status: 'pending'
        },
        {
            id: 'attractions',
            title: 'Finding Attractions',
            description: 'Discovering top-rated attractions and activities',
            icon: <MapPin className="w-5 h-5" />,
            status: 'pending'
        },
        {
            id: 'weather',
            title: 'Weather Analysis',
            description: 'Checking weather forecasts for optimal planning',
            icon: <Cloud className="w-5 h-5" />,
            status: 'pending'
        },
        {
            id: 'validation',
            title: 'Validating Availability',
            description: 'Verifying real-time availability and pricing',
            icon: <CheckCircle className="w-5 h-5" />,
            status: 'pending'
        },
        {
            id: 'generation',
            title: 'AI Generation',
            description: 'Creating your personalized itinerary',
            icon: <Brain className="w-5 h-5" />,
            status: 'pending'
        }
    ])

    const updateStepStatus = (stepId: string, status: RAGStep['status'], data?: any) => {
        setRagSteps(prev => prev.map(step =>
            step.id === stepId
                ? { ...step, status, data }
                : step
        ))
    }

    const executeRAGPlanning = async (request: RAGTripRequest) => {
        setIsPlanning(true)
        setError(null)
        setCurrentStep(0)

        try {
            // Step 1: Get RAG Context
            updateStepStatus('context', 'loading')
            setCurrentStep(0)

            const contextRequest = {
                destination: request.destination,
                user_id: request.user_id,
                start_date: request.start_date,
                end_date: request.end_date,
                budget: request.budget,
                travelers: request.travelers,
                interests: request.interests,
                preferences: request.preferences
            }

            const context = await ragService.getRAGContext(contextRequest)
            setRagContext(context)
            updateStepStatus('context', 'completed', context)

            // Step 2: Process Attractions
            updateStepStatus('attractions', 'loading')
            setCurrentStep(1)

            await new Promise(resolve => setTimeout(resolve, 1000)) // Simulate processing
            updateStepStatus('attractions', 'completed', context.attractions)

            // Step 3: Weather Analysis
            updateStepStatus('weather', 'loading')
            setCurrentStep(2)

            await new Promise(resolve => setTimeout(resolve, 800))
            updateStepStatus('weather', 'completed', context.weather)

            // Step 4: Validation
            updateStepStatus('validation', 'loading')
            setCurrentStep(3)

            if (context.hotels.length > 0) {
                await ragService.validateAvailability('hotels', context.hotels, request.start_date)
            }
            if (context.attractions.length > 0) {
                await ragService.validateAvailability('attractions', context.attractions, request.start_date)
            }

            updateStepStatus('validation', 'completed')

            // Step 5: AI Generation
            updateStepStatus('generation', 'loading')
            setCurrentStep(4)

            const itinerary = await ragService.planTrip(request)
            setFinalItinerary(itinerary)
            updateStepStatus('generation', 'completed', itinerary)

            onSuccess?.(itinerary)

        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : 'An unexpected error occurred'
            setError(errorMessage)
            updateStepStatus(ragSteps[currentStep]?.id || 'unknown', 'error')
            onError?.(errorMessage)
        } finally {
            setIsPlanning(false)
        }
    }

    const formatContextStats = (context: RAGTripContext | null) => {
        if (!context) return null

        return {
            attractions: context.attractions?.length || 0,
            hotels: context.hotels?.length || 0,
            events: context.localEvents?.length || 0,
            transportOptions: context.transportation?.length || 0,
            weatherDays: context.weather?.forecast?.length || 0
        }
    }

    const getStepProgress = () => {
        const completedSteps = ragSteps.filter(step => step.status === 'completed').length
        return (completedSteps / ragSteps.length) * 100
    }

    return (
        <div className="w-full max-w-4xl mx-auto">
            {/* RAG Process Header */}
            <div className="bg-gradient-to-r from-blue-600 to-purple-600 text-white p-6 rounded-t-lg">
                <div className="flex items-center space-x-3 mb-4">
                    <Zap className="w-8 h-8" />
                    <h2 className="text-2xl font-bold">RAG-Enhanced Trip Planning</h2>
                </div>
                <p className="text-blue-100 mb-4">
                    Using Retrieval-Augmented Generation for real-time, accurate travel recommendations
                </p>

                {/* Progress Bar */}
                <div className="w-full bg-blue-800 rounded-full h-2 mb-2">
                    <motion.div
                        className="bg-yellow-400 h-2 rounded-full"
                        initial={{ width: 0 }}
                        animate={{ width: `${getStepProgress()}%` }}
                        transition={{ duration: 0.5 }}
                    />
                </div>
                <div className="text-sm text-blue-100">
                    {Math.round(getStepProgress())}% Complete
                </div>
            </div>

            {/* RAG Steps */}
            <div className="bg-white border-x border-gray-200 p-6">
                <div className="space-y-4">
                    {ragSteps.map((step, index) => (
                        <motion.div
                            key={step.id}
                            className={`flex items-center space-x-4 p-4 rounded-lg border-2 transition-all duration-300 ${step.status === 'completed'
                                    ? 'border-green-200 bg-green-50'
                                    : step.status === 'loading'
                                        ? 'border-blue-200 bg-blue-50'
                                        : step.status === 'error'
                                            ? 'border-red-200 bg-red-50'
                                            : 'border-gray-200 bg-gray-50'
                                }`}
                            initial={{ opacity: 0, x: -20 }}
                            animate={{ opacity: 1, x: 0 }}
                            transition={{ delay: index * 0.1 }}
                        >
                            <div className={`flex-shrink-0 w-12 h-12 rounded-full flex items-center justify-center ${step.status === 'completed'
                                    ? 'bg-green-500 text-white'
                                    : step.status === 'loading'
                                        ? 'bg-blue-500 text-white'
                                        : step.status === 'error'
                                            ? 'bg-red-500 text-white'
                                            : 'bg-gray-300 text-gray-600'
                                }`}>
                                {step.status === 'loading' ? (
                                    <RefreshCw className="w-5 h-5 animate-spin" />
                                ) : step.status === 'error' ? (
                                    <AlertCircle className="w-5 h-5" />
                                ) : step.status === 'completed' ? (
                                    <CheckCircle className="w-5 h-5" />
                                ) : (
                                    step.icon
                                )}
                            </div>

                            <div className="flex-1">
                                <h3 className="font-semibold text-gray-900">{step.title}</h3>
                                <p className="text-sm text-gray-600">{step.description}</p>

                                {/* Step-specific data display */}
                                {step.status === 'completed' && step.id === 'context' && step.data && (
                                    <div className="mt-2 flex flex-wrap gap-2">
                                        <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">
                                            {step.data.attractions?.length || 0} attractions
                                        </span>
                                        <span className="text-xs bg-green-100 text-green-800 px-2 py-1 rounded">
                                            {step.data.hotels?.length || 0} hotels
                                        </span>
                                        <span className="text-xs bg-purple-100 text-purple-800 px-2 py-1 rounded">
                                            {step.data.localEvents?.length || 0} events
                                        </span>
                                    </div>
                                )}

                                {step.status === 'completed' && step.id === 'weather' && step.data && (
                                    <div className="mt-2">
                                        <span className="text-xs bg-yellow-100 text-yellow-800 px-2 py-1 rounded">
                                            {step.data.forecast?.length || 0} day forecast available
                                        </span>
                                    </div>
                                )}
                            </div>

                            <div className="flex-shrink-0">
                                {step.status === 'completed' && (
                                    <div className="text-green-600 text-sm font-medium">
                                        ✓ Complete
                                    </div>
                                )}
                                {step.status === 'loading' && (
                                    <div className="text-blue-600 text-sm font-medium">
                                        Processing...
                                    </div>
                                )}
                                {step.status === 'error' && (
                                    <div className="text-red-600 text-sm font-medium">
                                        Failed
                                    </div>
                                )}
                            </div>
                        </motion.div>
                    ))}
                </div>
            </div>

            {/* Error Display */}
            {error && (
                <motion.div
                    className="bg-red-50 border border-red-200 text-red-700 p-4 mx-6 rounded-lg"
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                >
                    <div className="flex items-center space-x-2">
                        <AlertCircle className="w-5 h-5" />
                        <p className="font-semibold">Planning Failed</p>
                    </div>
                    <p className="text-sm mt-1">{error}</p>
                </motion.div>
            )}

            {/* RAG Context Summary */}
            {ragContext && (
                <motion.div
                    className="bg-gray-50 border-x border-gray-200 p-6"
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                >
                    <h3 className="font-semibold text-gray-900 mb-4 flex items-center">
                        <Database className="w-5 h-5 mr-2" />
                        Retrieved Context Summary
                    </h3>

                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                        <div className="text-center">
                            <div className="text-2xl font-bold text-blue-600">
                                {ragContext.attractions?.length || 0}
                            </div>
                            <div className="text-sm text-gray-600">Attractions</div>
                        </div>
                        <div className="text-center">
                            <div className="text-2xl font-bold text-green-600">
                                {ragContext.hotels?.length || 0}
                            </div>
                            <div className="text-sm text-gray-600">Hotels</div>
                        </div>
                        <div className="text-center">
                            <div className="text-2xl font-bold text-purple-600">
                                {ragContext.localEvents?.length || 0}
                            </div>
                            <div className="text-sm text-gray-600">Local Events</div>
                        </div>
                        <div className="text-center">
                            <div className="text-2xl font-bold text-orange-600">
                                {ragContext.weather?.forecast?.length || 0}
                            </div>
                            <div className="text-sm text-gray-600">Weather Days</div>
                        </div>
                    </div>
                </motion.div>
            )}

            {/* Final Itinerary Preview */}
            {finalItinerary && (
                <motion.div
                    className="bg-white border-x border-gray-200 p-6"
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                >
                    <div className="flex items-center justify-between mb-4">
                        <h3 className="font-semibold text-gray-900 flex items-center">
                            <Sparkles className="w-5 h-5 mr-2 text-yellow-500" />
                            RAG-Enhanced Itinerary Generated
                        </h3>
                        {finalItinerary.rag_enhanced && (
                            <span className="bg-gradient-to-r from-blue-500 to-purple-600 text-white px-3 py-1 rounded-full text-sm font-medium">
                                RAG Enhanced
                            </span>
                        )}
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <div className="bg-blue-50 p-4 rounded-lg">
                            <div className="text-lg font-bold text-blue-600">
                                {finalItinerary.destination}
                            </div>
                            <div className="text-sm text-blue-800">
                                {finalItinerary.duration} days • {finalItinerary.travelers} travelers
                            </div>
                        </div>

                        <div className="bg-green-50 p-4 rounded-lg">
                            <div className="text-lg font-bold text-green-600">
                                ${finalItinerary.budget?.toLocaleString()}
                            </div>
                            <div className="text-sm text-green-800">Total Budget</div>
                        </div>

                        <div className="bg-purple-50 p-4 rounded-lg">
                            <div className="text-lg font-bold text-purple-600">
                                {finalItinerary.data_sources?.length || 0}
                            </div>
                            <div className="text-sm text-purple-800">Data Sources</div>
                        </div>
                    </div>

                    {finalItinerary.tips && finalItinerary.tips.length > 0 && (
                        <div className="mt-4">
                            <h4 className="font-medium text-gray-900 mb-2">AI-Generated Tips</h4>
                            <div className="space-y-1">
                                {finalItinerary.tips.slice(0, 3).map((tip, index) => (
                                    <div key={index} className="flex items-center space-x-2 text-sm text-gray-600">
                                        <Info className="w-4 h-4 text-blue-500" />
                                        <span>{tip}</span>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}
                </motion.div>
            )}

            {/* Action Buttons */}
            <div className="bg-gray-50 border border-gray-200 rounded-b-lg p-6">
                <div className="flex justify-center">
                    {!isPlanning && !finalItinerary && initialData && (
                        <button
                            onClick={() => executeRAGPlanning(initialData as RAGTripRequest)}
                            className="bg-gradient-to-r from-blue-600 to-purple-600 text-white px-8 py-3 rounded-lg font-semibold hover:from-blue-700 hover:to-purple-700 transition-all duration-200 flex items-center space-x-2"
                        >
                            <Brain className="w-5 h-5" />
                            <span>Start RAG Planning</span>
                        </button>
                    )}

                    {finalItinerary && (
                        <div className="flex space-x-4">
                            <button
                                onClick={() => {
                                    setFinalItinerary(null)
                                    setRagContext(null)
                                    setRagSteps(prev => prev.map(step => ({ ...step, status: 'pending' })))
                                    setCurrentStep(0)
                                }}
                                className="bg-gray-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-gray-700 transition-colors"
                            >
                                Plan Another Trip
                            </button>
                            <button
                                onClick={() => onSuccess?.(finalItinerary)}
                                className="bg-gradient-to-r from-green-600 to-blue-600 text-white px-6 py-3 rounded-lg font-semibold hover:from-green-700 hover:to-blue-700 transition-all duration-200"
                            >
                                Use This Itinerary
                            </button>
                        </div>
                    )}
                </div>
            </div>
        </div>
    )
}

export { RAGTripPlanner }