'use client'

import React, { useState, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import {
    Download,
    Mail,
    MessageSquare,
    Share2,
    FileText,
    Calendar,
    Code,
    Globe,
    Check,
    AlertCircle,
    Clock,
    User,
    Phone,
    Copy,
    ExternalLink,
    Smartphone,
    Printer,
    QrCode,
    Cloud,
    Archive
} from 'lucide-react'

// Types for itinerary delivery
interface DeliveryOption {
    type: 'download' | 'email' | 'sms' | 'share'
    format?: 'pdf' | 'ics' | 'json' | 'html'
    label: string
    icon: any
    description: string
    available: boolean
}

interface DeliveryRequest {
    tripId: string
    format: string
    deliveryMethod: string
    recipient?: string
    phoneNumber?: string
    email?: string
    locale?: string
    includeTickets?: boolean
    includeWeather?: boolean
    includeContacts?: boolean
}

interface DeliveryResult {
    success: boolean
    downloadUrl?: string
    message: string
    trackingId?: string
    deliveryTime?: Date
}

interface ItineraryDeliveryProps {
    tripId: string
    tripTitle: string
    currentItinerary?: any
    userEmail?: string
    userPhone?: string
    locale?: string
    className?: string
}

export default function ItineraryDelivery({
    tripId,
    tripTitle,
    currentItinerary,
    userEmail,
    userPhone,
    locale = 'en',
    className = ''
}: ItineraryDeliveryProps) {
    const [selectedFormat, setSelectedFormat] = useState<'pdf' | 'ics' | 'json' | 'html'>('pdf')
    const [deliveryStatus, setDeliveryStatus] = useState<{
        [key: string]: 'idle' | 'processing' | 'success' | 'error'
    }>({})
    const [deliveryResults, setDeliveryResults] = useState<{
        [key: string]: DeliveryResult
    }>({})
    const [emailRecipient, setEmailRecipient] = useState(userEmail || '')
    const [phoneRecipient, setPhoneRecipient] = useState(userPhone || '')
    const [showAdvancedOptions, setShowAdvancedOptions] = useState(false)
    const [deliveryOptions, setDeliveryOptions] = useState({
        includeTickets: true,
        includeWeather: true,
        includeContacts: true,
        includeMap: true,
        includePhotos: false,
        compressFiles: false
    })
    const [shareUrl, setShareUrl] = useState('')
    const [qrCode, setQrCode] = useState('')

    // Available delivery formats
    const formats = [
        {
            type: 'pdf',
            label: 'PDF Document',
            icon: FileText,
            description: 'Professional printable document with full itinerary details',
            features: ['Printable', 'Offline viewing', 'Professional layout', 'Includes maps']
        },
        {
            type: 'ics',
            label: 'Calendar File',
            icon: Calendar,
            description: 'Import directly into your calendar application',
            features: ['Calendar integration', 'Reminders', 'Location data', 'Time zones']
        },
        {
            type: 'json',
            label: 'JSON Data',
            icon: Code,
            description: 'Raw data format for developers and integrations',
            features: ['Machine readable', 'API integration', 'Custom processing', 'Complete data']
        },
        {
            type: 'html',
            label: 'Web Page',
            icon: Globe,
            description: 'Interactive web page with live updates',
            features: ['Interactive maps', 'Real-time updates', 'Mobile friendly', 'Shareable link']
        }
    ]

    // Generate and deliver itinerary
    const handleDelivery = async (method: string, format: string = selectedFormat) => {
        const key = `${method}_${format}`
        setDeliveryStatus(prev => ({ ...prev, [key]: 'processing' }))

        try {
            const request: DeliveryRequest = {
                tripId,
                format,
                deliveryMethod: method,
                locale,
                ...deliveryOptions
            }

            // Add recipient information based on method
            if (method === 'email') {
                request.email = emailRecipient
            } else if (method === 'sms') {
                request.phoneNumber = phoneRecipient
            }

            const response = await fetch('/api/trips/deliver', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(request)
            })

            const result = await response.json()

            if (result.success) {
                setDeliveryStatus(prev => ({ ...prev, [key]: 'success' }))
                setDeliveryResults(prev => ({ ...prev, [key]: result }))

                // Auto-download for download requests
                if (method === 'download' && result.downloadUrl) {
                    const link = document.createElement('a')
                    link.href = result.downloadUrl
                    link.download = `${tripTitle}_itinerary.${format}`
                    document.body.appendChild(link)
                    link.click()
                    document.body.removeChild(link)
                }

                // Generate share URL for HTML format
                if (format === 'html' && result.downloadUrl) {
                    setShareUrl(result.downloadUrl)
                    generateQRCode(result.downloadUrl)
                }
            } else {
                throw new Error(result.message || 'Delivery failed')
            }
        } catch (error) {
            console.error('Delivery error:', error)
            setDeliveryStatus(prev => ({ ...prev, [key]: 'error' }))
            setDeliveryResults(prev => ({
                ...prev,
                [key]: {
                    success: false,
                    message: error instanceof Error ? error.message : 'Unknown error'
                }
            }))
        }
    }

    // Generate QR code for sharing
    const generateQRCode = async (url: string) => {
        try {
            const response = await fetch('/api/qr-code', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ url, size: 200 })
            })

            if (response.ok) {
                const blob = await response.blob()
                const qrUrl = URL.createObjectURL(blob)
                setQrCode(qrUrl)
            }
        } catch (error) {
            console.error('QR code generation error:', error)
        }
    }

    // Copy share URL to clipboard
    const copyShareUrl = useCallback(async () => {
        if (shareUrl) {
            try {
                await navigator.clipboard.writeText(shareUrl)
                // Show temporary success message
                const key = 'copy_url'
                setDeliveryStatus(prev => ({ ...prev, [key]: 'success' }))
                setTimeout(() => {
                    setDeliveryStatus(prev => ({ ...prev, [key]: 'idle' }))
                }, 2000)
            } catch (error) {
                console.error('Copy failed:', error)
            }
        }
    }, [shareUrl])

    // Get status icon
    const getStatusIcon = (status: string) => {
        switch (status) {
            case 'processing':
                return <Clock className="w-4 h-4 animate-spin" />
            case 'success':
                return <Check className="w-4 h-4 text-green-600" />
            case 'error':
                return <AlertCircle className="w-4 h-4 text-red-600" />
            default:
                return null
        }
    }

    // Get button state
    const getButtonState = (method: string, format: string = selectedFormat) => {
        const key = `${method}_${format}`
        const status = deliveryStatus[key]
        const result = deliveryResults[key]

        return {
            status,
            result,
            disabled: status === 'processing',
            className: status === 'success' ? 'bg-green-600 hover:bg-green-700' :
                status === 'error' ? 'bg-red-600 hover:bg-red-700' :
                    'bg-blue-600 hover:bg-blue-700'
        }
    }

    return (
        <div className={`bg-white rounded-lg shadow-lg border border-gray-200 ${className}`}>
            {/* Header */}
            <div className="p-6 border-b border-gray-200">
                <div className="flex items-center justify-between">
                    <div>
                        <h3 className="text-lg font-semibold text-gray-900 flex items-center">
                            <Download className="w-5 h-5 mr-2" />
                            Download & Share Itinerary
                        </h3>
                        <p className="text-gray-600 mt-1">
                            Get your itinerary in multiple formats and share with others
                        </p>
                    </div>
                    <button
                        onClick={() => setShowAdvancedOptions(!showAdvancedOptions)}
                        className="text-sm text-blue-600 hover:text-blue-800"
                    >
                        Advanced Options
                    </button>
                </div>
            </div>

            {/* Format Selection */}
            <div className="p-6 border-b border-gray-200">
                <h4 className="font-medium text-gray-900 mb-4">Choose Format</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                    {formats.map((format) => (
                        <motion.button
                            key={format.type}
                            onClick={() => setSelectedFormat(format.type as any)}
                            className={`p-4 border rounded-lg text-left transition-all ${selectedFormat === format.type
                                    ? 'border-blue-500 bg-blue-50 shadow-md'
                                    : 'border-gray-200 hover:border-gray-300 hover:shadow-sm'
                                }`}
                            whileHover={{ scale: 1.02 }}
                            whileTap={{ scale: 0.98 }}
                        >
                            <div className="flex items-center mb-2">
                                <format.icon className={`w-5 h-5 mr-2 ${selectedFormat === format.type ? 'text-blue-600' : 'text-gray-600'
                                    }`} />
                                <span className="font-medium text-gray-900">{format.label}</span>
                            </div>
                            <p className="text-sm text-gray-600 mb-3">{format.description}</p>
                            <div className="space-y-1">
                                {format.features.map((feature, index) => (
                                    <div key={index} className="flex items-center text-xs text-gray-500">
                                        <Check className="w-3 h-3 mr-1 text-green-500" />
                                        {feature}
                                    </div>
                                ))}
                            </div>
                        </motion.button>
                    ))}
                </div>
            </div>

            {/* Advanced Options */}
            <AnimatePresence>
                {showAdvancedOptions && (
                    <motion.div
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: 'auto' }}
                        exit={{ opacity: 0, height: 0 }}
                        className="border-b border-gray-200 p-6 bg-gray-50"
                    >
                        <h4 className="font-medium text-gray-900 mb-4">Advanced Options</h4>
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                            {Object.entries(deliveryOptions).map(([key, value]) => (
                                <label key={key} className="flex items-center">
                                    <input
                                        type="checkbox"
                                        checked={value}
                                        onChange={(e) => setDeliveryOptions(prev => ({
                                            ...prev,
                                            [key]: e.target.checked
                                        }))}
                                        className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                    />
                                    <span className="ml-2 text-sm text-gray-700 capitalize">
                                        {key.replace(/([A-Z])/g, ' $1').toLowerCase()}
                                    </span>
                                </label>
                            ))}
                        </div>
                    </motion.div>
                )}
            </AnimatePresence>

            {/* Delivery Actions */}
            <div className="p-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    {/* Download Section */}
                    <div>
                        <h4 className="font-medium text-gray-900 mb-4">Download</h4>
                        <div className="space-y-3">
                            <motion.button
                                onClick={() => handleDelivery('download')}
                                className={`w-full flex items-center justify-center px-4 py-3 text-white rounded-lg transition-colors ${getButtonState('download').className
                                    }`}
                                disabled={getButtonState('download').disabled}
                                whileHover={{ scale: 1.02 }}
                                whileTap={{ scale: 0.98 }}
                            >
                                {getStatusIcon(getButtonState('download').status) || <Download className="w-4 h-4 mr-2" />}
                                Download {selectedFormat.toUpperCase()}
                            </motion.button>

                            {/* Quick download buttons for other formats */}
                            <div className="grid grid-cols-2 gap-2">
                                {formats.filter(f => f.type !== selectedFormat).map(format => {
                                    const state = getButtonState('download', format.type)
                                    return (
                                        <button
                                            key={format.type}
                                            onClick={() => handleDelivery('download', format.type)}
                                            className="flex items-center justify-center px-3 py-2 text-sm text-gray-700 border border-gray-300 rounded hover:bg-gray-50 transition-colors"
                                            disabled={state.disabled}
                                        >
                                            {getStatusIcon(state.status) || <format.icon className="w-3 h-3 mr-1" />}
                                            {format.type.toUpperCase()}
                                        </button>
                                    )
                                })}
                            </div>

                            {/* Download results */}
                            {Object.entries(deliveryResults).filter(([key]) => key.startsWith('download')).map(([key, result]) => (
                                <div key={key} className="p-3 bg-gray-50 rounded text-sm">
                                    {result.success ? (
                                        <p className="text-green-700">✓ Download completed successfully</p>
                                    ) : (
                                        <p className="text-red-700">✗ {result.message}</p>
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* Send Section */}
                    <div>
                        <h4 className="font-medium text-gray-900 mb-4">Send & Share</h4>
                        <div className="space-y-4">
                            {/* Email */}
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-2">
                                    Email Address
                                </label>
                                <div className="flex space-x-2">
                                    <input
                                        type="email"
                                        value={emailRecipient}
                                        onChange={(e) => setEmailRecipient(e.target.value)}
                                        placeholder="Enter email address"
                                        className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                                    />
                                    <motion.button
                                        onClick={() => handleDelivery('email')}
                                        className={`px-4 py-2 text-white rounded-md transition-colors flex items-center ${getButtonState('email').className
                                            }`}
                                        disabled={getButtonState('email').disabled || !emailRecipient}
                                        whileHover={{ scale: 1.02 }}
                                        whileTap={{ scale: 0.98 }}
                                    >
                                        {getStatusIcon(getButtonState('email').status) || <Mail className="w-4 h-4 mr-1" />}
                                        Send
                                    </motion.button>
                                </div>
                                {deliveryResults.email_pdf && (
                                    <p className={`mt-2 text-sm ${deliveryResults.email_pdf.success ? 'text-green-600' : 'text-red-600'
                                        }`}>
                                        {deliveryResults.email_pdf.message}
                                    </p>
                                )}
                            </div>

                            {/* SMS */}
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-2">
                                    Phone Number
                                </label>
                                <div className="flex space-x-2">
                                    <input
                                        type="tel"
                                        value={phoneRecipient}
                                        onChange={(e) => setPhoneRecipient(e.target.value)}
                                        placeholder="+1 (555) 123-4567"
                                        className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                                    />
                                    <motion.button
                                        onClick={() => handleDelivery('sms')}
                                        className={`px-4 py-2 text-white rounded-md transition-colors flex items-center ${getButtonState('sms').className
                                            }`}
                                        disabled={getButtonState('sms').disabled || !phoneRecipient}
                                        whileHover={{ scale: 1.02 }}
                                        whileTap={{ scale: 0.98 }}
                                    >
                                        {getStatusIcon(getButtonState('sms').status) || <MessageSquare className="w-4 h-4 mr-1" />}
                                        Send
                                    </motion.button>
                                </div>
                                {deliveryResults.sms_pdf && (
                                    <p className={`mt-2 text-sm ${deliveryResults.sms_pdf.success ? 'text-green-600' : 'text-red-600'
                                        }`}>
                                        {deliveryResults.sms_pdf.message}
                                    </p>
                                )}
                            </div>

                            {/* Share Link */}
                            <div>
                                <div className="flex justify-between items-center mb-2">
                                    <label className="block text-sm font-medium text-gray-700">
                                        Share Link
                                    </label>
                                    <motion.button
                                        onClick={() => handleDelivery('share', 'html')}
                                        className="px-3 py-1 text-sm text-blue-600 hover:text-blue-800"
                                        disabled={getButtonState('share', 'html').disabled}
                                    >
                                        {getButtonState('share', 'html').status === 'processing' ? 'Generating...' : 'Generate Link'}
                                    </motion.button>
                                </div>

                                {shareUrl && (
                                    <div className="space-y-3">
                                        <div className="flex space-x-2">
                                            <input
                                                type="text"
                                                value={shareUrl}
                                                readOnly
                                                className="flex-1 px-3 py-2 border border-gray-300 rounded-md bg-gray-50"
                                            />
                                            <button
                                                onClick={copyShareUrl}
                                                className="px-3 py-2 border border-gray-300 rounded-md hover:bg-gray-50 flex items-center"
                                            >
                                                {getButtonState('copy_url').status === 'success' ? (
                                                    <Check className="w-4 h-4 text-green-600" />
                                                ) : (
                                                    <Copy className="w-4 h-4" />
                                                )}
                                            </button>
                                            <a
                                                href={shareUrl}
                                                target="_blank"
                                                rel="noopener noreferrer"
                                                className="px-3 py-2 border border-gray-300 rounded-md hover:bg-gray-50 flex items-center"
                                            >
                                                <ExternalLink className="w-4 h-4" />
                                            </a>
                                        </div>

                                        {/* QR Code */}
                                        {qrCode && (
                                            <div className="flex justify-center">
                                                <div className="text-center">
                                                    <img src={qrCode} alt="QR Code" className="w-32 h-32 border rounded" />
                                                    <p className="text-xs text-gray-500 mt-2">Scan to access itinerary</p>
                                                </div>
                                            </div>
                                        )}
                                    </div>
                                )}
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Footer with additional options */}
            <div className="px-6 py-4 bg-gray-50 border-t border-gray-200 rounded-b-lg">
                <div className="flex items-center justify-between text-sm text-gray-600">
                    <div className="flex items-center space-x-4">
                        <span className="flex items-center">
                            <Cloud className="w-4 h-4 mr-1" />
                            Cloud Storage
                        </span>
                        <span className="flex items-center">
                            <Smartphone className="w-4 h-4 mr-1" />
                            Mobile Friendly
                        </span>
                        <span className="flex items-center">
                            <Printer className="w-4 h-4 mr-1" />
                            Print Ready
                        </span>
                    </div>
                    <span className="text-xs">
                        Files expire in 30 days
                    </span>
                </div>
            </div>
        </div>
    )
}