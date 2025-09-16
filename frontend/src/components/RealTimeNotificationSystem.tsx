'use client'

import React, { useState, useEffect, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import {
    Bell,
    X,
    AlertTriangle,
    Info,
    CheckCircle,
    CloudRain,
    Clock,
    Plane,
    MapPin,
    Calendar,
    User,
    Settings,
    Volume2,
    VolumeX,
    Download,
    ExternalLink,
    Globe,
    Wifi,
    WifiOff
} from 'lucide-react'
import { getMessaging, onMessage, getToken } from 'firebase/messaging'
import { getFirestore, doc, setDoc, serverTimestamp } from 'firebase/firestore'
import app from '@/lib/firebase'

// Types for notifications
interface Notification {
    id: string
    type: 'weather_alert' | 'itinerary_update' | 'trip_reminder' | 'delay_alert' | 'booking_confirm' | 'general'
    title: string
    body: string
    priority: 'high' | 'normal' | 'low'
    timestamp: Date
    read: boolean
    actionUrl?: string
    locale?: string
    data?: {
        tripId?: string
        activityId?: string
        severity?: string
        weatherCondition?: string
        delayMinutes?: number
        newTime?: string
        location?: string
        [key: string]: any
    }
}

interface NotificationSettings {
    enabled: boolean
    sound: boolean
    weatherAlerts: boolean
    delayAlerts: boolean
    itineraryUpdates: boolean
    tripReminders: boolean
    bookingConfirmations: boolean
    pushNotifications: boolean
    emailNotifications: boolean
    smsNotifications: boolean
}

interface RealTimeNotificationSystemProps {
    userId?: string
    locale?: string
    onNotificationAction?: (notification: Notification) => void
    className?: string
}

export default function RealTimeNotificationSystem({
    userId,
    locale = 'en',
    onNotificationAction,
    className = ''
}: RealTimeNotificationSystemProps) {
    const [notifications, setNotifications] = useState<Notification[]>([])
    const [unreadCount, setUnreadCount] = useState(0)
    const [isOpen, setIsOpen] = useState(false)
    const [isOnline, setIsOnline] = useState(true)
    const [settings, setSettings] = useState<NotificationSettings>({
        enabled: true,
        sound: true,
        weatherAlerts: true,
        delayAlerts: true,
        itineraryUpdates: true,
        tripReminders: true,
        bookingConfirmations: true,
        pushNotifications: true,
        emailNotifications: false,
        smsNotifications: false
    })
    const [showSettings, setShowSettings] = useState(false)
    const [deviceToken, setDeviceToken] = useState<string | null>(null)

    // Initialize Firebase messaging and register device token
    useEffect(() => {
        const initializeMessaging = async () => {
            try {
                const messaging = getMessaging(app)

                // Request permission for notifications
                const permission = await Notification.requestPermission()
                if (permission === 'granted') {
                    // Get FCM token
                    const token = await getToken(messaging, {
                        vapidKey: process.env.NEXT_PUBLIC_VAPID_KEY
                    })

                    if (token && userId) {
                        setDeviceToken(token)
                        await registerDeviceToken(token)
                    }

                    // Listen for foreground messages
                    onMessage(messaging, (payload) => {
                        handleIncomingNotification(payload)
                    })
                }
            } catch (error) {
                console.error('Error initializing messaging:', error)
            }
        }

        if (userId && typeof window !== 'undefined') {
            initializeMessaging()
        }
    }, [userId])

    // Monitor online status
    useEffect(() => {
        const handleOnline = () => setIsOnline(true)
        const handleOffline = () => setIsOnline(false)

        window.addEventListener('online', handleOnline)
        window.addEventListener('offline', handleOffline)

        return () => {
            window.removeEventListener('online', handleOnline)
            window.removeEventListener('offline', handleOffline)
        }
    }, [])

    // Load notification settings from localStorage
    useEffect(() => {
        const savedSettings = localStorage.getItem('notificationSettings')
        if (savedSettings) {
            setSettings(JSON.parse(savedSettings))
        }
    }, [])

    // Save notification settings to localStorage
    useEffect(() => {
        localStorage.setItem('notificationSettings', JSON.stringify(settings))
    }, [settings])

    // Register device token with backend
    const registerDeviceToken = async (token: string) => {
        try {
            const response = await fetch('/api/notifications/register-device', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    userId,
                    deviceToken: token,
                    platform: 'web',
                    locale
                })
            })

            if (!response.ok) {
                throw new Error('Failed to register device token')
            }

            console.log('Device token registered successfully')
        } catch (error) {
            console.error('Error registering device token:', error)
        }
    }

    // Handle incoming notifications
    const handleIncomingNotification = useCallback((payload: any) => {
        const notification: Notification = {
            id: payload.messageId || Date.now().toString(),
            type: payload.data?.type || 'general',
            title: payload.notification?.title || payload.data?.title || 'New Notification',
            body: payload.notification?.body || payload.data?.body || '',
            priority: payload.data?.priority || 'normal',
            timestamp: new Date(),
            read: false,
            actionUrl: payload.data?.actionUrl,
            locale: payload.data?.locale || locale,
            data: payload.data || {}
        }

        setNotifications(prev => [notification, ...prev])
        setUnreadCount(prev => prev + 1)

        // Play sound if enabled
        if (settings.sound && settings.enabled) {
            playNotificationSound(notification.type)
        }

        // Show browser notification if tab is not active
        if (document.hidden && 'Notification' in window && Notification.permission === 'granted') {
            const browserNotification = new Notification(notification.title, {
                body: notification.body,
                icon: '/notification-icon.png', // Use a static image path
                tag: notification.id,
                badge: '/notification-badge.png'
            })

            browserNotification.onclick = () => {
                window.focus()
                handleNotificationClick(notification)
                browserNotification.close()
            }
        }
    }, [settings, locale])

    // Play notification sound based on type
    const playNotificationSound = (type: string) => {
        const audio = new Audio()
        switch (type) {
            case 'weather_alert':
            case 'delay_alert':
                audio.src = '/sounds/alert.mp3'
                break
            case 'booking_confirm':
                audio.src = '/sounds/success.mp3'
                break
            default:
                audio.src = '/sounds/notification.mp3'
        }
        audio.play().catch(() => { }) // Ignore errors
    }

    // Get icon for notification type
    const getNotificationIcon = (type: string) => {
        switch (type) {
            case 'weather_alert':
                return CloudRain
            case 'delay_alert':
                return Clock
            case 'itinerary_update':
                return MapPin
            case 'trip_reminder':
                return Calendar
            case 'booking_confirm':
                return CheckCircle
            default:
                return Info
        }
    }

    // Get color scheme for notification type
    const getNotificationColors = (type: string, priority: string) => {
        const base = {
            'weather_alert': 'bg-amber-50 border-amber-200 text-amber-800',
            'delay_alert': 'bg-red-50 border-red-200 text-red-800',
            'itinerary_update': 'bg-blue-50 border-blue-200 text-blue-800',
            'trip_reminder': 'bg-green-50 border-green-200 text-green-800',
            'booking_confirm': 'bg-emerald-50 border-emerald-200 text-emerald-800',
            'general': 'bg-gray-50 border-gray-200 text-gray-800'
        }

        if (priority === 'high') {
            return 'bg-red-50 border-red-300 text-red-900'
        }

        return base[type as keyof typeof base] || base.general
    }

    // Handle notification click
    const handleNotificationClick = (notification: Notification) => {
        // Mark as read
        setNotifications(prev =>
            prev.map(n => n.id === notification.id ? { ...n, read: true } : n)
        )
        setUnreadCount(prev => Math.max(0, prev - 1))

        // Call action handler
        if (onNotificationAction) {
            onNotificationAction(notification)
        }

        // Navigate to action URL if available
        if (notification.actionUrl) {
            window.open(notification.actionUrl, '_blank')
        }
    }

    // Mark notification as read
    const markAsRead = (notificationId: string) => {
        setNotifications(prev =>
            prev.map(n => n.id === notificationId ? { ...n, read: true } : n)
        )
        setUnreadCount(prev => Math.max(0, prev - 1))
    }

    // Mark all notifications as read
    const markAllAsRead = () => {
        setNotifications(prev => prev.map(n => ({ ...n, read: true })))
        setUnreadCount(0)
    }

    // Clear notification
    const clearNotification = (notificationId: string) => {
        setNotifications(prev => prev.filter(n => n.id !== notificationId))
        setUnreadCount(prev => {
            const notification = notifications.find(n => n.id === notificationId)
            return notification && !notification.read ? Math.max(0, prev - 1) : prev
        })
    }

    // Clear all notifications
    const clearAllNotifications = () => {
        setNotifications([])
        setUnreadCount(0)
    }

    // Toggle notification setting
    const toggleSetting = (key: keyof NotificationSettings) => {
        setSettings(prev => ({ ...prev, [key]: !prev[key] }))
    }

    return (
        <div className={`relative ${className}`}>
            {/* Notification Bell */}
            <button
                onClick={() => setIsOpen(!isOpen)}
                className="relative p-3 text-gray-600 hover:text-gray-800 hover:bg-gray-100 rounded-full transition-colors"
                aria-label="Notifications"
            >
                <Bell className="w-6 h-6" />
                {unreadCount > 0 && (
                    <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                        {unreadCount > 99 ? '99+' : unreadCount}
                    </span>
                )}
                {!isOnline && (
                    <WifiOff className="absolute -bottom-1 -right-1 w-3 h-3 text-red-500" />
                )}
            </button>

            {/* Notification Panel */}
            <AnimatePresence>
                {isOpen && (
                    <motion.div
                        initial={{ opacity: 0, y: -10, scale: 0.95 }}
                        animate={{ opacity: 1, y: 0, scale: 1 }}
                        exit={{ opacity: 0, y: -10, scale: 0.95 }}
                        className="absolute right-0 top-full mt-2 w-96 bg-white rounded-lg shadow-xl border border-gray-200 z-50"
                    >
                        {/* Header */}
                        <div className="p-4 border-b border-gray-200 flex items-center justify-between">
                            <div className="flex items-center space-x-2">
                                <h3 className="font-semibold text-gray-900">Notifications</h3>
                                {!isOnline && (
                                    <span className="text-xs text-red-500 flex items-center">
                                        <WifiOff className="w-3 h-3 mr-1" />
                                        Offline
                                    </span>
                                )}
                            </div>
                            <div className="flex items-center space-x-2">
                                <button
                                    onClick={() => setShowSettings(!showSettings)}
                                    className="text-gray-500 hover:text-gray-700"
                                    aria-label="Notification Settings"
                                >
                                    <Settings className="w-4 h-4" />
                                </button>
                                <button
                                    onClick={() => setIsOpen(false)}
                                    className="text-gray-500 hover:text-gray-700"
                                    aria-label="Close"
                                >
                                    <X className="w-4 h-4" />
                                </button>
                            </div>
                        </div>

                        {/* Settings Panel */}
                        <AnimatePresence>
                            {showSettings && (
                                <motion.div
                                    initial={{ opacity: 0, height: 0 }}
                                    animate={{ opacity: 1, height: 'auto' }}
                                    exit={{ opacity: 0, height: 0 }}
                                    className="border-b border-gray-200 p-4 bg-gray-50"
                                >
                                    <h4 className="font-medium text-gray-900 mb-3">Notification Settings</h4>
                                    <div className="space-y-2">
                                        <label className="flex items-center">
                                            <input
                                                type="checkbox"
                                                checked={settings.enabled}
                                                onChange={() => toggleSetting('enabled')}
                                                className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                            />
                                            <span className="ml-2 text-sm text-gray-700">Enable notifications</span>
                                        </label>
                                        <label className="flex items-center">
                                            <input
                                                type="checkbox"
                                                checked={settings.sound}
                                                onChange={() => toggleSetting('sound')}
                                                className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                            />
                                            <span className="ml-2 text-sm text-gray-700">Play sound</span>
                                        </label>
                                        <label className="flex items-center">
                                            <input
                                                type="checkbox"
                                                checked={settings.weatherAlerts}
                                                onChange={() => toggleSetting('weatherAlerts')}
                                                className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                            />
                                            <span className="ml-2 text-sm text-gray-700">Weather alerts</span>
                                        </label>
                                        <label className="flex items-center">
                                            <input
                                                type="checkbox"
                                                checked={settings.delayAlerts}
                                                onChange={() => toggleSetting('delayAlerts')}
                                                className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                            />
                                            <span className="ml-2 text-sm text-gray-700">Delay alerts</span>
                                        </label>
                                        <label className="flex items-center">
                                            <input
                                                type="checkbox"
                                                checked={settings.itineraryUpdates}
                                                onChange={() => toggleSetting('itineraryUpdates')}
                                                className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                            />
                                            <span className="ml-2 text-sm text-gray-700">Itinerary updates</span>
                                        </label>
                                    </div>
                                </motion.div>
                            )}
                        </AnimatePresence>

                        {/* Actions */}
                        {notifications.length > 0 && (
                            <div className="p-2 border-b border-gray-200 flex justify-between">
                                <button
                                    onClick={markAllAsRead}
                                    className="text-sm text-blue-600 hover:text-blue-800"
                                >
                                    Mark all read
                                </button>
                                <button
                                    onClick={clearAllNotifications}
                                    className="text-sm text-red-600 hover:text-red-800"
                                >
                                    Clear all
                                </button>
                            </div>
                        )}

                        {/* Notifications List */}
                        <div className="max-h-96 overflow-y-auto">
                            {notifications.length === 0 ? (
                                <div className="p-8 text-center text-gray-500">
                                    <Bell className="w-8 h-8 mx-auto mb-2 opacity-50" />
                                    <p>No notifications yet</p>
                                </div>
                            ) : (
                                notifications.map((notification) => {
                                    const IconComponent = getNotificationIcon(notification.type)
                                    return (
                                        <motion.div
                                            key={notification.id}
                                            initial={{ opacity: 0, x: -20 }}
                                            animate={{ opacity: 1, x: 0 }}
                                            exit={{ opacity: 0, x: 20 }}
                                            className={`p-4 border-b border-gray-100 hover:bg-gray-50 cursor-pointer transition-colors ${!notification.read ? 'bg-blue-50' : ''
                                                }`}
                                            onClick={() => handleNotificationClick(notification)}
                                        >
                                            <div className="flex items-start space-x-3">
                                                <div className={`p-1 rounded-full ${getNotificationColors(notification.type, notification.priority)}`}>
                                                    <IconComponent className="w-4 h-4" />
                                                </div>
                                                <div className="flex-1 min-w-0">
                                                    <div className="flex items-center justify-between">
                                                        <h4 className={`text-sm font-medium truncate ${!notification.read ? 'text-gray-900' : 'text-gray-700'}`}>
                                                            {notification.title}
                                                        </h4>
                                                        <div className="flex items-center space-x-1">
                                                            {notification.actionUrl && (
                                                                <ExternalLink className="w-3 h-3 text-gray-400" />
                                                            )}
                                                            <button
                                                                onClick={(e) => {
                                                                    e.stopPropagation()
                                                                    clearNotification(notification.id)
                                                                }}
                                                                className="text-gray-400 hover:text-gray-600"
                                                            >
                                                                <X className="w-3 h-3" />
                                                            </button>
                                                        </div>
                                                    </div>
                                                    <p className="text-sm text-gray-600 line-clamp-2">
                                                        {notification.body}
                                                    </p>
                                                    <div className="flex items-center justify-between mt-1">
                                                        <span className="text-xs text-gray-500">
                                                            {notification.timestamp.toLocaleTimeString()}
                                                        </span>
                                                        {notification.priority === 'high' && (
                                                            <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-red-100 text-red-800">
                                                                <AlertTriangle className="w-3 h-3 mr-1" />
                                                                High Priority
                                                            </span>
                                                        )}
                                                    </div>
                                                </div>
                                            </div>
                                        </motion.div>
                                    )
                                })
                            )}
                        </div>
                    </motion.div>
                )}
            </AnimatePresence>
        </div>
    )
}