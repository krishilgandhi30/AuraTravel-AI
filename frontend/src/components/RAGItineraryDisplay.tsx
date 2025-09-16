'use client'

import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import {
    MapPin,
    Clock,
    Star,
    DollarSign,
    Calendar,
    Cloud,
    Sun,
    CloudRain,
    Thermometer,
    Wind,
    Users,
    Plane,
    Car,
    Train,
    Bus,
    Hotel,
    Utensils,
    Camera,
    Music,
    Shield,
    ExternalLink,
    ChevronDown,
    ChevronUp,
    Info,
    CheckCircle,
    AlertTriangle,
    Sparkles,
    Database
} from 'lucide-react'
import { RAGItineraryResponse, WeatherCondition, RAGHotel, TransportOption, LocalEvent, EMTItem } from '@/types'
import { ragService } from '@/lib/api/rag-service'

interface RAGItineraryDisplayProps {
    itinerary: RAGItineraryResponse
    onBooking?: (type: string, item: any) => void
    onModify?: () => void
}

export default function RAGItineraryDisplay({ itinerary, onBooking, onModify }: RAGItineraryDisplayProps) {
    const [activeDay, setActiveDay] = useState(1)
    const [showWeatherDetails, setShowWeatherDetails] = useState(false)
    const [showDataSources, setShowDataSources] = useState(false)

    const parsedItinerary = ragService.parseItineraryResponse(itinerary)

    const getWeatherIcon = (description: string) => {
        const desc = description.toLowerCase()
        if (desc.includes('rain')) return <CloudRain className="w-5 h-5" />
        if (desc.includes('cloud')) return <Cloud className="w-5 h-5" />
        if (desc.includes('sun') || desc.includes('clear')) return <Sun className="w-5 h-5" />
        return <Cloud className="w-5 h-5" />
    }

    const getTransportIcon = (type: string) => {
        switch (type.toLowerCase()) {
            case 'flight': return <Plane className="w-4 h-4" />
            case 'car': case 'taxi': return <Car className="w-4 h-4" />
            case 'train': return <Train className="w-4 h-4" />
            case 'bus': return <Bus className="w-4 h-4" />
            default: return <Car className="w-4 h-4" />
        }
    }

    const formatCurrency = (amount: number) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD'
        }).format(amount)
    }

    return (
        <div className="max-w-6xl mx-auto bg-white rounded-lg shadow-lg overflow-hidden">
            {/* Header */}
            <div className="bg-gradient-to-r from-blue-600 to-purple-600 text-white p-6">
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-3xl font-bold mb-2">{itinerary.destination}</h1>
                        <div className="flex items-center space-x-4 text-blue-100">
                            <div className="flex items-center space-x-1">
                                <Calendar className="w-4 h-4" />
                                <span>{itinerary.duration} days</span>
                            </div>
                            <div className="flex items-center space-x-1">
                                <Users className="w-4 h-4" />
                                <span>{itinerary.travelers} travelers</span>
                            </div>
                            <div className="flex items-center space-x-1">
                                <DollarSign className="w-4 h-4" />
                                <span>{formatCurrency(itinerary.budget)}</span>
                            </div>
                        </div>
                    </div>

                    <div className="flex items-center space-x-2">
                        {itinerary.rag_enhanced && (
                            <span className="bg-yellow-400 text-yellow-900 px-3 py-1 rounded-full text-sm font-medium flex items-center space-x-1">
                                <Sparkles className="w-4 h-4" />
                                <span>RAG Enhanced</span>
                            </span>
                        )}
                        {itinerary.ai_generated && (
                            <span className="bg-white bg-opacity-20 px-3 py-1 rounded-full text-sm font-medium">
                                AI Generated
                            </span>
                        )}
                    </div>
                </div>

                {/* Data Sources */}
                {itinerary.data_sources && itinerary.data_sources.length > 0 && (
                    <motion.div className="mt-4">
                        <button
                            onClick={() => setShowDataSources(!showDataSources)}
                            className="flex items-center space-x-2 text-blue-100 hover:text-white transition-colors"
                        >
                            <Database className="w-4 h-4" />
                            <span>Data Sources ({itinerary.data_sources.length})</span>
                            {showDataSources ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
                        </button>

                        <AnimatePresence>
                            {showDataSources && (
                                <motion.div
                                    initial={{ opacity: 0, height: 0 }}
                                    animate={{ opacity: 1, height: 'auto' }}
                                    exit={{ opacity: 0, height: 0 }}
                                    className="mt-2 flex flex-wrap gap-2"
                                >
                                    {itinerary.data_sources.map((source, index) => (
                                        <span
                                            key={index}
                                            className="bg-white bg-opacity-20 px-2 py-1 rounded text-xs"
                                        >
                                            {source.replace(/_/g, ' ')}
                                        </span>
                                    ))}
                                </motion.div>
                            )}
                        </AnimatePresence>
                    </motion.div>
                )}
            </div>

            {/* Weather Information */}
            {itinerary.weather_info && (
                <div className="bg-blue-50 border-b p-4">
                    <button
                        onClick={() => setShowWeatherDetails(!showWeatherDetails)}
                        className="flex items-center space-x-2 text-blue-800 hover:text-blue-900 w-full"
                    >
                        <div className="flex items-center space-x-2">
                            {getWeatherIcon(itinerary.weather_info.current?.description || '')}
                            <span className="font-medium">Weather Forecast</span>
                            <span className="text-sm text-blue-600">
                                {itinerary.weather_info.current?.temperature}°C
                            </span>
                        </div>
                        <div className="flex-1" />
                        {showWeatherDetails ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
                    </button>

                    <AnimatePresence>
                        {showWeatherDetails && itinerary.weather_info.forecast && (
                            <motion.div
                                initial={{ opacity: 0, height: 0 }}
                                animate={{ opacity: 1, height: 'auto' }}
                                exit={{ opacity: 0, height: 0 }}
                                className="mt-4"
                            >
                                <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-7 gap-3">
                                    {itinerary.weather_info.forecast.slice(0, 7).map((day, index) => (
                                        <div key={index} className="bg-white p-3 rounded-lg text-center">
                                            <div className="text-xs text-gray-600 mb-1">
                                                {new Date(day.date).toLocaleDateString('en-US', { weekday: 'short' })}
                                            </div>
                                            <div className="flex justify-center mb-1">
                                                {getWeatherIcon(day.description)}
                                            </div>
                                            <div className="text-sm font-medium">{day.temperature}°C</div>
                                            <div className="text-xs text-gray-500">{day.humidity}%</div>
                                        </div>
                                    ))}
                                </div>

                                {itinerary.weather_info.tips && (
                                    <div className="mt-4 bg-yellow-50 border border-yellow-200 rounded-lg p-3">
                                        <h4 className="font-medium text-yellow-800 mb-2">Weather Tips</h4>
                                        <ul className="space-y-1">
                                            {itinerary.weather_info.tips.map((tip, index) => (
                                                <li key={index} className="text-sm text-yellow-700 flex items-center space-x-2">
                                                    <Info className="w-3 h-3" />
                                                    <span>{tip}</span>
                                                </li>
                                            ))}
                                        </ul>
                                    </div>
                                )}
                            </motion.div>
                        )}
                    </AnimatePresence>
                </div>
            )}

            <div className="flex">
                {/* Day Navigation */}
                <div className="w-48 bg-gray-50 border-r">
                    <div className="p-4">
                        <h3 className="font-semibold text-gray-900 mb-3">Itinerary</h3>
                        <div className="space-y-2">
                            {parsedItinerary.dayPlans.map((day) => (
                                <button
                                    key={day.day}
                                    onClick={() => setActiveDay(day.day)}
                                    className={`w-full text-left p-3 rounded-lg transition-colors ${activeDay === day.day
                                            ? 'bg-blue-100 text-blue-800 border border-blue-200'
                                            : 'bg-white text-gray-700 hover:bg-gray-100'
                                        }`}
                                >
                                    <div className="font-medium">Day {day.day}</div>
                                    <div className="text-sm text-gray-500">
                                        {new Date(day.date).toLocaleDateString('en-US', {
                                            month: 'short',
                                            day: 'numeric'
                                        })}
                                    </div>
                                </button>
                            ))}
                        </div>
                    </div>
                </div>

                {/* Day Content */}
                <div className="flex-1 p-6">
                    {parsedItinerary.dayPlans.map((day) => (
                        <AnimatePresence key={day.day}>
                            {activeDay === day.day && (
                                <motion.div
                                    initial={{ opacity: 0, x: 20 }}
                                    animate={{ opacity: 1, x: 0 }}
                                    exit={{ opacity: 0, x: -20 }}
                                    transition={{ duration: 0.3 }}
                                >
                                    <h2 className="text-2xl font-bold text-gray-900 mb-4">
                                        Day {day.day} - {new Date(day.date).toLocaleDateString('en-US', {
                                            weekday: 'long',
                                            month: 'long',
                                            day: 'numeric'
                                        })}
                                    </h2>

                                    <div className="space-y-6">
                                        {day.activities.map((activity, index) => (
                                            <div key={index} className="bg-white border rounded-lg p-4 shadow-sm">
                                                <div className="flex items-center space-x-3 mb-3">
                                                    <div className="bg-blue-100 text-blue-600 px-3 py-1 rounded-full text-sm font-medium capitalize">
                                                        {activity.time}
                                                    </div>
                                                    <Clock className="w-4 h-4 text-gray-400" />
                                                </div>
                                                <p className="text-gray-800">{activity.description}</p>
                                            </div>
                                        ))}
                                    </div>
                                </motion.div>
                            )}
                        </AnimatePresence>
                    ))}
                </div>
            </div>

            {/* Hotel Recommendations */}
            {itinerary.recommended_hotels && itinerary.recommended_hotels.length > 0 && (
                <div className="bg-gray-50 border-t p-6">
                    <h3 className="font-semibold text-gray-900 mb-4 flex items-center">
                        <Hotel className="w-5 h-5 mr-2" />
                        Recommended Hotels
                    </h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                        {itinerary.recommended_hotels.map((hotel, index) => (
                            <div key={index} className="bg-white border rounded-lg p-4">
                                <div className="flex items-center justify-between mb-2">
                                    <h4 className="font-medium">{hotel.name}</h4>
                                    <div className="flex items-center space-x-1">
                                        <Star className="w-4 h-4 text-yellow-500" />
                                        <span className="text-sm">{hotel.rating}</span>
                                    </div>
                                </div>
                                <div className="text-sm text-gray-600 mb-2">
                                    {formatCurrency(hotel.price_per_night)}/night
                                </div>
                                <div className="flex items-center justify-between">
                                    <span className={`text-xs px-2 py-1 rounded-full ${hotel.availability ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                                        }`}>
                                        {hotel.availability ? 'Available' : 'Not Available'}
                                    </span>
                                    {hotel.booking_url && (
                                        <button
                                            onClick={() => onBooking?.('hotel', hotel)}
                                            className="text-blue-600 hover:text-blue-800 text-sm flex items-center space-x-1"
                                        >
                                            <span>Book</span>
                                            <ExternalLink className="w-3 h-3" />
                                        </button>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}

            {/* Transportation */}
            {itinerary.transportation && itinerary.transportation.length > 0 && (
                <div className="bg-white border-t p-6">
                    <h3 className="font-semibold text-gray-900 mb-4 flex items-center">
                        <Plane className="w-5 h-5 mr-2" />
                        Transportation Options
                    </h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        {itinerary.transportation.map((transport, index) => (
                            <div key={index} className="bg-gray-50 border rounded-lg p-4">
                                <div className="flex items-center space-x-3 mb-2">
                                    {getTransportIcon(transport.type)}
                                    <span className="font-medium capitalize">{transport.type}</span>
                                    <span className="text-sm text-gray-600">{transport.provider}</span>
                                </div>
                                <div className="text-sm text-gray-600 mb-2">
                                    {formatCurrency(transport.price)}
                                </div>
                                <div className="flex items-center justify-between">
                                    <span className={`text-xs px-2 py-1 rounded-full ${transport.availability ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                                        }`}>
                                        {transport.availability ? 'Available' : 'Not Available'}
                                    </span>
                                    {transport.booking_url && (
                                        <button
                                            onClick={() => onBooking?.('transport', transport)}
                                            className="text-blue-600 hover:text-blue-800 text-sm flex items-center space-x-1"
                                        >
                                            <span>Book</span>
                                            <ExternalLink className="w-3 h-3" />
                                        </button>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}

            {/* Local Events */}
            {itinerary.local_events && itinerary.local_events.length > 0 && (
                <div className="bg-gray-50 border-t p-6">
                    <h3 className="font-semibold text-gray-900 mb-4 flex items-center">
                        <Music className="w-5 h-5 mr-2" />
                        Local Events
                    </h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        {itinerary.local_events.map((event, index) => (
                            <div key={index} className="bg-white border rounded-lg p-4">
                                <h4 className="font-medium mb-1">{event.name}</h4>
                                <div className="text-sm text-gray-600 mb-2">
                                    {new Date(event.date).toLocaleDateString('en-US', {
                                        month: 'short',
                                        day: 'numeric',
                                        year: 'numeric'
                                    })}
                                </div>
                                <p className="text-sm text-gray-700 mb-2">{event.description}</p>
                                <div className="flex items-center justify-between">
                                    <span className="text-sm font-medium text-green-600">
                                        {event.price === 0 ? 'Free' : formatCurrency(event.price)}
                                    </span>
                                    <button className="text-blue-600 hover:text-blue-800 text-sm">
                                        Learn More
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}

            {/* EMT Services */}
            {itinerary.emt_services && itinerary.emt_services.length > 0 && (
                <div className="bg-red-50 border-t p-6">
                    <h3 className="font-semibold text-red-900 mb-4 flex items-center">
                        <Shield className="w-5 h-5 mr-2" />
                        Emergency Medical Services
                    </h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        {itinerary.emt_services.map((service, index) => (
                            <div key={index} className="bg-white border border-red-200 rounded-lg p-4">
                                <h4 className="font-medium text-red-900 mb-1">{service.name}</h4>
                                <div className="text-sm text-red-700 mb-2 capitalize">{service.type}</div>
                                <p className="text-sm text-gray-700 mb-2">{service.description}</p>
                                <div className="flex items-center justify-between">
                                    <span className={`text-xs px-2 py-1 rounded-full ${service.available ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                                        }`}>
                                        {service.available ? 'Available' : 'Not Available'}
                                    </span>
                                    <span className="text-sm text-red-600 font-medium">{service.contact}</span>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}

            {/* Tips */}
            {itinerary.tips && itinerary.tips.length > 0 && (
                <div className="bg-yellow-50 border-t p-6">
                    <h3 className="font-semibold text-yellow-900 mb-4 flex items-center">
                        <Info className="w-5 h-5 mr-2" />
                        AI-Generated Travel Tips
                    </h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        {itinerary.tips.map((tip, index) => (
                            <div key={index} className="flex items-start space-x-3 bg-white border border-yellow-200 rounded-lg p-3">
                                <CheckCircle className="w-5 h-5 text-yellow-600 mt-0.5 flex-shrink-0" />
                                <span className="text-sm text-gray-700">{tip}</span>
                            </div>
                        ))}
                    </div>
                </div>
            )}

            {/* Action Buttons */}
            <div className="bg-white border-t p-6">
                <div className="flex justify-center space-x-4">
                    <button
                        onClick={onModify}
                        className="bg-gray-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-gray-700 transition-colors"
                    >
                        Modify Trip
                    </button>
                    <button
                        onClick={() => window.print()}
                        className="bg-blue-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-blue-700 transition-colors"
                    >
                        Save Itinerary
                    </button>
                </div>
            </div>
        </div>
    )
}

export { RAGItineraryDisplay }