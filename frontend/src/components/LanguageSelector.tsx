'use client'

import React, { useState, useEffect, useRef } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import {
    Globe,
    ChevronDown,
    Check,
    Settings,
    MapPin,
    Clock,
    DollarSign,
    Calendar,
    Users,
    Smartphone,
    Monitor,
    Volume2,
    Type,
    RotateCcw
} from 'lucide-react'

// Language configuration interface
interface LanguageConfig {
    code: string
    name: string
    nativeName: string
    flag: string
    currency: string
    currencySymbol: string
    dateFormat: string
    timeFormat: string
    numberFormat: string
    rtl: boolean
    timezone: string
    popularity: number
    region: string
    regionCode: string
}

// Available languages with comprehensive configuration
const languages: LanguageConfig[] = [
    {
        code: 'en',
        name: 'English',
        nativeName: 'English',
        flag: '🇺🇸',
        currency: 'USD',
        currencySymbol: '$',
        dateFormat: 'MM/DD/YYYY',
        timeFormat: '12h',
        numberFormat: '123,456',
        rtl: false,
        timezone: 'America/New_York',
        popularity: 100,
        region: 'North America',
        regionCode: 'NA'
    },
    {
        code: 'hi',
        name: 'Hindi',
        nativeName: 'हिंदी',
        flag: '🇮🇳',
        currency: 'INR',
        currencySymbol: '₹',
        dateFormat: 'DD/MM/YYYY',
        timeFormat: '12h',
        numberFormat: '1,23,456',
        rtl: false,
        timezone: 'Asia/Kolkata',
        popularity: 95,
        region: 'South Asia',
        regionCode: 'SA'
    },
    {
        code: 'bn',
        name: 'Bengali',
        nativeName: 'বাংলা',
        flag: '🇧🇩',
        currency: 'INR',
        currencySymbol: '₹',
        dateFormat: 'DD/MM/YYYY',
        timeFormat: '12h',
        numberFormat: '1,23,456',
        rtl: false,
        timezone: 'Asia/Kolkata',
        popularity: 85,
        region: 'South Asia',
        regionCode: 'SA'
    },
    {
        code: 'ta',
        name: 'Tamil',
        nativeName: 'தமிழ்',
        flag: '🇮🇳',
        currency: 'INR',
        currencySymbol: '₹',
        dateFormat: 'DD/MM/YYYY',
        timeFormat: '12h',
        numberFormat: '1,23,456',
        rtl: false,
        timezone: 'Asia/Kolkata',
        popularity: 80,
        region: 'South Asia',
        regionCode: 'SA'
    },
    {
        code: 'mr',
        name: 'Marathi',
        nativeName: 'मराठी',
        flag: '🇮🇳',
        currency: 'INR',
        currencySymbol: '₹',
        dateFormat: 'DD/MM/YYYY',
        timeFormat: '12h',
        numberFormat: '1,23,456',
        rtl: false,
        timezone: 'Asia/Kolkata',
        popularity: 75,
        region: 'South Asia',
        regionCode: 'SA'
    },
    {
        code: 'gu',
        name: 'Gujarati',
        nativeName: 'ગુજરાતી',
        flag: '🇮🇳',
        currency: 'INR',
        currencySymbol: '₹',
        dateFormat: 'DD/MM/YYYY',
        timeFormat: '12h',
        numberFormat: '1,23,456',
        rtl: false,
        timezone: 'Asia/Kolkata',
        popularity: 70,
        region: 'South Asia',
        regionCode: 'SA'
    },
    {
        code: 'te',
        name: 'Telugu',
        nativeName: 'తెలుగు',
        flag: '🇮🇳',
        currency: 'INR',
        currencySymbol: '₹',
        dateFormat: 'DD/MM/YYYY',
        timeFormat: '12h',
        numberFormat: '1,23,456',
        rtl: false,
        timezone: 'Asia/Kolkata',
        popularity: 72,
        region: 'South Asia',
        regionCode: 'SA'
    },
    {
        code: 'kn',
        name: 'Kannada',
        nativeName: 'ಕನ್ನಡ',
        flag: '🇮🇳',
        currency: 'INR',
        currencySymbol: '₹',
        dateFormat: 'DD/MM/YYYY',
        timeFormat: '12h',
        numberFormat: '1,23,456',
        rtl: false,
        timezone: 'Asia/Kolkata',
        popularity: 68,
        region: 'South Asia',
        regionCode: 'SA'
    },
    {
        code: 'ml',
        name: 'Malayalam',
        nativeName: 'മലയാളം',
        flag: '🇮🇳',
        currency: 'INR',
        currencySymbol: '₹',
        dateFormat: 'DD/MM/YYYY',
        timeFormat: '12h',
        numberFormat: '1,23,456',
        rtl: false,
        timezone: 'Asia/Kolkata',
        popularity: 65,
        region: 'South Asia',
        regionCode: 'SA'
    },
    {
        code: 'pa',
        name: 'Punjabi',
        nativeName: 'ਪੰਜਾਬੀ',
        flag: '🇮🇳',
        currency: 'INR',
        currencySymbol: '₹',
        dateFormat: 'DD/MM/YYYY',
        timeFormat: '12h',
        numberFormat: '1,23,456',
        rtl: false,
        timezone: 'Asia/Kolkata',
        popularity: 63,
        region: 'South Asia',
        regionCode: 'SA'
    }
]

interface LanguageSelectorProps {
    currentLanguage?: string
    onLanguageChange?: (language: LanguageConfig) => void
    showRegionalSettings?: boolean
    className?: string
    userId?: string
}

export default function LanguageSelector({
    currentLanguage = 'en',
    onLanguageChange,
    showRegionalSettings = true,
    className = '',
    userId
}: LanguageSelectorProps) {
    const [isOpen, setIsOpen] = useState(false)
    const [searchTerm, setSearchTerm] = useState('')
    const [selectedLanguage, setSelectedLanguage] = useState<LanguageConfig>(
        languages.find(lang => lang.code === currentLanguage) || languages[0]
    )
    const [showSettings, setShowSettings] = useState(false)
    const [regionalSettings, setRegionalSettings] = useState({
        autoDetectLocation: true,
        useLocalCurrency: true,
        useLocalTimeZone: true,
        useLocalDateFormat: true,
        useLocalNumberFormat: true
    })
    const [detectedLocation, setDetectedLocation] = useState<{
        country?: string
        region?: string
        timezone?: string
        currency?: string
    }>({})
    const dropdownRef = useRef<HTMLDivElement>(null)

    // Auto-detect user location and language
    useEffect(() => {
        const detectLocation = async () => {
            try {
                // Detect timezone
                const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone

                // Detect locale preferences
                const locale = navigator.language || navigator.languages[0]
                const [langCode] = locale.split('-')

                // Try to detect country/region from timezone or locale
                const regionInfo = {
                    timezone,
                    country: locale.split('-')[1] || 'Unknown',
                    region: timezone.split('/')[0] || 'Unknown'
                }

                setDetectedLocation(regionInfo)

                // Auto-select language if available and auto-detect is enabled
                if (regionalSettings.autoDetectLocation) {
                    const detectedLang = languages.find(lang =>
                        lang.code === langCode || lang.code === locale
                    )
                    if (detectedLang && detectedLang.code !== selectedLanguage.code) {
                        handleLanguageSelect(detectedLang)
                    }
                }
            } catch (error) {
                console.error('Location detection failed:', error)
            }
        }

        detectLocation()
    }, [regionalSettings.autoDetectLocation])

    // Close dropdown when clicking outside
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
                setIsOpen(false)
                setShowSettings(false)
            }
        }

        document.addEventListener('mousedown', handleClickOutside)
        return () => document.removeEventListener('mousedown', handleClickOutside)
    }, [])

    // Save language preference
    const saveLanguagePreference = async (language: LanguageConfig) => {
        try {
            if (userId) {
                await fetch('/api/user/language-preference', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        userId,
                        language: language.code,
                        regionalSettings,
                        detectedLocation
                    })
                })
            }

            // Save to localStorage as backup
            localStorage.setItem('preferredLanguage', JSON.stringify({
                language: language.code,
                regionalSettings,
                timestamp: new Date().toISOString()
            }))
        } catch (error) {
            console.error('Failed to save language preference:', error)
        }
    }

    // Handle language selection
    const handleLanguageSelect = async (language: LanguageConfig) => {
        setSelectedLanguage(language)
        setIsOpen(false)
        setSearchTerm('')

        // Update regional settings based on language
        if (regionalSettings.useLocalCurrency || regionalSettings.useLocalTimeZone ||
            regionalSettings.useLocalDateFormat || regionalSettings.useLocalNumberFormat) {
            const updatedSettings = { ...regionalSettings }

            // You could update other settings here based on the selected language
            setRegionalSettings(updatedSettings)
        }

        // Save preference
        await saveLanguagePreference(language)

        // Notify parent component
        if (onLanguageChange) {
            onLanguageChange(language)
        }

        // Update page direction for RTL languages
        document.documentElement.dir = language.rtl ? 'rtl' : 'ltr'
        document.documentElement.lang = language.code
    }

    // Filter languages based on search
    const filteredLanguages = languages.filter(lang =>
        lang.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        lang.nativeName.toLowerCase().includes(searchTerm.toLowerCase()) ||
        lang.code.toLowerCase().includes(searchTerm.toLowerCase())
    )

    // Group languages by region
    const groupedLanguages = filteredLanguages.reduce((groups, lang) => {
        const region = lang.region
        if (!groups[region]) groups[region] = []
        groups[region].push(lang)
        return groups
    }, {} as Record<string, LanguageConfig[]>)

    // Format example values
    const formatExamples = {
        currency: `${selectedLanguage.currencySymbol}1,234.56`,
        date: new Date().toLocaleDateString('en-US', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric'
        }).replace(/\//g, selectedLanguage.dateFormat.includes('DD/MM') ? '/' : '/'),
        time: new Date().toLocaleTimeString('en-US', {
            hour12: selectedLanguage.timeFormat === '12h'
        }),
        number: selectedLanguage.numberFormat === '1,23,456' ? '1,23,456.78' : '123,456.78'
    }

    return (
        <div className={`relative ${className}`} ref={dropdownRef}>
            {/* Language Selector Button */}
            <button
                onClick={() => setIsOpen(!isOpen)}
                className="flex items-center px-3 py-2 text-sm text-gray-700 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
                aria-label="Select Language"
            >
                <Globe className="w-4 h-4 mr-2" />
                <span className="mr-1">{selectedLanguage.flag}</span>
                <span className="hidden sm:inline">{selectedLanguage.nativeName}</span>
                <span className="sm:hidden">{selectedLanguage.code.toUpperCase()}</span>
                <ChevronDown className={`w-4 h-4 ml-1 transition-transform ${isOpen ? 'rotate-180' : ''}`} />
            </button>

            {/* Dropdown Menu */}
            <AnimatePresence>
                {isOpen && (
                    <motion.div
                        initial={{ opacity: 0, y: -10, scale: 0.95 }}
                        animate={{ opacity: 1, y: 0, scale: 1 }}
                        exit={{ opacity: 0, y: -10, scale: 0.95 }}
                        className="absolute right-0 top-full mt-2 w-80 bg-white rounded-lg shadow-xl border border-gray-200 z-50"
                    >
                        {/* Header */}
                        <div className="p-4 border-b border-gray-200">
                            <div className="flex items-center justify-between mb-3">
                                <h3 className="font-semibold text-gray-900">Select Language</h3>
                                {showRegionalSettings && (
                                    <button
                                        onClick={() => setShowSettings(!showSettings)}
                                        className="text-gray-500 hover:text-gray-700"
                                        aria-label="Regional Settings"
                                    >
                                        <Settings className="w-4 h-4" />
                                    </button>
                                )}
                            </div>

                            {/* Search */}
                            <input
                                type="text"
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                                placeholder="Search languages..."
                                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                            />
                        </div>

                        {/* Regional Settings */}
                        <AnimatePresence>
                            {showSettings && (
                                <motion.div
                                    initial={{ opacity: 0, height: 0 }}
                                    animate={{ opacity: 1, height: 'auto' }}
                                    exit={{ opacity: 0, height: 0 }}
                                    className="border-b border-gray-200 p-4 bg-gray-50"
                                >
                                    <h4 className="font-medium text-gray-900 mb-3">Regional Settings</h4>
                                    <div className="space-y-3">
                                        <label className="flex items-center justify-between">
                                            <span className="text-sm text-gray-700">Auto-detect location</span>
                                            <input
                                                type="checkbox"
                                                checked={regionalSettings.autoDetectLocation}
                                                onChange={(e) => setRegionalSettings(prev => ({
                                                    ...prev,
                                                    autoDetectLocation: e.target.checked
                                                }))}
                                                className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                            />
                                        </label>
                                        <label className="flex items-center justify-between">
                                            <span className="text-sm text-gray-700">Use local currency</span>
                                            <input
                                                type="checkbox"
                                                checked={regionalSettings.useLocalCurrency}
                                                onChange={(e) => setRegionalSettings(prev => ({
                                                    ...prev,
                                                    useLocalCurrency: e.target.checked
                                                }))}
                                                className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                            />
                                        </label>
                                        <label className="flex items-center justify-between">
                                            <span className="text-sm text-gray-700">Use local timezone</span>
                                            <input
                                                type="checkbox"
                                                checked={regionalSettings.useLocalTimeZone}
                                                onChange={(e) => setRegionalSettings(prev => ({
                                                    ...prev,
                                                    useLocalTimeZone: e.target.checked
                                                }))}
                                                className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                            />
                                        </label>
                                        <label className="flex items-center justify-between">
                                            <span className="text-sm text-gray-700">Use local date format</span>
                                            <input
                                                type="checkbox"
                                                checked={regionalSettings.useLocalDateFormat}
                                                onChange={(e) => setRegionalSettings(prev => ({
                                                    ...prev,
                                                    useLocalDateFormat: e.target.checked
                                                }))}
                                                className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                            />
                                        </label>
                                    </div>

                                    {/* Current Settings Preview */}
                                    <div className="mt-4 p-3 bg-white rounded border">
                                        <h5 className="text-xs font-medium text-gray-700 mb-2">Format Examples:</h5>
                                        <div className="grid grid-cols-2 gap-2 text-xs">
                                            <div className="flex items-center">
                                                <DollarSign className="w-3 h-3 mr-1 text-gray-500" />
                                                <span>{formatExamples.currency}</span>
                                            </div>
                                            <div className="flex items-center">
                                                <Calendar className="w-3 h-3 mr-1 text-gray-500" />
                                                <span>{formatExamples.date}</span>
                                            </div>
                                            <div className="flex items-center">
                                                <Clock className="w-3 h-3 mr-1 text-gray-500" />
                                                <span>{formatExamples.time}</span>
                                            </div>
                                            <div className="flex items-center">
                                                <Type className="w-3 h-3 mr-1 text-gray-500" />
                                                <span>{formatExamples.number}</span>
                                            </div>
                                        </div>
                                    </div>
                                </motion.div>
                            )}
                        </AnimatePresence>

                        {/* Language List */}
                        <div className="max-h-64 overflow-y-auto">
                            {Object.entries(groupedLanguages).map(([region, langs]) => (
                                <div key={region}>
                                    <div className="px-4 py-2 text-xs font-medium text-gray-500 bg-gray-50 sticky top-0">
                                        {region}
                                    </div>
                                    {langs
                                        .sort((a, b) => b.popularity - a.popularity)
                                        .map((language) => (
                                            <motion.button
                                                key={language.code}
                                                onClick={() => handleLanguageSelect(language)}
                                                className={`w-full px-4 py-3 text-left hover:bg-gray-50 transition-colors flex items-center justify-between ${selectedLanguage.code === language.code ? 'bg-blue-50 text-blue-700' : 'text-gray-900'
                                                    }`}
                                                whileHover={{ x: 4 }}
                                            >
                                                <div className="flex items-center">
                                                    <span className="text-lg mr-3">{language.flag}</span>
                                                    <div>
                                                        <div className="font-medium">{language.nativeName}</div>
                                                        <div className="text-sm text-gray-500">{language.name}</div>
                                                    </div>
                                                </div>
                                                <div className="flex items-center space-x-2">
                                                    <div className="text-right text-xs text-gray-500">
                                                        <div>{language.currencySymbol} • {language.dateFormat}</div>
                                                        <div>{language.timezone.split('/')[1]}</div>
                                                    </div>
                                                    {selectedLanguage.code === language.code && (
                                                        <Check className="w-4 h-4 text-blue-600" />
                                                    )}
                                                </div>
                                            </motion.button>
                                        ))}
                                </div>
                            ))}
                        </div>

                        {/* Detected Location Info */}
                        {detectedLocation.timezone && (
                            <div className="p-3 border-t border-gray-200 bg-gray-50 text-xs text-gray-600">
                                <div className="flex items-center">
                                    <MapPin className="w-3 h-3 mr-1" />
                                    <span>Detected: {detectedLocation.timezone}</span>
                                    {regionalSettings.autoDetectLocation && (
                                        <span className="ml-2 text-green-600">• Auto-enabled</span>
                                    )}
                                </div>
                            </div>
                        )}

                        {/* Footer */}
                        <div className="p-3 border-t border-gray-200 flex items-center justify-between text-xs text-gray-500">
                            <span>{filteredLanguages.length} languages available</span>
                            <button
                                onClick={() => {
                                    setSearchTerm('')
                                    setShowSettings(false)
                                }}
                                className="text-blue-600 hover:text-blue-800 flex items-center"
                            >
                                <RotateCcw className="w-3 h-3 mr-1" />
                                Reset
                            </button>
                        </div>
                    </motion.div>
                )}
            </AnimatePresence>
        </div>
    )
}