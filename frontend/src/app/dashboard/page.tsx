'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { motion } from 'framer-motion'
import {
    Plus,
    MapPin,
    Calendar,
    Users,
    Star,
    Clock,
    Plane,
    Camera,
    Heart,
    Share2,
    Edit,
    Trash2,
    Filter,
    Search,
    Bell,
    Settings,
    LogOut,
    User
} from 'lucide-react'

interface Trip {
    id: string
    title: string
    destination: string
    startDate: string
    endDate: string
    status: 'upcoming' | 'completed' | 'draft'
    travelers: number
    image: string
    budget: number
    activities: number
}

export default function DashboardPage() {
    const [trips, setTrips] = useState<Trip[]>([
        {
            id: '1',
            title: 'Romantic Getaway to Paris',
            destination: 'Paris, France',
            startDate: '2024-10-15',
            endDate: '2024-10-22',
            status: 'upcoming',
            travelers: 2,
            image: '/images/paris.jpg',
            budget: 3500,
            activities: 12
        },
        {
            id: '2',
            title: 'Adventure in Tokyo',
            destination: 'Tokyo, Japan',
            startDate: '2024-08-10',
            endDate: '2024-08-20',
            status: 'completed',
            travelers: 1,
            image: '/images/tokyo.jpg',
            budget: 4200,
            activities: 18
        },
        {
            id: '3',
            title: 'Beach Vacation Draft',
            destination: 'Bali, Indonesia',
            startDate: '',
            endDate: '',
            status: 'draft',
            travelers: 4,
            image: '/images/bali.jpg',
            budget: 0,
            activities: 0
        }
    ])

    const [filter, setFilter] = useState<'all' | 'upcoming' | 'completed' | 'draft'>('all')
    const [searchTerm, setSearchTerm] = useState('')

    const filteredTrips = trips.filter(trip => {
        const matchesFilter = filter === 'all' || trip.status === filter
        const matchesSearch = trip.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
            trip.destination.toLowerCase().includes(searchTerm.toLowerCase())
        return matchesFilter && matchesSearch
    })

    const stats = {
        totalTrips: trips.length,
        upcomingTrips: trips.filter(trip => trip.status === 'upcoming').length,
        completedTrips: trips.filter(trip => trip.status === 'completed').length,
        totalBudget: trips.reduce((sum, trip) => sum + trip.budget, 0)
    }

    const formatDate = (dateString: string) => {
        if (!dateString) return 'TBD'
        return new Date(dateString).toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric'
        })
    }

    const getStatusColor = (status: Trip['status']) => {
        switch (status) {
            case 'upcoming': return 'bg-blue-100 text-blue-800'
            case 'completed': return 'bg-green-100 text-green-800'
            case 'draft': return 'bg-gray-100 text-gray-800'
            default: return 'bg-gray-100 text-gray-800'
        }
    }

    return (
        <div className="min-h-screen bg-gray-50">
            {/* Navigation */}
            <nav className="bg-white shadow-sm border-b border-gray-200">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <div className="flex justify-between items-center h-16">
                        <Link href="/" className="flex items-center space-x-2">
                            <Plane className="w-8 h-8 text-primary-600" />
                            <span className="text-2xl font-bold text-gray-900">AuraTravel</span>
                            <span className="text-sm bg-primary-100 text-primary-800 px-2 py-1 rounded-full">AI</span>
                        </Link>

                        <div className="flex items-center space-x-4">
                            <button className="p-2 text-gray-600 hover:text-primary-600 relative">
                                <Bell className="w-5 h-5" />
                                <span className="absolute -top-1 -right-1 w-3 h-3 bg-red-500 rounded-full"></span>
                            </button>
                            <button className="p-2 text-gray-600 hover:text-primary-600">
                                <Settings className="w-5 h-5" />
                            </button>
                            <div className="flex items-center space-x-2">
                                <div className="w-8 h-8 bg-primary-600 rounded-full flex items-center justify-center">
                                    <User className="w-4 h-4 text-white" />
                                </div>
                                <span className="text-sm font-medium text-gray-700">John Doe</span>
                            </div>
                            <button className="p-2 text-gray-600 hover:text-red-600">
                                <LogOut className="w-5 h-5" />
                            </button>
                        </div>
                    </div>
                </div>
            </nav>

            <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
                {/* Welcome Section */}
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.6 }}
                    className="mb-8"
                >
                    <h1 className="text-3xl font-bold text-gray-900 mb-2">Welcome back, John!</h1>
                    <p className="text-gray-600">Ready for your next adventure? Let's plan something amazing.</p>
                </motion.div>

                {/* Stats Cards */}
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.6, delay: 0.1 }}
                    className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8"
                >
                    <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm font-medium text-gray-600">Total Trips</p>
                                <p className="text-2xl font-bold text-gray-900">{stats.totalTrips}</p>
                            </div>
                            <div className="p-3 bg-primary-100 rounded-full">
                                <MapPin className="w-6 h-6 text-primary-600" />
                            </div>
                        </div>
                    </div>

                    <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm font-medium text-gray-600">Upcoming</p>
                                <p className="text-2xl font-bold text-gray-900">{stats.upcomingTrips}</p>
                            </div>
                            <div className="p-3 bg-blue-100 rounded-full">
                                <Calendar className="w-6 h-6 text-blue-600" />
                            </div>
                        </div>
                    </div>

                    <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm font-medium text-gray-600">Completed</p>
                                <p className="text-2xl font-bold text-gray-900">{stats.completedTrips}</p>
                            </div>
                            <div className="p-3 bg-green-100 rounded-full">
                                <Star className="w-6 h-6 text-green-600" />
                            </div>
                        </div>
                    </div>

                    <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm font-medium text-gray-600">Total Budget</p>
                                <p className="text-2xl font-bold text-gray-900">${stats.totalBudget.toLocaleString()}</p>
                            </div>
                            <div className="p-3 bg-yellow-100 rounded-full">
                                <Clock className="w-6 h-6 text-yellow-600" />
                            </div>
                        </div>
                    </div>
                </motion.div>

                {/* Quick Actions */}
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.6, delay: 0.2 }}
                    className="mb-8"
                >
                    <div className="flex flex-col sm:flex-row gap-4">
                        <Link href="/plan" className="flex-1">
                            <motion.button
                                whileHover={{ scale: 1.02 }}
                                whileTap={{ scale: 0.98 }}
                                className="w-full bg-primary-600 text-white py-4 px-6 rounded-xl font-semibold hover:bg-primary-700 transition-colors flex items-center justify-center gap-2"
                            >
                                <Plus className="w-5 h-5" />
                                Plan New Trip
                            </motion.button>
                        </Link>
                        <motion.button
                            whileHover={{ scale: 1.02 }}
                            whileTap={{ scale: 0.98 }}
                            className="bg-white border border-gray-300 text-gray-700 py-4 px-6 rounded-xl font-semibold hover:bg-gray-50 transition-colors flex items-center justify-center gap-2"
                        >
                            <Camera className="w-5 h-5" />
                            Browse Inspiration
                        </motion.button>
                    </div>
                </motion.div>

                {/* Trips Section */}
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.6, delay: 0.3 }}
                >
                    <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
                        {/* Section Header */}
                        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-6">
                            <h2 className="text-xl font-bold text-gray-900 mb-4 sm:mb-0">Your Trips</h2>

                            <div className="flex flex-col sm:flex-row gap-3">
                                {/* Search */}
                                <div className="relative">
                                    <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                                    <input
                                        type="text"
                                        placeholder="Search trips..."
                                        value={searchTerm}
                                        onChange={(e) => setSearchTerm(e.target.value)}
                                        className="pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm"
                                    />
                                </div>

                                {/* Filter */}
                                <div className="relative">
                                    <Filter className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                                    <select
                                        value={filter}
                                        onChange={(e) => setFilter(e.target.value as any)}
                                        className="pl-10 pr-8 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm appearance-none bg-white"
                                    >
                                        <option value="all">All Trips</option>
                                        <option value="upcoming">Upcoming</option>
                                        <option value="completed">Completed</option>
                                        <option value="draft">Drafts</option>
                                    </select>
                                </div>
                            </div>
                        </div>

                        {/* Trips Grid */}
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                            {filteredTrips.map((trip, index) => (
                                <motion.div
                                    key={trip.id}
                                    initial={{ opacity: 0, y: 20 }}
                                    animate={{ opacity: 1, y: 0 }}
                                    transition={{ duration: 0.5, delay: index * 0.1 }}
                                    className="border border-gray-200 rounded-xl overflow-hidden hover:shadow-lg transition-shadow"
                                >
                                    {/* Trip Image */}
                                    <div className="relative h-48 bg-gradient-to-r from-primary-400 to-secondary-400">
                                        <div className="absolute inset-0 bg-black/20"></div>
                                        <div className="absolute top-4 left-4">
                                            <span className={`px-3 py-1 rounded-full text-xs font-semibold ${getStatusColor(trip.status)}`}>
                                                {trip.status.charAt(0).toUpperCase() + trip.status.slice(1)}
                                            </span>
                                        </div>
                                        <div className="absolute top-4 right-4 flex gap-2">
                                            <button className="p-2 bg-white/20 rounded-full backdrop-blur-sm hover:bg-white/30 transition-colors">
                                                <Heart className="w-4 h-4 text-white" />
                                            </button>
                                            <button className="p-2 bg-white/20 rounded-full backdrop-blur-sm hover:bg-white/30 transition-colors">
                                                <Share2 className="w-4 h-4 text-white" />
                                            </button>
                                        </div>
                                    </div>

                                    {/* Trip Content */}
                                    <div className="p-6">
                                        <h3 className="text-lg font-semibold text-gray-900 mb-2">{trip.title}</h3>
                                        <div className="flex items-center text-gray-600 mb-3">
                                            <MapPin className="w-4 h-4 mr-1" />
                                            <span className="text-sm">{trip.destination}</span>
                                        </div>

                                        <div className="space-y-2 mb-4">
                                            <div className="flex items-center justify-between text-sm">
                                                <span className="text-gray-600">Dates:</span>
                                                <span className="font-medium">
                                                    {formatDate(trip.startDate)} - {formatDate(trip.endDate)}
                                                </span>
                                            </div>
                                            <div className="flex items-center justify-between text-sm">
                                                <span className="text-gray-600">Travelers:</span>
                                                <span className="font-medium">{trip.travelers}</span>
                                            </div>
                                            {trip.budget > 0 && (
                                                <div className="flex items-center justify-between text-sm">
                                                    <span className="text-gray-600">Budget:</span>
                                                    <span className="font-medium">${trip.budget.toLocaleString()}</span>
                                                </div>
                                            )}
                                        </div>

                                        {/* Actions */}
                                        <div className="flex gap-2">
                                            <motion.button
                                                whileHover={{ scale: 1.02 }}
                                                whileTap={{ scale: 0.98 }}
                                                className="flex-1 bg-primary-600 text-white py-2 px-4 rounded-lg text-sm font-semibold hover:bg-primary-700 transition-colors"
                                            >
                                                {trip.status === 'draft' ? 'Continue Planning' : 'View Details'}
                                            </motion.button>
                                            <button className="p-2 text-gray-600 hover:text-primary-600 border border-gray-300 rounded-lg hover:border-primary-300 transition-colors">
                                                <Edit className="w-4 h-4" />
                                            </button>
                                            <button className="p-2 text-gray-600 hover:text-red-600 border border-gray-300 rounded-lg hover:border-red-300 transition-colors">
                                                <Trash2 className="w-4 h-4" />
                                            </button>
                                        </div>
                                    </div>
                                </motion.div>
                            ))}
                        </div>

                        {filteredTrips.length === 0 && (
                            <div className="text-center py-12">
                                <div className="text-gray-400 mb-4">
                                    <MapPin className="w-12 h-12 mx-auto" />
                                </div>
                                <h3 className="text-lg font-semibold text-gray-900 mb-2">No trips found</h3>
                                <p className="text-gray-600 mb-6">
                                    {searchTerm ? 'Try adjusting your search terms.' : 'Start planning your first adventure!'}
                                </p>
                                <Link href="/plan">
                                    <motion.button
                                        whileHover={{ scale: 1.02 }}
                                        whileTap={{ scale: 0.98 }}
                                        className="bg-primary-600 text-white py-3 px-6 rounded-lg font-semibold hover:bg-primary-700 transition-colors inline-flex items-center gap-2"
                                    >
                                        <Plus className="w-4 h-4" />
                                        Plan New Trip
                                    </motion.button>
                                </Link>
                            </div>
                        )}
                    </div>
                </motion.div>
            </div>
        </div>
    )
}