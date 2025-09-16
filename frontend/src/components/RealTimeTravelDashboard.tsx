'use client'

import React, { useState, useEffect, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import {
    Activity,
    Zap,
    TrendingUp,
    Clock,
    MapPin,
    CloudRain,
    Users,
    Bell,
    Download,
    Globe,
    Sparkles,
    RefreshCw,
    Settings,
    BarChart3,
    Wifi,
    WifiOff,
    Battery,
    Signal,
    AlertTriangle,
    CheckCircle,
    Info
} from 'lucide-react'

// Import our newly created components
import RealTimeNotificationSystem from './RealTimeNotificationSystem'
import DynamicReplanningUI from './DynamicReplanningUI'
import ItineraryDelivery from './ItineraryDelivery'
import LanguageSelector from './LanguageSelector'

// Types for real-time dashboard
interface TripStatus {
    tripId: string
    status: 'planning' | 'active' | 'completed' | 'cancelled'
    progress: number
    currentActivity?: string
    nextActivity?: string
    location?: {
        latitude: number
        longitude: number
        address: string
    }
    weather?: {
        condition: string
        temperature: number
        description: string
    }
    timeline: {
        started: Date
        estimatedEnd: Date
        actualEnd?: Date
    }
}

interface RealTimeMetrics {
    connectionStatus: 'online' | 'offline' | 'poor'
    lastUpdate: Date
    monitoringActive: boolean
    activeAlerts: number
    totalNotifications: number
    replanningEvents: number
    deliveryCount: number
    languageDetected?: string
}

interface DashboardProps {
    tripId: string
    userId?: string
    initialLanguage?: string
    className?: string
}

export default function RealTimeTravelDashboard({
    tripId,
    userId,
    initialLanguage = 'en',
    className = ''
}: DashboardProps) {
    const [tripStatus, setTripStatus] = useState<TripStatus | null>(null)
    const [metrics, setMetrics] = useState<RealTimeMetrics>({
        connectionStatus: 'online',
        lastUpdate: new Date(),
        monitoringActive: true,
        activeAlerts: 0,
        totalNotifications: 0,
        replanningEvents: 0,
        deliveryCount: 0
    })
    const [currentLanguage, setCurrentLanguage] = useState(initialLanguage)
    const [activeTab, setActiveTab] = useState<'overview' | 'notifications' | 'replanning' | 'delivery' | 'settings'>('overview')
    const [isExpanded, setIsExpanded] = useState(true)
    const [autoRefresh, setAutoRefresh] = useState(true)
    const [refreshInterval, setRefreshInterval] = useState(30) // seconds
    const [dashboardSettings, setDashboardSettings] = useState({
        showMetrics: true,
        showQuickActions: true,
        compactMode: false,
        realTimeUpdates: true,
        soundEnabled: true,
        animationsEnabled: true
    })

    // Real-time connection monitoring
    useEffect(() => {
        const checkConnection = () => {
            const online = navigator.onLine
            setMetrics(prev => ({
                ...prev,
                connectionStatus: online ? 'online' : 'offline'
            }))
        }

        window.addEventListener('online', checkConnection)
        window.addEventListener('offline', checkConnection)

        return () => {
            window.removeEventListener('online', checkConnection)
            window.removeEventListener('offline', checkConnection)
        }
    }, [])

    // Auto-refresh trip status
    useEffect(() => {
        if (!autoRefresh || !tripId) return

        const interval = setInterval(async () => {
            try {
                const response = await fetch(`/api/trips/${tripId}/status`)
                const data = await response.json()

                setTripStatus(data.tripStatus)
                setMetrics(prev => ({
                    ...prev,
                    lastUpdate: new Date(),
                    activeAlerts: data.activeAlerts || 0
                }))
            } catch (error) {
                console.error('Failed to fetch trip status:', error)
                setMetrics(prev => ({
                    ...prev,
                    connectionStatus: 'poor'
                }))
            }
        }, refreshInterval * 1000)

        return () => clearInterval(interval)
    }, [autoRefresh, tripId, refreshInterval])

    // Handle language change
    const handleLanguageChange = useCallback((language: any) => {
        setCurrentLanguage(language.code)

        // Update metrics
        setMetrics(prev => ({
            ...prev,
            languageDetected: language.code
        }))

        // Save preference
        localStorage.setItem('dashboardLanguage', language.code)
    }, [])

    // Handle notification events
    const handleNotificationAction = useCallback((notification: any) => {
        setMetrics(prev => ({
            ...prev,
            totalNotifications: prev.totalNotifications + 1
        }))

        // Auto-switch to relevant tab based on notification type
        if (notification.type === 'weather_alert' || notification.type === 'delay_alert') {
            setActiveTab('replanning')
        } else if (notification.type === 'itinerary_update') {
            setActiveTab('delivery')
        }
    }, [])

    // Handle replanning events
    const handleReplanAccept = useCallback((option: any) => {
        setMetrics(prev => ({
            ...prev,
            replanningEvents: prev.replanningEvents + 1
        }))
    }, [])

    // Handle delivery completion
    const handleDeliveryComplete = useCallback(() => {
        setMetrics(prev => ({
            ...prev,
            deliveryCount: prev.deliveryCount + 1
        }))
    }, [])

    // Get connection status color
    const getConnectionColor = (status: string) => {
        switch (status) {
            case 'online':
                return 'text-green-600'
            case 'poor':
                return 'text-yellow-600'
            case 'offline':
                return 'text-red-600'
            default:
                return 'text-gray-600'
        }
    }

    // Get connection icon
    const getConnectionIcon = (status: string) => {
        switch (status) {
            case 'online':
                return Wifi
            case 'poor':
                return Signal
            case 'offline':
                return WifiOff
            default:
                return Wifi
        }
    }

    // Dashboard tabs configuration
    const tabs = [
        {
            id: 'overview',
            label: 'Overview',
            icon: Activity,
            badge: metrics.activeAlerts > 0 ? metrics.activeAlerts : undefined
        },
        {
            id: 'notifications',
            label: 'Notifications',
            icon: Bell,
            badge: metrics.totalNotifications > 0 ? metrics.totalNotifications : undefined
        },
        {
            id: 'replanning',
            label: 'Replanning',
            icon: RefreshCw,
            badge: metrics.replanningEvents > 0 ? metrics.replanningEvents : undefined
        },
        {
            id: 'delivery',
            label: 'Delivery',
            icon: Download,
            badge: metrics.deliveryCount > 0 ? '✓' : undefined
        },
        {
            id: 'settings',
            label: 'Settings',
            icon: Settings,
            badge: undefined
        }
    ]

    return (
        <div className={`bg-white rounded-lg shadow-lg border border-gray-200 ${dashboardSettings.compactMode ? 'text-sm' : ''} ${className}`}>
            {/* Header */}
            <div className="p-4 border-b border-gray-200 bg-gradient-to-r from-blue-50 to-indigo-50">
                <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-3">
                        <div className="p-2 bg-blue-100 rounded-full">
                            <Zap className="w-5 h-5 text-blue-600" />
                        </div>
                        <div>
                            <h2 className="text-lg font-semibold text-gray-900 flex items-center">
                                Real-Time Travel Dashboard
                                {dashboardSettings.realTimeUpdates && (
                                    <span className="ml-2 w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
                                )}
                            </h2>
                            <p className="text-sm text-gray-600">
                                Trip ID: {tripId} • {tripStatus?.status || 'Loading...'}
                            </p>
                        </div>
                    </div>

                    <div className="flex items-center space-x-3">
                        {/* Connection Status */}
                        <div className={`flex items-center text-sm ${getConnectionColor(metrics.connectionStatus)}`}>
                            {React.createElement(getConnectionIcon(metrics.connectionStatus), {
                                className: 'w-4 h-4 mr-1'
                            })}
                            <span className="capitalize">{metrics.connectionStatus}</span>
                        </div>

                        {/* Language Selector */}
                        <LanguageSelector
                            currentLanguage={currentLanguage}
                            onLanguageChange={handleLanguageChange}
                            userId={userId}
                        />

                        {/* Expand/Collapse */}
                        <button
                            onClick={() => setIsExpanded(!isExpanded)}
                            className="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded"
                        >
                            <motion.div
                                animate={{ rotate: isExpanded ? 180 : 0 }}
                                transition={{ duration: 0.2 }}
                            >
                                <TrendingUp className="w-4 h-4" />
                            </motion.div>
                        </button>
                    </div>
                </div>

                {/* Real-time Metrics */}
                <AnimatePresence>
                    {dashboardSettings.showMetrics && isExpanded && (
                        <motion.div
                            initial={{ opacity: 0, height: 0 }}
                            animate={{ opacity: 1, height: 'auto' }}
                            exit={{ opacity: 0, height: 0 }}
                            className="mt-4 grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-3"
                        >
                            <div className="text-center p-2 bg-white rounded border">
                                <Activity className="w-4 h-4 mx-auto mb-1 text-gray-600" />
                                <p className="text-xs text-gray-500">Status</p>
                                <p className="font-medium">{tripStatus?.status || 'N/A'}</p>
                            </div>
                            <div className="text-center p-2 bg-white rounded border">
                                <Bell className="w-4 h-4 mx-auto mb-1 text-blue-600" />
                                <p className="text-xs text-gray-500">Alerts</p>
                                <p className="font-medium">{metrics.activeAlerts}</p>
                            </div>
                            <div className="text-center p-2 bg-white rounded border">
                                <RefreshCw className="w-4 h-4 mx-auto mb-1 text-orange-600" />
                                <p className="text-xs text-gray-500">Replans</p>
                                <p className="font-medium">{metrics.replanningEvents}</p>
                            </div>
                            <div className="text-center p-2 bg-white rounded border">
                                <Download className="w-4 h-4 mx-auto mb-1 text-green-600" />
                                <p className="text-xs text-gray-500">Downloads</p>
                                <p className="font-medium">{metrics.deliveryCount}</p>
                            </div>
                            <div className="text-center p-2 bg-white rounded border">
                                <Clock className="w-4 h-4 mx-auto mb-1 text-purple-600" />
                                <p className="text-xs text-gray-500">Updated</p>
                                <p className="font-medium">{metrics.lastUpdate.toLocaleTimeString()}</p>
                            </div>
                            <div className="text-center p-2 bg-white rounded border">
                                <Globe className="w-4 h-4 mx-auto mb-1 text-indigo-600" />
                                <p className="text-xs text-gray-500">Language</p>
                                <p className="font-medium">{currentLanguage.toUpperCase()}</p>
                            </div>
                        </motion.div>
                    )}
                </AnimatePresence>
            </div>

            {/* Tab Navigation */}
            <AnimatePresence>
                {isExpanded && (
                    <motion.div
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: 'auto' }}
                        exit={{ opacity: 0, height: 0 }}
                        className="border-b border-gray-200"
                    >
                        <nav className="flex space-x-0">
                            {tabs.map((tab) => (
                                <button
                                    key={tab.id}
                                    onClick={() => setActiveTab(tab.id as any)}
                                    className={`flex items-center px-4 py-3 text-sm font-medium border-b-2 transition-colors relative ${activeTab === tab.id
                                            ? 'border-blue-500 text-blue-600 bg-blue-50'
                                            : 'border-transparent text-gray-500 hover:text-gray-700 hover:bg-gray-50'
                                        }`}
                                >
                                    <tab.icon className="w-4 h-4 mr-2" />
                                    {tab.label}
                                    {tab.badge && (
                                        <span className="ml-2 px-2 py-0.5 text-xs bg-red-100 text-red-800 rounded-full">
                                            {tab.badge}
                                        </span>
                                    )}
                                </button>
                            ))}
                        </nav>
                    </motion.div>
                )}
            </AnimatePresence>

            {/* Tab Content */}
            <AnimatePresence mode="wait">
                {isExpanded && (
                    <motion.div
                        key={activeTab}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -20 }}
                        transition={{ duration: dashboardSettings.animationsEnabled ? 0.3 : 0 }}
                        className="p-6"
                    >
                        {activeTab === 'overview' && (
                            <div className="space-y-6">
                                {/* Trip Progress */}
                                {tripStatus && (
                                    <div className="p-4 bg-gradient-to-r from-blue-50 to-indigo-50 rounded-lg">
                                        <h3 className="font-semibold text-gray-900 mb-3">Trip Progress</h3>
                                        <div className="space-y-3">
                                            <div className="flex justify-between text-sm">
                                                <span>Overall Progress</span>
                                                <span>{tripStatus.progress}%</span>
                                            </div>
                                            <div className="w-full bg-gray-200 rounded-full h-2">
                                                <motion.div
                                                    className="bg-blue-600 h-2 rounded-full"
                                                    initial={{ width: 0 }}
                                                    animate={{ width: `${tripStatus.progress}%` }}
                                                    transition={{ duration: 1 }}
                                                />
                                            </div>
                                            {tripStatus.currentActivity && (
                                                <div className="flex items-center text-sm text-gray-600">
                                                    <MapPin className="w-4 h-4 mr-1" />
                                                    <span>Current: {tripStatus.currentActivity}</span>
                                                </div>
                                            )}
                                            {tripStatus.nextActivity && (
                                                <div className="flex items-center text-sm text-gray-600">
                                                    <Clock className="w-4 h-4 mr-1" />
                                                    <span>Next: {tripStatus.nextActivity}</span>
                                                </div>
                                            )}
                                        </div>
                                    </div>
                                )}

                                {/* Quick Actions */}
                                {dashboardSettings.showQuickActions && (
                                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                                        <button
                                            onClick={() => setActiveTab('notifications')}
                                            className="p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow text-left"
                                        >
                                            <Bell className="w-6 h-6 text-blue-600 mb-2" />
                                            <h4 className="font-medium">View Notifications</h4>
                                            <p className="text-sm text-gray-600">Check alerts and updates</p>
                                        </button>
                                        <button
                                            onClick={() => setActiveTab('replanning')}
                                            className="p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow text-left"
                                        >
                                            <RefreshCw className="w-6 h-6 text-orange-600 mb-2" />
                                            <h4 className="font-medium">Dynamic Replanning</h4>
                                            <p className="text-sm text-gray-600">Adapt to changes</p>
                                        </button>
                                        <button
                                            onClick={() => setActiveTab('delivery')}
                                            className="p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow text-left"
                                        >
                                            <Download className="w-6 h-6 text-green-600 mb-2" />
                                            <h4 className="font-medium">Download Itinerary</h4>
                                            <p className="text-sm text-gray-600">Get files and share</p>
                                        </button>
                                    </div>
                                )}

                                {/* Weather Widget */}
                                {tripStatus?.weather && (
                                    <div className="p-4 bg-gradient-to-r from-sky-50 to-blue-50 rounded-lg">
                                        <h3 className="font-semibold text-gray-900 mb-3">Current Weather</h3>
                                        <div className="flex items-center">
                                            <CloudRain className="w-8 h-8 text-blue-600 mr-3" />
                                            <div>
                                                <p className="font-medium">{tripStatus.weather.temperature}°C</p>
                                                <p className="text-sm text-gray-600">{tripStatus.weather.description}</p>
                                            </div>
                                        </div>
                                    </div>
                                )}
                            </div>
                        )}

                        {activeTab === 'notifications' && (
                            <RealTimeNotificationSystem
                                userId={userId}
                                locale={currentLanguage}
                                onNotificationAction={handleNotificationAction}
                            />
                        )}

                        {activeTab === 'replanning' && (
                            <DynamicReplanningUI
                                tripId={tripId}
                                currentItinerary={tripStatus}
                                locale={currentLanguage}
                                onReplanAccept={handleReplanAccept}
                            />
                        )}

                        {activeTab === 'delivery' && (
                            <ItineraryDelivery
                                tripId={tripId}
                                tripTitle={`Trip ${tripId}`}
                                currentItinerary={tripStatus}
                                locale={currentLanguage}
                            />
                        )}

                        {activeTab === 'settings' && (
                            <div className="space-y-6">
                                <h3 className="font-semibold text-gray-900">Dashboard Settings</h3>

                                <div className="space-y-4">
                                    <label className="flex items-center justify-between">
                                        <span className="text-sm font-medium text-gray-700">Show metrics</span>
                                        <input
                                            type="checkbox"
                                            checked={dashboardSettings.showMetrics}
                                            onChange={(e) => setDashboardSettings(prev => ({
                                                ...prev,
                                                showMetrics: e.target.checked
                                            }))}
                                            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                        />
                                    </label>

                                    <label className="flex items-center justify-between">
                                        <span className="text-sm font-medium text-gray-700">Show quick actions</span>
                                        <input
                                            type="checkbox"
                                            checked={dashboardSettings.showQuickActions}
                                            onChange={(e) => setDashboardSettings(prev => ({
                                                ...prev,
                                                showQuickActions: e.target.checked
                                            }))}
                                            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                        />
                                    </label>

                                    <label className="flex items-center justify-between">
                                        <span className="text-sm font-medium text-gray-700">Compact mode</span>
                                        <input
                                            type="checkbox"
                                            checked={dashboardSettings.compactMode}
                                            onChange={(e) => setDashboardSettings(prev => ({
                                                ...prev,
                                                compactMode: e.target.checked
                                            }))}
                                            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                        />
                                    </label>

                                    <label className="flex items-center justify-between">
                                        <span className="text-sm font-medium text-gray-700">Real-time updates</span>
                                        <input
                                            type="checkbox"
                                            checked={dashboardSettings.realTimeUpdates}
                                            onChange={(e) => setDashboardSettings(prev => ({
                                                ...prev,
                                                realTimeUpdates: e.target.checked
                                            }))}
                                            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                        />
                                    </label>

                                    <label className="flex items-center justify-between">
                                        <span className="text-sm font-medium text-gray-700">Auto refresh</span>
                                        <input
                                            type="checkbox"
                                            checked={autoRefresh}
                                            onChange={(e) => setAutoRefresh(e.target.checked)}
                                            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                        />
                                    </label>

                                    <div className="space-y-2">
                                        <label className="block text-sm font-medium text-gray-700">
                                            Refresh interval (seconds)
                                        </label>
                                        <input
                                            type="range"
                                            min="10"
                                            max="300"
                                            value={refreshInterval}
                                            onChange={(e) => setRefreshInterval(Number(e.target.value))}
                                            className="w-full"
                                        />
                                        <span className="text-sm text-gray-500">{refreshInterval}s</span>
                                    </div>
                                </div>
                            </div>
                        )}
                    </motion.div>
                )}
            </AnimatePresence>
        </div>
    )
}