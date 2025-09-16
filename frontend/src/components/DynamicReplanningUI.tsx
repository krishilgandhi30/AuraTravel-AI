'use client'

import React, { useState, useEffect, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import {
    RefreshCw,
    AlertTriangle,
    CheckCircle,
    Clock,
    CloudRain,
    Sun,
    MapPin,
    Calendar,
    DollarSign,
    Users,
    ArrowRight,
    ArrowLeft,
    X,
    Settings,
    Zap,
    Activity,
    TrendingUp,
    Target,
    Sparkles,
    Play,
    Pause,
    SkipForward
} from 'lucide-react'

// Types for dynamic replanning
interface ReplanningTrigger {
    type: 'weather' | 'delay' | 'availability' | 'user_preference'
    severity: 'low' | 'medium' | 'high' | 'critical'
    description: string
    affectedActivities: string[]
    timestamp: Date
    location?: string
    weatherCondition?: string
    delayMinutes?: number
}

interface ReplanningOption {
    id: string
    title: string
    description: string
    impact: 'minimal' | 'moderate' | 'major'
    costChange: number
    timeChange: number
    confidence: number
    changes: ItineraryChange[]
    pros: string[]
    cons: string[]
    aiRecommendation: boolean
}

interface ItineraryChange {
    type: 'add' | 'remove' | 'modify' | 'reschedule'
    activity: string
    originalTime?: string
    newTime?: string
    reason: string
    location?: string
}

interface DynamicReplanningUIProps {
    tripId: string
    currentItinerary?: any
    locale?: string
    onReplanAccept?: (option: ReplanningOption) => void
    onReplanReject?: (reason: string) => void
    className?: string
}

export default function DynamicReplanningUI({
    tripId,
    currentItinerary,
    locale = 'en',
    onReplanAccept,
    onReplanReject,
    className = ''
}: DynamicReplanningUIProps) {
    const [isMonitoring, setIsMonitoring] = useState(true)
    const [triggers, setTriggers] = useState<ReplanningTrigger[]>([])
    const [currentReplanning, setCurrentReplanning] = useState<{
        trigger: ReplanningTrigger
        options: ReplanningOption[]
        loading: boolean
    } | null>(null)
    const [showHistory, setShowHistory] = useState(false)
    const [replanningHistory, setReplanningHistory] = useState<Array<{
        id: string
        timestamp: Date
        trigger: ReplanningTrigger
        selectedOption?: ReplanningOption
        status: 'accepted' | 'rejected' | 'auto-applied'
    }>>([])
    const [autoReplanEnabled, setAutoReplanEnabled] = useState(false)
    const [lastMonitorUpdate, setLastMonitorUpdate] = useState(new Date())

    // Start monitoring for replanning triggers
    useEffect(() => {
        if (!isMonitoring || !tripId) return

        const monitoringInterval = setInterval(async () => {
            try {
                const response = await fetch(`/api/trips/${tripId}/monitoring-status`)
                const data = await response.json()

                if (data.newTriggers && data.newTriggers.length > 0) {
                    setTriggers(prev => [...data.newTriggers, ...prev])

                    // Process critical triggers immediately
                    const criticalTrigger = data.newTriggers.find((t: ReplanningTrigger) =>
                        t.severity === 'critical'
                    )

                    if (criticalTrigger) {
                        handleTriggerReplanning(criticalTrigger)
                    }
                }

                setLastMonitorUpdate(new Date())
            } catch (error) {
                console.error('Error checking monitoring status:', error)
            }
        }, 15000) // Check every 15 seconds

        return () => clearInterval(monitoringInterval)
    }, [isMonitoring, tripId])

    // Handle replanning trigger
    const handleTriggerReplanning = async (trigger: ReplanningTrigger) => {
        if (currentReplanning) return // Already processing a replanning

        setCurrentReplanning({
            trigger,
            options: [],
            loading: true
        })

        try {
            const response = await fetch('/api/trips/dynamic-replan', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    tripId,
                    trigger,
                    currentItinerary,
                    locale
                })
            })

            const data = await response.json()

            setCurrentReplanning(prev => prev ? {
                ...prev,
                options: data.options || [],
                loading: false
            } : null)

            // Auto-apply if enabled and high confidence
            if (autoReplanEnabled && data.options?.length > 0) {
                const bestOption = data.options.find((opt: ReplanningOption) =>
                    opt.aiRecommendation && opt.confidence > 0.8
                )

                if (bestOption) {
                    setTimeout(() => handleAcceptReplanning(bestOption, true), 3000)
                }
            }
        } catch (error) {
            console.error('Error generating replanning options:', error)
            setCurrentReplanning(null)
        }
    }

    // Accept replanning option
    const handleAcceptReplanning = (option: ReplanningOption, autoApplied = false) => {
        if (!currentReplanning) return

        // Add to history
        setReplanningHistory(prev => [{
            id: Date.now().toString(),
            timestamp: new Date(),
            trigger: currentReplanning.trigger,
            selectedOption: option,
            status: autoApplied ? 'auto-applied' : 'accepted'
        }, ...prev])

        // Clear current replanning
        setCurrentReplanning(null)

        // Call callback
        if (onReplanAccept) {
            onReplanAccept(option)
        }
    }

    // Reject replanning
    const handleRejectReplanning = (reason: string) => {
        if (!currentReplanning) return

        // Add to history
        setReplanningHistory(prev => [{
            id: Date.now().toString(),
            timestamp: new Date(),
            trigger: currentReplanning.trigger,
            status: 'rejected'
        }, ...prev])

        // Clear current replanning
        setCurrentReplanning(null)

        // Call callback
        if (onReplanReject) {
            onReplanReject(reason)
        }
    }

    // Get trigger icon
    const getTriggerIcon = (type: string) => {
        switch (type) {
            case 'weather':
                return CloudRain
            case 'delay':
                return Clock
            case 'availability':
                return MapPin
            case 'user_preference':
                return Users
            default:
                return AlertTriangle
        }
    }

    // Get severity color
    const getSeverityColor = (severity: string) => {
        switch (severity) {
            case 'critical':
                return 'bg-red-500 text-white'
            case 'high':
                return 'bg-orange-500 text-white'
            case 'medium':
                return 'bg-yellow-500 text-white'
            case 'low':
                return 'bg-blue-500 text-white'
            default:
                return 'bg-gray-500 text-white'
        }
    }

    // Get impact color
    const getImpactColor = (impact: string) => {
        switch (impact) {
            case 'major':
                return 'text-red-600 bg-red-50'
            case 'moderate':
                return 'text-orange-600 bg-orange-50'
            case 'minimal':
                return 'text-green-600 bg-green-50'
            default:
                return 'text-gray-600 bg-gray-50'
        }
    }

    return (
        <div className={`bg-white rounded-lg shadow-lg border border-gray-200 ${className}`}>
            {/* Header */}
            <div className="p-4 border-b border-gray-200">
                <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-3">
                        <div className={`p-2 rounded-full ${isMonitoring ? 'bg-green-100 text-green-600' : 'bg-gray-100 text-gray-600'}`}>
                            <Activity className={`w-5 h-5 ${isMonitoring ? 'animate-pulse' : ''}`} />
                        </div>
                        <div>
                            <h3 className="text-lg font-semibold text-gray-900">Dynamic Replanning</h3>
                            <p className="text-sm text-gray-500">
                                {isMonitoring ? (
                                    <>
                                        <span className="text-green-600">● Active</span>
                                        {' - Last check: '}
                                        {lastMonitorUpdate.toLocaleTimeString()}
                                    </>
                                ) : (
                                    <span className="text-gray-600">● Paused</span>
                                )}
                            </p>
                        </div>
                    </div>
                    <div className="flex items-center space-x-2">
                        <button
                            onClick={() => setShowHistory(!showHistory)}
                            className="px-3 py-1 text-sm text-gray-600 hover:text-gray-800 hover:bg-gray-100 rounded"
                        >
                            History
                        </button>
                        <button
                            onClick={() => setIsMonitoring(!isMonitoring)}
                            className={`p-2 rounded-full ${isMonitoring ? 'bg-red-100 text-red-600 hover:bg-red-200' : 'bg-green-100 text-green-600 hover:bg-green-200'}`}
                        >
                            {isMonitoring ? <Pause className="w-4 h-4" /> : <Play className="w-4 h-4" />}
                        </button>
                    </div>
                </div>

                {/* Auto-replan toggle */}
                <div className="mt-3 flex items-center justify-between">
                    <label className="flex items-center space-x-2">
                        <input
                            type="checkbox"
                            checked={autoReplanEnabled}
                            onChange={(e) => setAutoReplanEnabled(e.target.checked)}
                            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                        />
                        <span className="text-sm text-gray-700">Auto-apply high-confidence changes</span>
                    </label>
                    {triggers.length > 0 && (
                        <span className="text-sm text-gray-500">
                            {triggers.length} active trigger{triggers.length !== 1 ? 's' : ''}
                        </span>
                    )}
                </div>
            </div>

            {/* Current Replanning Modal */}
            <AnimatePresence>
                {currentReplanning && (
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
                    >
                        <motion.div
                            initial={{ scale: 0.9, opacity: 0 }}
                            animate={{ scale: 1, opacity: 1 }}
                            exit={{ scale: 0.9, opacity: 0 }}
                            className="bg-white rounded-lg shadow-xl max-w-4xl w-full mx-4 max-h-[90vh] overflow-y-auto"
                        >
                            {/* Modal Header */}
                            <div className="p-6 border-b border-gray-200">
                                <div className="flex items-center justify-between">
                                    <div className="flex items-center space-x-3">
                                        <div className={`p-2 rounded-full ${getSeverityColor(currentReplanning.trigger.severity)}`}>
                                            {React.createElement(getTriggerIcon(currentReplanning.trigger.type), {
                                                className: 'w-5 h-5'
                                            })}
                                        </div>
                                        <div>
                                            <h3 className="text-xl font-semibold text-gray-900">
                                                Replanning Required
                                            </h3>
                                            <p className="text-gray-600">
                                                {currentReplanning.trigger.description}
                                            </p>
                                        </div>
                                    </div>
                                    <button
                                        onClick={() => setCurrentReplanning(null)}
                                        className="text-gray-500 hover:text-gray-700"
                                    >
                                        <X className="w-6 h-6" />
                                    </button>
                                </div>

                                {/* Trigger Details */}
                                <div className="mt-4 p-4 bg-gray-50 rounded-lg">
                                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                                        <div>
                                            <span className="font-medium text-gray-700">Type:</span>
                                            <p className="text-gray-900 capitalize">{currentReplanning.trigger.type}</p>
                                        </div>
                                        <div>
                                            <span className="font-medium text-gray-700">Severity:</span>
                                            <p className={`capitalize ${currentReplanning.trigger.severity === 'critical' ? 'text-red-600' :
                                                    currentReplanning.trigger.severity === 'high' ? 'text-orange-600' :
                                                        currentReplanning.trigger.severity === 'medium' ? 'text-yellow-600' :
                                                            'text-blue-600'
                                                }`}>
                                                {currentReplanning.trigger.severity}
                                            </p>
                                        </div>
                                        <div>
                                            <span className="font-medium text-gray-700">Time:</span>
                                            <p className="text-gray-900">
                                                {currentReplanning.trigger.timestamp.toLocaleTimeString()}
                                            </p>
                                        </div>
                                        <div>
                                            <span className="font-medium text-gray-700">Affected:</span>
                                            <p className="text-gray-900">
                                                {currentReplanning.trigger.affectedActivities.length} activities
                                            </p>
                                        </div>
                                    </div>
                                </div>
                            </div>

                            {/* Replanning Options */}
                            <div className="p-6">
                                {currentReplanning.loading ? (
                                    <div className="flex items-center justify-center py-12">
                                        <div className="text-center">
                                            <RefreshCw className="w-8 h-8 animate-spin text-blue-600 mx-auto mb-4" />
                                            <p className="text-gray-600">Generating replanning options...</p>
                                        </div>
                                    </div>
                                ) : (
                                    <div className="space-y-4">
                                        <h4 className="text-lg font-semibold text-gray-900">
                                            Recommended Options
                                        </h4>
                                        {currentReplanning.options.map((option, index) => (
                                            <motion.div
                                                key={option.id}
                                                initial={{ opacity: 0, y: 20 }}
                                                animate={{ opacity: 1, y: 0 }}
                                                transition={{ delay: index * 0.1 }}
                                                className={`border rounded-lg p-4 hover:shadow-md transition-shadow ${option.aiRecommendation ? 'border-blue-200 bg-blue-50' : 'border-gray-200'
                                                    }`}
                                            >
                                                <div className="flex items-start justify-between">
                                                    <div className="flex-1">
                                                        <div className="flex items-center space-x-2 mb-2">
                                                            <h5 className="font-semibold text-gray-900">
                                                                {option.title}
                                                            </h5>
                                                            {option.aiRecommendation && (
                                                                <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                                                                    <Sparkles className="w-3 h-3 mr-1" />
                                                                    AI Recommended
                                                                </span>
                                                            )}
                                                            <span className={`px-2 py-1 rounded-full text-xs ${getImpactColor(option.impact)}`}>
                                                                {option.impact} impact
                                                            </span>
                                                        </div>
                                                        <p className="text-gray-600 mb-3">
                                                            {option.description}
                                                        </p>

                                                        {/* Changes Summary */}
                                                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                                                            <div className="text-center p-3 bg-white rounded border">
                                                                <DollarSign className="w-5 h-5 mx-auto mb-1 text-gray-600" />
                                                                <p className="text-sm font-medium">
                                                                    {option.costChange >= 0 ? '+' : ''}
                                                                    ₹{Math.abs(option.costChange)}
                                                                </p>
                                                                <p className="text-xs text-gray-500">Cost change</p>
                                                            </div>
                                                            <div className="text-center p-3 bg-white rounded border">
                                                                <Clock className="w-5 h-5 mx-auto mb-1 text-gray-600" />
                                                                <p className="text-sm font-medium">
                                                                    {option.timeChange >= 0 ? '+' : ''}
                                                                    {Math.abs(option.timeChange)}h
                                                                </p>
                                                                <p className="text-xs text-gray-500">Time change</p>
                                                            </div>
                                                            <div className="text-center p-3 bg-white rounded border">
                                                                <TrendingUp className="w-5 h-5 mx-auto mb-1 text-gray-600" />
                                                                <p className="text-sm font-medium">
                                                                    {Math.round(option.confidence * 100)}%
                                                                </p>
                                                                <p className="text-xs text-gray-500">Confidence</p>
                                                            </div>
                                                        </div>

                                                        {/* Pros and Cons */}
                                                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                                                            <div>
                                                                <h6 className="font-medium text-green-700 mb-2">Pros:</h6>
                                                                <ul className="text-sm text-gray-600 space-y-1">
                                                                    {option.pros.map((pro, i) => (
                                                                        <li key={i} className="flex items-start">
                                                                            <CheckCircle className="w-3 h-3 text-green-500 mt-0.5 mr-2 flex-shrink-0" />
                                                                            {pro}
                                                                        </li>
                                                                    ))}
                                                                </ul>
                                                            </div>
                                                            <div>
                                                                <h6 className="font-medium text-orange-700 mb-2">Considerations:</h6>
                                                                <ul className="text-sm text-gray-600 space-y-1">
                                                                    {option.cons.map((con, i) => (
                                                                        <li key={i} className="flex items-start">
                                                                            <AlertTriangle className="w-3 h-3 text-orange-500 mt-0.5 mr-2 flex-shrink-0" />
                                                                            {con}
                                                                        </li>
                                                                    ))}
                                                                </ul>
                                                            </div>
                                                        </div>

                                                        {/* Detailed Changes */}
                                                        <div className="border-t pt-3">
                                                            <h6 className="font-medium text-gray-700 mb-2">Changes:</h6>
                                                            <div className="space-y-2">
                                                                {option.changes.map((change, i) => (
                                                                    <div key={i} className="flex items-center text-sm">
                                                                        <span className={`px-2 py-1 rounded text-xs mr-2 ${change.type === 'add' ? 'bg-green-100 text-green-800' :
                                                                                change.type === 'remove' ? 'bg-red-100 text-red-800' :
                                                                                    change.type === 'modify' ? 'bg-blue-100 text-blue-800' :
                                                                                        'bg-yellow-100 text-yellow-800'
                                                                            }`}>
                                                                            {change.type}
                                                                        </span>
                                                                        <span className="text-gray-900">{change.activity}</span>
                                                                        {change.originalTime && change.newTime && (
                                                                            <span className="text-gray-500 ml-2">
                                                                                {change.originalTime} → {change.newTime}
                                                                            </span>
                                                                        )}
                                                                    </div>
                                                                ))}
                                                            </div>
                                                        </div>
                                                    </div>

                                                    {/* Action Buttons */}
                                                    <div className="ml-4 flex flex-col space-y-2">
                                                        <button
                                                            onClick={() => handleAcceptReplanning(option)}
                                                            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors flex items-center space-x-1"
                                                        >
                                                            <CheckCircle className="w-4 h-4" />
                                                            <span>Accept</span>
                                                        </button>
                                                        <button
                                                            onClick={() => handleRejectReplanning('User preferred alternative')}
                                                            className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 transition-colors"
                                                        >
                                                            Skip
                                                        </button>
                                                    </div>
                                                </div>
                                            </motion.div>
                                        ))}

                                        {/* Reject All Button */}
                                        <div className="flex justify-center pt-4 border-t">
                                            <button
                                                onClick={() => handleRejectReplanning('User rejected all options')}
                                                className="px-6 py-2 text-red-600 hover:text-red-800 font-medium"
                                            >
                                                Keep current itinerary
                                            </button>
                                        </div>
                                    </div>
                                )}
                            </div>
                        </motion.div>
                    </motion.div>
                )}
            </AnimatePresence>

            {/* History Panel */}
            <AnimatePresence>
                {showHistory && (
                    <motion.div
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: 'auto' }}
                        exit={{ opacity: 0, height: 0 }}
                        className="border-t border-gray-200 p-4 bg-gray-50"
                    >
                        <h4 className="font-medium text-gray-900 mb-3">Replanning History</h4>
                        {replanningHistory.length === 0 ? (
                            <p className="text-gray-500 text-center py-4">No replanning history yet</p>
                        ) : (
                            <div className="space-y-3 max-h-60 overflow-y-auto">
                                {replanningHistory.map((entry) => (
                                    <div key={entry.id} className="flex items-center justify-between p-3 bg-white rounded border">
                                        <div className="flex items-center space-x-3">
                                            <div className={`p-1 rounded-full ${getSeverityColor(entry.trigger.severity)}`}>
                                                {React.createElement(getTriggerIcon(entry.trigger.type), {
                                                    className: 'w-3 h-3'
                                                })}
                                            </div>
                                            <div>
                                                <p className="text-sm font-medium text-gray-900">
                                                    {entry.trigger.description}
                                                </p>
                                                <p className="text-xs text-gray-500">
                                                    {entry.timestamp.toLocaleString()}
                                                </p>
                                            </div>
                                        </div>
                                        <span className={`px-2 py-1 rounded text-xs ${entry.status === 'accepted' ? 'bg-green-100 text-green-800' :
                                                entry.status === 'auto-applied' ? 'bg-blue-100 text-blue-800' :
                                                    'bg-red-100 text-red-800'
                                            }`}>
                                            {entry.status}
                                        </span>
                                    </div>
                                ))}
                            </div>
                        )}
                    </motion.div>
                )}
            </AnimatePresence>

            {/* Recent Triggers */}
            {triggers.length > 0 && !currentReplanning && (
                <div className="p-4 border-t border-gray-200">
                    <h4 className="font-medium text-gray-900 mb-3">Recent Triggers</h4>
                    <div className="space-y-2">
                        {triggers.slice(0, 3).map((trigger, index) => (
                            <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded">
                                <div className="flex items-center space-x-3">
                                    <div className={`p-1 rounded-full ${getSeverityColor(trigger.severity)}`}>
                                        {React.createElement(getTriggerIcon(trigger.type), {
                                            className: 'w-3 h-3'
                                        })}
                                    </div>
                                    <div>
                                        <p className="text-sm font-medium text-gray-900">
                                            {trigger.description}
                                        </p>
                                        <p className="text-xs text-gray-500">
                                            {trigger.timestamp.toLocaleTimeString()}
                                        </p>
                                    </div>
                                </div>
                                <button
                                    onClick={() => handleTriggerReplanning(trigger)}
                                    className="px-3 py-1 text-sm text-blue-600 hover:text-blue-800 font-medium"
                                >
                                    Review
                                </button>
                            </div>
                        ))}
                    </div>
                </div>
            )}
        </div>
    )
}