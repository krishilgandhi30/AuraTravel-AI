'use client'

import { useState } from 'react'
import Link from 'next/link'
import { motion } from 'framer-motion'
import {
    MapPin,
    Calendar,
    Users,
    DollarSign,
    Plane,
    ArrowRight,
    Camera,
    Utensils,
    Mountain,
    Building,
    Sun,
    Snowflake,
    Heart,
    ChevronDown,
    Plus,
    X
} from 'lucide-react'

interface TripPreferences {
    destination: string
    startDate: string
    endDate: string
    travelers: number
    budget: string
    travelStyle: string[]
    interests: string[]
    accommodation: string
    transportation: string
}

export default function PlanPage() {
    const [step, setStep] = useState(1)
    const [preferences, setPreferences] = useState<TripPreferences>({
        destination: '',
        startDate: '',
        endDate: '',
        travelers: 1,
        budget: '',
        travelStyle: [],
        interests: [],
        accommodation: '',
        transportation: ''
    })

    const budgetOptions = [
        { value: 'budget', label: 'Budget ($)', range: 'Under $1,000' },
        { value: 'mid-range', label: 'Mid-range ($$)', range: '$1,000 - $3,000' },
        { value: 'luxury', label: 'Luxury ($$$)', range: '$3,000 - $7,000' },
        { value: 'ultra-luxury', label: 'Ultra Luxury ($$$$)', range: '$7,000+' }
    ]

    const travelStyles = [
        { id: 'adventure', label: 'Adventure', icon: <Mountain className="w-5 h-5" /> },
        { id: 'relaxation', label: 'Relaxation', icon: <Sun className="w-5 h-5" /> },
        { id: 'cultural', label: 'Cultural', icon: <Building className="w-5 h-5" /> },
        { id: 'foodie', label: 'Foodie', icon: <Utensils className="w-5 h-5" /> },
        { id: 'photography', label: 'Photography', icon: <Camera className="w-5 h-5" /> },
        { id: 'romantic', label: 'Romantic', icon: <Heart className="w-5 h-5" /> }
    ]

    const interests = [
        'Beaches', 'Mountains', 'Cities', 'History', 'Art & Museums',
        'Nightlife', 'Shopping', 'Nature & Wildlife', 'Local Cuisine',
        'Festivals & Events', 'Architecture', 'Adventure Sports'
    ]

    const accommodationTypes = [
        { value: 'hotel', label: 'Hotel' },
        { value: 'resort', label: 'Resort' },
        { value: 'apartment', label: 'Apartment/Airbnb' },
        { value: 'hostel', label: 'Hostel' },
        { value: 'villa', label: 'Villa' },
        { value: 'boutique', label: 'Boutique Hotel' }
    ]

    const transportationOptions = [
        { value: 'flight', label: 'Flight' },
        { value: 'car', label: 'Car/Road Trip' },
        { value: 'train', label: 'Train' },
        { value: 'bus', label: 'Bus' },
        { value: 'cruise', label: 'Cruise' }
    ]

    const handleInputChange = (field: keyof TripPreferences, value: any) => {
        setPreferences(prev => ({
            ...prev,
            [field]: value
        }))
    }

    const handleArrayToggle = (field: 'travelStyle' | 'interests', value: string) => {
        setPreferences(prev => ({
            ...prev,
            [field]: prev[field].includes(value)
                ? prev[field].filter(item => item !== value)
                : [...prev[field], value]
        }))
    }

    const handleSubmit = async () => {
        console.log('Submitting trip preferences:', preferences)
        // Handle trip planning submission
    }

    const renderStep1 = () => (
        <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.5 }}
            className="space-y-6"
        >
            <div>
                <h2 className="text-2xl font-bold text-gray-900 mb-6">Where would you like to go?</h2>

                {/* Destination Input */}
                <div className="space-y-4">
                    <div className="relative">
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Destination
                        </label>
                        <div className="relative">
                            <MapPin className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
                            <input
                                type="text"
                                value={preferences.destination}
                                onChange={(e) => handleInputChange('destination', e.target.value)}
                                placeholder="e.g., Paris, Tokyo, New York..."
                                className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                            />
                        </div>
                    </div>

                    {/* Date Selection */}
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">
                                Start Date
                            </label>
                            <div className="relative">
                                <Calendar className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
                                <input
                                    type="date"
                                    value={preferences.startDate}
                                    onChange={(e) => handleInputChange('startDate', e.target.value)}
                                    className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                                />
                            </div>
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">
                                End Date
                            </label>
                            <div className="relative">
                                <Calendar className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
                                <input
                                    type="date"
                                    value={preferences.endDate}
                                    onChange={(e) => handleInputChange('endDate', e.target.value)}
                                    className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                                />
                            </div>
                        </div>
                    </div>

                    {/* Number of Travelers */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Number of Travelers
                        </label>
                        <div className="relative">
                            <Users className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
                            <select
                                value={preferences.travelers}
                                onChange={(e) => handleInputChange('travelers', parseInt(e.target.value))}
                                className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent appearance-none"
                            >
                                {[1, 2, 3, 4, 5, 6, 7, 8].map(num => (
                                    <option key={num} value={num}>
                                        {num} {num === 1 ? 'Traveler' : 'Travelers'}
                                    </option>
                                ))}
                            </select>
                            <ChevronDown className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5 pointer-events-none" />
                        </div>
                    </div>
                </div>
            </div>
        </motion.div>
    )

    const renderStep2 = () => (
        <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.5 }}
            className="space-y-6"
        >
            <div>
                <h2 className="text-2xl font-bold text-gray-900 mb-6">What's your budget and style?</h2>

                {/* Budget Selection */}
                <div className="space-y-4">
                    <label className="block text-sm font-medium text-gray-700 mb-3">
                        Budget Range
                    </label>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                        {budgetOptions.map((option) => (
                            <motion.div
                                key={option.value}
                                whileHover={{ scale: 1.02 }}
                                whileTap={{ scale: 0.98 }}
                                onClick={() => handleInputChange('budget', option.value)}
                                className={`p-4 border-2 rounded-lg cursor-pointer transition-all ${preferences.budget === option.value
                                        ? 'border-primary-500 bg-primary-50'
                                        : 'border-gray-200 hover:border-primary-300'
                                    }`}
                            >
                                <div className="flex items-center justify-between">
                                    <div>
                                        <h3 className="font-semibold text-gray-900">{option.label}</h3>
                                        <p className="text-sm text-gray-600">{option.range}</p>
                                    </div>
                                    <DollarSign className="w-5 h-5 text-primary-500" />
                                </div>
                            </motion.div>
                        ))}
                    </div>
                </div>

                {/* Travel Style */}
                <div className="space-y-4">
                    <label className="block text-sm font-medium text-gray-700 mb-3">
                        Travel Style (Select all that apply)
                    </label>
                    <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                        {travelStyles.map((style) => (
                            <motion.div
                                key={style.id}
                                whileHover={{ scale: 1.02 }}
                                whileTap={{ scale: 0.98 }}
                                onClick={() => handleArrayToggle('travelStyle', style.id)}
                                className={`p-3 border-2 rounded-lg cursor-pointer transition-all ${preferences.travelStyle.includes(style.id)
                                        ? 'border-primary-500 bg-primary-50'
                                        : 'border-gray-200 hover:border-primary-300'
                                    }`}
                            >
                                <div className="flex items-center space-x-2">
                                    <div className="text-primary-500">{style.icon}</div>
                                    <span className="text-sm font-medium text-gray-900">{style.label}</span>
                                </div>
                            </motion.div>
                        ))}
                    </div>
                </div>
            </div>
        </motion.div>
    )

    const renderStep3 = () => (
        <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.5 }}
            className="space-y-6"
        >
            <div>
                <h2 className="text-2xl font-bold text-gray-900 mb-6">What are your interests?</h2>

                {/* Interests Selection */}
                <div className="space-y-4">
                    <label className="block text-sm font-medium text-gray-700 mb-3">
                        Select your interests
                    </label>
                    <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                        {interests.map((interest) => (
                            <motion.div
                                key={interest}
                                whileHover={{ scale: 1.02 }}
                                whileTap={{ scale: 0.98 }}
                                onClick={() => handleArrayToggle('interests', interest)}
                                className={`p-3 border-2 rounded-lg cursor-pointer transition-all ${preferences.interests.includes(interest)
                                        ? 'border-primary-500 bg-primary-50'
                                        : 'border-gray-200 hover:border-primary-300'
                                    }`}
                            >
                                <span className="text-sm font-medium text-gray-900">{interest}</span>
                            </motion.div>
                        ))}
                    </div>
                </div>

                {/* Accommodation */}
                <div className="space-y-4">
                    <label className="block text-sm font-medium text-gray-700 mb-3">
                        Preferred Accommodation
                    </label>
                    <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                        {accommodationTypes.map((type) => (
                            <motion.div
                                key={type.value}
                                whileHover={{ scale: 1.02 }}
                                whileTap={{ scale: 0.98 }}
                                onClick={() => handleInputChange('accommodation', type.value)}
                                className={`p-3 border-2 rounded-lg cursor-pointer transition-all ${preferences.accommodation === type.value
                                        ? 'border-primary-500 bg-primary-50'
                                        : 'border-gray-200 hover:border-primary-300'
                                    }`}
                            >
                                <span className="text-sm font-medium text-gray-900">{type.label}</span>
                            </motion.div>
                        ))}
                    </div>
                </div>

                {/* Transportation */}
                <div className="space-y-4">
                    <label className="block text-sm font-medium text-gray-700 mb-3">
                        Preferred Transportation
                    </label>
                    <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                        {transportationOptions.map((option) => (
                            <motion.div
                                key={option.value}
                                whileHover={{ scale: 1.02 }}
                                whileTap={{ scale: 0.98 }}
                                onClick={() => handleInputChange('transportation', option.value)}
                                className={`p-3 border-2 rounded-lg cursor-pointer transition-all ${preferences.transportation === option.value
                                        ? 'border-primary-500 bg-primary-50'
                                        : 'border-gray-200 hover:border-primary-300'
                                    }`}
                            >
                                <span className="text-sm font-medium text-gray-900">{option.label}</span>
                            </motion.div>
                        ))}
                    </div>
                </div>
            </div>
        </motion.div>
    )

    return (
        <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50">
            {/* Navigation */}
            <nav className="bg-white shadow-sm">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <div className="flex justify-between items-center h-16">
                        <Link href="/" className="flex items-center space-x-2">
                            <Plane className="w-8 h-8 text-primary-600" />
                            <span className="text-2xl font-bold text-gray-900">AuraTravel</span>
                            <span className="text-sm bg-primary-100 text-primary-800 px-2 py-1 rounded-full">AI</span>
                        </Link>
                        <div className="flex items-center space-x-4">
                            <Link href="/dashboard" className="text-gray-600 hover:text-primary-600">
                                Dashboard
                            </Link>
                            <Link href="/login" className="bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700">
                                Sign In
                            </Link>
                        </div>
                    </div>
                </div>
            </nav>

            <div className="max-w-4xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
                {/* Progress Steps */}
                <div className="mb-8">
                    <div className="flex items-center justify-center space-x-8">
                        {[1, 2, 3].map((stepNumber) => (
                            <div key={stepNumber} className="flex items-center">
                                <div className={`w-10 h-10 rounded-full flex items-center justify-center font-semibold ${step >= stepNumber
                                        ? 'bg-primary-600 text-white'
                                        : 'bg-gray-200 text-gray-600'
                                    }`}>
                                    {stepNumber}
                                </div>
                                {stepNumber < 3 && (
                                    <div className={`w-16 h-1 ml-4 ${step > stepNumber ? 'bg-primary-600' : 'bg-gray-200'
                                        }`} />
                                )}
                            </div>
                        ))}
                    </div>
                    <div className="flex justify-center mt-4">
                        <div className="text-center">
                            <h1 className="text-3xl font-bold text-gray-900">Plan Your Perfect Trip</h1>
                            <p className="text-gray-600 mt-2">Let AI create a personalized itinerary just for you</p>
                        </div>
                    </div>
                </div>

                {/* Form Container */}
                <div className="bg-white rounded-2xl shadow-xl p-8">
                    {step === 1 && renderStep1()}
                    {step === 2 && renderStep2()}
                    {step === 3 && renderStep3()}

                    {/* Navigation Buttons */}
                    <div className="flex justify-between mt-8 pt-6 border-t border-gray-200">
                        <motion.button
                            whileHover={{ scale: 1.02 }}
                            whileTap={{ scale: 0.98 }}
                            onClick={() => setStep(Math.max(1, step - 1))}
                            disabled={step === 1}
                            className={`px-6 py-3 rounded-lg font-semibold ${step === 1
                                    ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                                    : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                                }`}
                        >
                            Previous
                        </motion.button>

                        <motion.button
                            whileHover={{ scale: 1.02 }}
                            whileTap={{ scale: 0.98 }}
                            onClick={() => {
                                if (step < 3) {
                                    setStep(step + 1)
                                } else {
                                    handleSubmit()
                                }
                            }}
                            className="bg-primary-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-primary-700 flex items-center gap-2"
                        >
                            {step === 3 ? 'Generate Itinerary' : 'Next'}
                            <ArrowRight className="w-4 h-4" />
                        </motion.button>
                    </div>
                </div>
            </div>
        </div>
    )
}